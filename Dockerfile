FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /build ./internal/cmd \
    && go clean -cache -modcache

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build .

EXPOSE 8080
CMD ["./build"]