# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20

# 安装必要的依赖（包括 FFmpeg 开发库）
RUN apt-get update && \
    apt-get install -y ffmpeg libavcodec-dev libavformat-dev libavutil-dev && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 将本地代码复制到容器中
COPY . .

# 设置 CGO 编译参数，指向 FFmpeg 的头文件和库文件路径
ENV CGO_CFLAGS="-I/usr/include"
ENV CGO_LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lavcodec -lavformat -lavutil"

# 下载 Go 依赖包
RUN go mod download

# 编译 Go 应用程序
RUN go build -o app .

# 暴露端口（根据你的需求）
EXPOSE 8080

# 启动应用程序
CMD ["./app"]
