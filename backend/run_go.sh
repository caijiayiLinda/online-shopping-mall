#!/bin/bash
while true; do
    nohup go run main.go &
    wait $!
    echo "程序已退出，正在重新启动..."
done