FROM golang:1.25-alpine AS builder

# Устанавливаем protoc и плагины
RUN apk add --no-cache protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем proto файлы и генерируем код
COPY proto ./proto
RUN mkdir -p proto/gen/go
RUN protoc --proto_path=proto \
    --go_out=proto/gen/go --go_opt=paths=source_relative \
    --go-grpc_out=proto/gen/go --go-grpc_opt=paths=source_relative \
    proto/auth.proto

COPY . .

RUN go build -o tasks-service ./services/tasks/cmd/tasks

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/tasks-service .
COPY --from=builder /app/proto/gen/go ./proto/gen/go

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

USER appuser

EXPOSE 8082

CMD ["./tasks-service"]