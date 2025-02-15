FROM golang:1.22 AS builder

WORKDIR /app
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./internal/cmd/main.go --output ./docs

RUN CGO_ENABLED=0 GOOS=linux go build -o /build ./internal/cmd \
    && go clean -cache -modcache

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build .
COPY --from=builder /app/docs ./docs

EXPOSE 8080
CMD ["./build"]