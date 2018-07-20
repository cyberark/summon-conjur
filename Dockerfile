FROM golang:1.10
MAINTAINER Conjur Inc

RUN apt-get update && apt-get install -y jq less 
RUN go get -u github.com/jstemmer/go-junit-report
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get github.com/playscale/goconvey

RUN mkdir -p /go/src/github.com/cyberark/summon-conjur/output
WORKDIR /go/src/github.com/cyberark/summon-conjur

COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 8080
