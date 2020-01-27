FROM golang:1.13-alpine

RUN apk update && apk add git

RUN go get -v golang.org/x/tools/cmd/godoc

WORKDIR /go/src/github.com/dcrichards/go-fhrs

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
