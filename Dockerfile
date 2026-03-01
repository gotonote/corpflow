FROM golang:1.21-alpine AS builder

WORKDIR /app

# 设置 Go 国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 复制所有源代码
COPY . ./

# 使用 -mod=mod 自动下载缺失的依赖，绕过 go.sum 验证
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o server ./cmd/server
