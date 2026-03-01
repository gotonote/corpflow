FROM golang:1.21-alpine AS builder

WORKDIR /app

# 先只复制 go.mod，下载依赖后才会生成 go.sum
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/internal ./internal

EXPOSE 8080

CMD ["./server"]
