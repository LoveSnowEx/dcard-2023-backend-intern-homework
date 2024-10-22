FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./tmp/main .

FROM caddy:2.8-alpine

COPY --from=builder /app/tmp/main ./main

COPY Caddyfile /etc/caddy/Caddyfile

CMD ["sh", "-c", "caddy run --config /etc/caddy/Caddyfile & ./main"]
