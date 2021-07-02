package utils

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestGetToken(t *testing.T) {
	fmt.Println(GetSkipKey())
}

func TestGetSkipKey(t *testing.T) {
	keyword := "log, es"
	var ok bool
	for _, key := range strings.Split(keyword, ",") {
		fmt.Println(key)
		ok, _ = regexp.MatchString(key, keyword)
		if ok {
			break
		}
	}
	fmt.Println(ok)
}