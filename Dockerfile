FROM golang:1.21-alpine AS builder

WORKDIR /app

# 设置 Go 国内代理
RUN go env -w GOPROXY=https://goproxy.cn,direct

# 创建初始的 go.sum（空文件即可，go get 会填充）
RUN touch go.sum

# 复制 go.mod 并使用 go get 获取依赖（会生成完整的 go.sum）
COPY go.mod ./
RUN go get ./...

# 复制所有源代码并构建
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
