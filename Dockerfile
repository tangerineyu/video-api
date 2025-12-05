FROM golang:1.25-alpine AS builder
#设置工作目录
WORKDIR /app
# 设置代理
ENV GOPROXY=https://goproxy.cn,direct
# 复制go.mod和go.sum文件
COPY go.mod go.sum ./
RUN go mod download
# 复制源代码
COPY . .
# 编译为名为video-api的二进制此文件
RUN go build -o video-api main.go

#运行阶段
FROM alpine:latest
# 设置工作目录
WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# 安装必要的系统库
RUN apk --no-cache add ca-certificates tzdata
# 从构建阶段复制二进制文件
COPY --from=builder /app/video-api .
# 创建上传目录
RUN mkdir -p uploads/avatars uploads/videos
# 暴露端口
EXPOSE 8080
# 运行二进制文件
CMD ["./video-api"]

