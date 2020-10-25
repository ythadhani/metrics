FROM golang:1.15.3-alpine AS build

WORKDIR /go/src/github.com/ythadhani/metrics/

ADD ./metrics /go/src/github.com/ythadhani/metrics

EXPOSE 9090

ENTRYPOINT ["./metrics"]