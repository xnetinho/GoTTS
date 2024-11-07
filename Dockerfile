FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
RUN mkdir -p /app/voices

EXPOSE 8080
CMD ["./main"]