FROM golang:1.11.4-alpine

MAINTAINER Conjur Inc


RUN apk add --no-cache bash \
                       build-base \
                       curl \
                       git \
                       jq \
                       less && \
    go get -u github.com/jstemmer/go-junit-report && \
    go get -u github.com/smartystreets/goconvey && \
    mkdir -p /summon-conjur/output

WORKDIR /summon-conjur

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o summon-conjur cmd/main.go

EXPOSE 8080
