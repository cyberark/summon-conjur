FROM golang:1.8
MAINTAINER Conjur Inc

RUN go get -u github.com/jstemmer/go-junit-report
RUN go get github.com/tools/godep
RUN go get github.com/smartystreets/goconvey

RUN mkdir -p /go/src/github.com/cyberark/summon-conjur/output
WORKDIR /go/src/github.com/cyberark/summon-conjur

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 8080
