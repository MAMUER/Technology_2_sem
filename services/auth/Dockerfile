FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto ./proto
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/auth-service ./services/auth/cmd/auth

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder --chown=appuser:appgroup /app/bin/auth-service .
COPY --from=builder --chown=appuser:appgroup /app/proto/gen/go ./proto/gen/go

USER appuser

EXPOSE 8081 50051

CMD ["./auth-service"]