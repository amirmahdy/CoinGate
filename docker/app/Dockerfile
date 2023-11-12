FROM golang:1.20.4 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY docker/bash/start.sh .
COPY docker/bash/wait-for.sh .
COPY db/migrations ./db/migrations