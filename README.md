# go-blog后端项目 :blue_book:

## 技术栈
Gin + Gorm + Mysql + Redis + Elasticsearch + RabbitMQ

## 项目描述
用户博客论坛，采用Gin实现用户注册，登录和发文功能，支持阅读、收藏、评论和文章搜索
1. 基于JWT实现用户权限校验，通过Redis维护Token黑名单限制多地点登录同一账号
2. 使用Elasticsearch存储文章，支持关键词检索和内容检索，可以按照浏览量、收藏量等排序
3. 利用RabbitMQ异步更新Elasticsearch中文章的浏览量、评论量等信息,采用批量更新
4. 使用Cron设置定时任务，爬取日历和热搜新闻信息并存储到Redis中
5. 基于Docker Compose​​实现一键式容器化部署


## 部署
需要安装docker和docker compose
### 后端
修改`config.yaml`配置文件和`docker-compose.yaml`文件
```shell
docker compose build
docker compose up
```