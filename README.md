# ss-assistant

[英文版](./README.en.md)

## 介绍

ss-assistant 是一个基于 go 语言实现的被动回复聊天机器人, 目前仅支持部署认证的服务号。

## 架构

![架构](./assest/ss-assistant-architecture.png)

## 实现和特性

- [x] 异步回复用户信息，不受微信公众平台 15 秒超时取消限制。
- [x] 首次关注，给用户发送文字推广信息。
- [x] 上下文对话功能，目前能够支持约 4K token 数的记忆功能。
- [x] 异步回复图文消息。
- [ ] 专业领域知识库的问答。

## 准备

- 认证的服务号
- OpenAI API Key
- Dokcer
- go 1.18+

## 部署

### 设置配置文件

1. 设置配置变量: 在 `conf/*` 中选择对应的环境，填写配置信息。

   下面是一个在本地环境部署的例子，请填写自己的服务对应的信息。

```
[grpc]
host = "localhost"
port = 6333

[database]
host = "127.0.0.1"
port = 3306
user = "root"
password = "123456"
type = "mysql"
db_name = "ss-assistant"
init_conn = 20
max_conn = 200

[redis]
host = "127.0.0.1"
port = 6379
user = "root"
password = "123456"
init_conn = 30
max_conn = 120

[wechat]
app_id = "20211107"
app_secret = "20210117"
token = "****"
encoding_aes_key = "****"

[openai]
api_key = "this is secret"
```

### 安装数据库

1. 拉取 docker 镜像

- Mysql: `docker pull mysql`
- Redis: `docker pull redis`

2. 创建数据挂载文件夹

- 新建`/data/mysql/data`文件夹，用来在本地挂载 mysql 数据
- 新建`/data/mysql/log`文件夹，用来在本地挂载 mysql 日志数据
- 新建`/data/redis/data`文件夹，用来在本地挂载 redis 数据

3. 启动 docker 镜像

- mysql

```bash
docker run -d --restart=always --name mysql \
    -v /data/mysql/data:/var/lib/mysql \
    -v /data/mysql/log:/var/log/mysql \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=123456 \
    mysql \
    --character-set-server=utf8mb4 \
    --collation-server=utf8mb4_general_ci
```

- redis

```bash
docker run -p -d 6379:6379 --name redis --privileged=true \
    -v /data/redis/redis.conf:/etc/redis/redis.conf \
    -v /data/redis/data:/data \
    redis
```

```
docker run -d --name ss-assistant --log-driver=loki --log-opt loki-tag=ss-assistant -p 80:80 --env-file /root/github/ss-assistant/.env ss-assistant
```

### 编译

```
go build
```

### 部署

目前默认部署在本地 local 环境,如果是线上部署请记得切换环境。你可以通过 `-e` 来选择部署环境。

```
go run main.go -e local
```

## 给个星星！⭐

如果你喜欢或正在使用这个项目来学习或开始你的解决方案，请给它一个星星。谢谢！

## 许可证

ss-assistant 在 MIT 许可证下发布。更多细节请查阅 [LICENSE](./LICENSE) 文件。
