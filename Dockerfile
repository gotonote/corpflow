FROM golang:1.21-alpine AS builder

WORKDIR /app

# 设置 Go 国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 和空的 go.sum
COPY go.mod ./
COPY go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码并构建
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
