FROM golang:1.14.4 as builder

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src
RUN go generate ./src

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/bin/main \
    ./src

FROM scratch as runner

WORKDIR /usr/local/bin/

COPY --from=builder /go/bin/main /main

ENTRYPOINT ["/main"]
