# ss-wecom-assistant

[Chinese Version](./README.md)

## Introduction

ss-wecom-assistant is the most emotionally-rich chatbot implemented in Go, currently supporting deployment on the WeChat Official Account platform.

## Architecture

![architecture](./assest/ss-wecom-assistant-architecture.png)

## Implementation and features

- [x] Asynchronously reply to user information, not subject to the WeChat public platform's 15-second timeout cancellation.
- [x] Send text promotional information to users when they follow for the first time.
- [x] Contextual conversation feature, currently able to support a memory function of about 4K tokens.
- [ ] Asynchronously reply to graphic news.
- [ ] Q&A for professional knowledge base.
- [ ] Sensitive word detection.

## Requirements

- Verified service account
- OpenAI API Key
- Dokcer
- Go 1.18+

## Deployment

### Set up the configuration file

1. Set configuration variables: Choose the corresponding environment in conf/\* and fill in the configuration information.

   Below is an example of deploying in a local environment, please fill in the information corresponding to your service.

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
db_name = "ss-wecom-assistant"
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

### Install Databases

1.  Pull docker images

- Mysql: `docker pull mysql`
- Redis: `docker pull redis`

2.  Create data mount folders

- Create a new folder `/data/mysql/data` for local mounting of mysql data storage.
- Create a new folder `/data/mysql/log` for local mounting of mysql log data.

3.  Start docker images

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

## Compilation

    go build

## Deployment

By default, it is deployed locally. If it is online deployment, please remember to switch the environment. You can use -e to select the deployment environment.

    go run main.go -e local

## Give a Star! ⭐

If you like or are using this project to learn or start your solution, please give it a star. Thanks!

## License

ss-wecom-assistant is released under the MIT License. For more details, please see the LICENSE file.
