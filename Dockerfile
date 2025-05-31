FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o chat_socio ./cmd/

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /app/chat_socio .
COPY --from=builder /app/migrations ./migrations

CMD ["./chat_socio"]