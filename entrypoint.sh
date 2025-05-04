#!/bin/bash

# 等待依赖服务就绪
/app/wait-for.sh mysql:3307 --timeout=60 -- echo "MySQL is up!"
/app/wait-for.sh redis:6379 --timeout=60 -- echo "Redis is up!"
/app/wait-for.sh elasticsearch:9200 --timeout=60 -- echo "Elasticsearch is up!"
/app/wait-for.sh rabbitmq:5672 --timeout=60 -- echo "RabbitMQ is up!"

# 执行初始化操作
#./main -sql   # 初始化数据库
#./main -es    # 初始化ES索引
#./main -admin # 创建管理员

# 启动主程序
exec ./main
