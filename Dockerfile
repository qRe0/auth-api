FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o auth-api cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

COPY --from=builder /app/auth-api .
COPY --from=builder /app/.env .
COPY --from=builder /app/internal/migrations /root/internal/migrations

CMD ["./auth-api"]