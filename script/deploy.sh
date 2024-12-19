#!/bin/bash

APP_NAME="go-signin-service"            # 应用名称
BUILD_DIR="."                            # 本地构建目录
TARGET_OS="linux"                        # 目标操作系统
TARGET_ARCH="amd64"                      # 目标架构
SERVER_USER="ubuntu"                     # 目标服务器用户名
SERVER_HOST="52.221.41.163"              # 目标服务器IP
DEPLOY_DIR="/data/go-signin-service/deploy"  # 部署目录
APP_PORT=8090                            # 应用运行端口

echo "Cleaning old build files..."
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

echo "Building application for $TARGET_OS/$TARGET_ARCH..."
GOOS=$TARGET_OS GOARCH=$TARGET_ARCH go build -o $BUILD_DIR/$APP_NAME main.go
if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi
echo "Build successful!"

# 上传二进制文件之前先删除远程的旧文件
echo "Deleting old binary on server..."
ssh $SERVER_USER@$SERVER_HOST "if [ -f $DEPLOY_DIR/$APP_NAME ]; then rm -f $DEPLOY_DIR/$APP_NAME; fi"
if [ $? -ne 0 ]; then
    echo "Failed to delete old binary on server!"
    exit 1
fi
echo "Old binary deleted successfully!"

echo "Uploading new binary to $SERVER_USER@$SERVER_HOST..."
scp $BUILD_DIR/$APP_NAME $SERVER_USER@$SERVER_HOST:$DEPLOY_DIR/
if [ $? -ne 0 ]; then
    echo "Upload failed!"
    exit 1
fi
echo "Upload successful!"

echo "Restarting service on server..."
ssh $SERVER_USER@$SERVER_HOST <<EOF
    echo "Stopping old service..."
    OLD_PID=\$(lsof -t -i:$APP_PORT)
    if [ ! -z "\$OLD_PID" ]; then
        kill -9 \$OLD_PID
        echo "Old service stopped!"
    else
        echo "No old service running."
    fi

    # 启动新服务
    echo "Starting new service..."
    cd $DEPLOY_DIR
    nohup ./$APP_NAME > app.log 2>&1 &
    echo "New service started!"
EOF

echo "Deployment completed!"
