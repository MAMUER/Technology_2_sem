FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tasks-service ./services/tasks/cmd/tasks

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tasks-service .

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

USER appuser

EXPOSE 8082

CMD ["./tasks-service"]