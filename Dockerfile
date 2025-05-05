FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o goProxy /app/cmd/server/main.go

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/goProxy .

CMD ["./goProxy"]
