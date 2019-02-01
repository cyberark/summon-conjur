FROM golang:1.12-rc-alpine

MAINTAINER Conjur Inc

RUN apk add --no-cache bash \
                       build-base \
                       git \
                       jq \
                       less
RUN go get -u github.com/jstemmer/go-junit-report
RUN go get github.com/playscale/goconvey

RUN mkdir -p /summon-conjur/output
WORKDIR /summon-conjur

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 8080
