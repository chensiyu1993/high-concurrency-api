#!/bin/bash

echo "=== 开始安装 ==="

# 安装所需软件
sudo apt update
sudo DEBIAN_FRONTEND=noninteractive apt install -y mysql-server redis-server golang git

# 设置MySQL密码和数据库
sudo mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'chensiyu';"
sudo mysql -e "CREATE DATABASE IF NOT EXISTS high_concurrency_db;"

# 启动服务
sudo systemctl restart mysql
sudo systemctl restart redis

# 克隆项目
git clone https://github.com/chensiyu1993/high-concurrency-api.git
cd high-concurrency-api

# 修改配置文件中的数据库密码
sed -i 's/your_password/chensiyu/g' config/config.yaml

# 运行项目
go mod tidy
go build -o api-server
./api-server

echo "=== 安装完成 ==="
echo "MySQL密码: chensiyu"
echo "服务已启动在: http://localhost:8080" 