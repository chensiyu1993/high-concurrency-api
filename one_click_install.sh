#!/bin/bash

echo "=== 开始安装 ==="

# 使用国内源安装软件
sudo apt update
sudo apt install -y wget curl unzip
sudo DEBIAN_FRONTEND=noninteractive apt install -y mysql-server redis-server

# 使用国内Go镜像
wget https://mirrors.aliyun.com/golang/go1.20.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 下载项目代码
curl -L https://github.com/chensiyu1993/high-concurrency-api/archive/refs/heads/master.zip -o master.zip
unzip master.zip
mv high-concurrency-api-master high-concurrency-api
cd high-concurrency-api

# 设置MySQL
sudo mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'chensiyu';"
sudo mysql -e "CREATE DATABASE IF NOT EXISTS high_concurrency_db;"

# 启动服务
sudo systemctl restart mysql
sudo systemctl restart redis

# 设置Go环境
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 编译
go mod tidy
go build -o api-server

# 创建系统服务
sudo tee /etc/systemd/system/api-server.service << EOF
[Unit]
Description=High Concurrency API Server
After=network.target mysql.service redis.service

[Service]
Type=simple
WorkingDirectory=$(pwd)
ExecStart=$(pwd)/api-server
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable api-server
sudo systemctl start api-server

# 等待服务启动
sleep 3

# 检查服务状态
sudo systemctl status api-server

echo "=== 安装完成 ==="
echo "MySQL密码: chensiyu"
echo "服务已启动在: http://服务器IP:8080"
echo ""
echo "管理命令："
echo "1. 查看服务状态：sudo systemctl status api-server"
echo "2. 重启服务：sudo systemctl restart api-server"
echo "3. 查看日志：sudo journalctl -u api-server -f" 