FROM golang:1.14.12-alpine3.12

RUN mkdir /go/src/task_checker

WORKDIR /go/src/task_checker

ADD . /go/src/task_checker

RUN go get -v golang.org/x/tools/gopls 
RUN go get -v github.com/go-delve/delve/cmd/dlv
RUN apk add gcc alpine-sdk
ENV PORT=8080
