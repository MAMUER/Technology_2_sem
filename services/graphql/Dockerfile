FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект, потому что GraphQL зависит от shared и других пакетов
COPY . .

# Правильный путь к main.go - он находится в services/graphql/cmd/graphql
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/graphql-service ./services/graphql/cmd/graphql

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/bin/graphql-service .

USER appuser

EXPOSE 8090

CMD ["./graphql-service"]