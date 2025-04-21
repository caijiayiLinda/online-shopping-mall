#!/bin/bash
export JWT_SECRET="secret-key" # TODO: Replace with a secure secret
while true; do
    nohup go run main.go &
    wait $!
    echo "程序已退出，正在重新启动..."
done
