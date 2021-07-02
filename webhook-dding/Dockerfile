FROM golang:alpine

ADD webhook-dding /
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add tzdata \
    && rm -rf /var/cache/apk/* \
    && chmod +x /webhook-dding

CMD /webhook-dding