FROM golang:1.8
MAINTAINER Conjur Inc

RUN go get -u github.com/jstemmer/go-junit-report
RUN go get github.com/tools/godep
RUN go get github.com/smartystreets/goconvey

WORKDIR /go/src/github.com/conjurinc/summon-conjur

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 8080

CMD /bin/bash
