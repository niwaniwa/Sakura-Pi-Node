FROM golang:latest

WORKDIR /app

# 前提インストール
RUN dpkg --add-architecture armhf
RUN apt update && apt upgrade -y
RUN apt install pkg-config
RUN apt-get install -y gcc-arm-linux-gnueabihf libusb-1.0-0-dev libusb-1.0-0-dev:armhf

ENV GOOS linux
ENV GOARCH arm
ENV GOARM 7
ENV CGO_ENABLED 1
ENV CC arm-linux-gnueabihf-gcc

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# ソースコードのコピー
COPY . .

# ビルド
RUN go build -o /main ./cmd

RUN ls -l /app

CMD [ "/main" ]