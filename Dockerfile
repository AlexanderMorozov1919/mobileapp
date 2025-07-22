FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /mobileapp ./cmd/app/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /mobileapp .
COPY --from=builder /app/.env .
EXPOSE ${APP_PORT:-8080}
ENTRYPOINT ["/app/mobileapp"]