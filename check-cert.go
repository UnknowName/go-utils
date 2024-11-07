package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultConcurrency = 1
	errExpiringShortly = "%s: ** '%s' expires in %d hours! **"
	errExpiringSoon    = "%s: '%s' expires in roughly %d days."
	errSunsetAlg       = "%s: '%s' expires after the sunset date for its signature algorithm '%s'."
	sendFmt            = "https://oapi.dingtalk.com/robot/send?access_token="
	contentType        = "application/json"
)

type sigAlgSunset struct {
	name      string    // Human readable name of signature algorithm
	sunsetsAt time.Time // Time the algorithm will be sunset
}

// sunsetSigAlgs is an algorithm to string mapping for signature algorithms
// which have been or are being deprecated.  See the following links to learn
// more about SHA1's inclusion on this list.
//
// - https://technet.microsoft.com/en-us/library/security/2880823.aspx
// - http://googleonlinesecurity.blogspot.com/2014/09/gradually-sunsetting-sha-1.html
var sunsetSigAlgs = map[x509.SignatureAlgorithm]sigAlgSunset{
	x509.MD2WithRSA: sigAlgSunset{
		name:      "MD2 with RSA",
		sunsetsAt: time.Now(),
	},
	x509.MD5WithRSA: sigAlgSunset{
		name:      "MD5 with RSA",
		sunsetsAt: time.Now(),
	},
	x509.SHA1WithRSA: sigAlgSunset{
		name:      "SHA1 with RSA",
		sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	x509.DSAWithSHA1: sigAlgSunset{
		name:      "DSA with SHA1",
		sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	x509.ECDSAWithSHA1: sigAlgSunset{
		name:      "ECDSA with SHA1",
		sunsetsAt: time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
	},
}

var (
	hostsFile   = flag.String("hosts", "hosts", "The path to the file containing a list of hosts to check.")
	warnYears   = flag.Int("years", 0, "Warn if the certificate will expire within this many years.")
	warnMonths  = flag.Int("months", 0, "Warn if the certificate will expire within this many months.")
	warnDays    = flag.Int("days", 0, "Warn if the certificate will expire within this many days.")
	checkSigAlg = flag.Bool("check-sig-alg", true, "Verify that non-root certificates are using a good signature algorithm.")
	concurrency = flag.Int("concurrency", defaultConcurrency, "Maximum number of hosts to check at once.")
	ddingToken  = flag.String("token", "", "dding token")
	sendRecords  map[string][]string
)

type certErrors struct {
	commonName string
	errs       []error
}


type hostResult struct {
	host  string
	err   error
	certs []certErrors
}

func main() {
	sendRecords = map[string][]string{}
	flag.Parse()
	if len(*hostsFile) == 0 {
		flag.Usage()
		return
	}
	if *warnYears < 0 {
		*warnYears = 0
	}
	if *warnMonths < 0 {
		*warnMonths = 0
	}
	if *warnDays < 0 {
		*warnDays = 0
	}
	if *warnYears == 0 && *warnMonths == 0 && *warnDays == 0 {
		*warnDays = 20
	}
	if *concurrency < 0 {
		*concurrency = defaultConcurrency
	}
	for {
		log.Println("start check")
		processHosts()
		for _, msgs := range sendRecords {
			message := "\n"
			for _, msg := range msgs {
				message = fmt.Sprintf("%s%s\n", message, msg)
			}
			// log.Println("message", message)
			sendMsg(*ddingToken, message)
		}
		time.Sleep(time.Hour * 24)
	}
}

func processHosts() {
	done := make(chan struct{})
	defer close(done)

	hosts := queueHosts(done)
	results := make(chan hostResult)
	var wg sync.WaitGroup
	wg.Add(*concurrency)
	for i := 0; i < *concurrency; i++ {
		go func() {
			processQueue(done, hosts, results)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()


	for r := range results {
		if r.err != nil {
			log.Printf("%s: %v\n", r.host, r.err)
			continue
		}
		for _, cert := range r.certs {
			if cert.errs == nil || len(cert.errs) == 0 {
				continue
			}
			msgs := sendRecords[cert.commonName]
			if msgs == nil {
				msgs = make([]string, 0)
			}
			for _, err := range cert.errs {
				msgs = append(msgs, err.Error())
			}
			sendRecords[cert.commonName] = msgs
		}
	}
}

func queueHosts(done <-chan struct{}) <-chan string {
	hosts := make(chan string)
	go func() {
		defer close(hosts)

		fileContents, err := os.ReadFile(*hostsFile)
		if err != nil {
			return
		}
		lines := strings.Split(string(fileContents), "\n")
		for _, line := range lines {
			host := strings.TrimSpace(line)
			if len(host) == 0 || host[0] == '#' {
				continue
			}
			select {
			case hosts <- host:
			case <-done:
				return
			}
		}
	}()
	return hosts
}


func processQueue(done <-chan struct{}, hosts <-chan string, results chan<- hostResult) {
	for host := range hosts {
		select {
		case results <- checkHost(host):
		case <-done:
			return
		}
	}
}

func checkHost(host string) (result hostResult) {
	result = hostResult{
		host:  host,
		certs: []certErrors{},
	}
	conn, err := tls.Dial("tcp", host, nil)
	if err != nil {
		result.err = err
		return
	}
	defer conn.Close()

	timeNow := time.Now()
	checkedCerts := make(map[string]struct{})
	for _, chain := range conn.ConnectionState().VerifiedChains {
		for certNum, cert := range chain {
			if _, checked := checkedCerts[string(cert.Signature)]; checked {
				continue
			}
			checkedCerts[string(cert.Signature)] = struct{}{}
			cErrs := []error{}

			// Check the expiration.
			if timeNow.AddDate(*warnYears, *warnMonths, *warnDays).After(cert.NotAfter) {
				expiresIn := int64(cert.NotAfter.Sub(timeNow).Hours())
				if expiresIn <= 48 {
					cErrs = append(cErrs, fmt.Errorf(errExpiringShortly, host, cert.Subject.CommonName, expiresIn))
				} else {
					cErrs = append(cErrs, fmt.Errorf(errExpiringSoon, host, cert.Subject.CommonName, expiresIn/24))
				}
			}

			// Check the signature algorithm, ignoring the root certificate.
			if alg, exists := sunsetSigAlgs[cert.SignatureAlgorithm]; *checkSigAlg && exists && certNum != len(chain)-1 {
				if cert.NotAfter.Equal(alg.sunsetsAt) || cert.NotAfter.After(alg.sunsetsAt) {
					cErrs = append(cErrs, fmt.Errorf(errSunsetAlg, host, cert.Subject.CommonName, alg.name))
				}
			}
			// log.Println("host", host, "commonName=", cert.Subject.CommonName)
			result.certs = append(result.certs, certErrors{
				commonName: cert.Subject.CommonName,
				errs:       cErrs,
			})
		}
	}
	return
}



func sendMsg(token, msg string) {
	log.Println("wait send msg", msg)
	_msg := NewTMessage(msg, nil, false)
	httpClient := http.Client{Timeout: time.Second * 5}
	url := fmt.Sprintf("%s%s", sendFmt, token)
	resp, err := httpClient.Post(url, contentType, bytes.NewBuffer(_msg.Encode()))
	if err != nil {
		log.Println("send failed", err)
		return
	}
	defer resp.Body.Close()
	_resp, _ := io.ReadAll(resp.Body)
	log.Println("发送响应", string(_resp))
}


type TMessage struct {
	MsgType string  `json:"msgtype"`
	Text    Content `json:"text"`
	At      At      `json:"at"`
	IsAtAll bool    `json:"isAtAll"`
}

type Content struct {
	Content string `json:"content"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
}

func NewTMessage(msg string, atMobiles []string, atAll bool) *TMessage {
	if atMobiles == nil {
		atMobiles = make([]string, 0)
	}
	atUsers := At{AtMobiles: atMobiles}
	text := Content{Content: msg}
	return &TMessage{
		MsgType: "text",
		Text:    text,
		At:      atUsers,
		IsAtAll: atAll,
	}
}

func (tm *TMessage) Encode() []byte {
	_bytes, err := json.Marshal(&tm)
	if err != nil {
		fmt.Println("Encode to json err ", err)
		return nil
	}
	return _bytes
}