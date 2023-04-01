# syntax=docker/dockerfile:1

## build
FROM golang:1.18.5-buster as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /token-fetcher

## Runner
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /token-fetcher /token-fetcher

ENTRYPOINT ["/token-fetcher"]