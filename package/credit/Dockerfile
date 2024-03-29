FROM golang:1.18.0 as builder

WORKDIR /go/app

RUN apt-get -y update && apt-get install -y ca-certificates

COPY ./pkg/ ./pkg/
COPY ./cmd/credit/ ./cmd/credit/
COPY ./go.mod  ./
COPY ./go.sum  ./

ARG CLOUD=aws
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN cd cmd/credit && \
    go build \
    -o /go/app/bin/main \
    -ldflags '-s -w -X main.target='$CLOUD

FROM scratch as runner

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/app/bin/main /app/main

EXPOSE 8080
ENTRYPOINT ["/app/main"]
