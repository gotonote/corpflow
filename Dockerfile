FROM golang:1.21-alpine AS builder

WORKDIR /app

# 设置 Go 国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制 go.mod 并下载依赖生成 go.sum
COPY go.mod ./
RUN go mod download

# 复制所有源代码（不包括 go.sum，让 Go 自动处理）
COPY . ./

# 运行 go mod tidy 确保依赖完整
RUN go mod tidy

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
