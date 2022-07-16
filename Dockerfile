FROM golang:1.18.1 AS builder
WORKDIR /builder
COPY ./ /builder

WORKDIR /builder
RUN go build -o ./deploy/bin/mqtt-metrics ./deploy/app

FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /builder/deploy/bin /app/
