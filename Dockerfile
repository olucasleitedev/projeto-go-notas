# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
ARG SERVICE=gateway
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/${SERVICE}

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
