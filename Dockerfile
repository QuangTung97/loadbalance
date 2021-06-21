FROM golang:1.16-buster as builder
WORKDIR /loadbalance
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o bin/server cmd/server/main.go
RUN go build -o bin/client cmd/client/main.go


## Today ubuntu is using minimalized image by default, using ubuntu for better compatible than alpine
FROM ubuntu:20.04
WORKDIR /loadbalance
COPY --from=builder /loadbalance/bin /loadbalance

EXPOSE 5001/tcp 7001/tcp 7001/udp
