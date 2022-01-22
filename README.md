# remix-api

本项目为[remix-mini-program](https://github.com/trumanwong/remix-mini-program.git)的后端项目，基于`go`开发。

## 安装指引

安装前，你的服务器需要先安装`docker`、`docker-compose`、`ffmpeg`、`python3`、`supervisor`、`rabbitmq`、`elasticsearch`(可选)。

```shell
$ cd /data/docker/golang
$ git clone git@github.com:trumanwong/remix-api.git
$ cd remix-api
$ cp config.yaml.example config.yaml # 修改你自己相应的配置
$ cd scripts
$ docker-compose up -d
# 进入容器
$ docker exec -it remix-api
$ cd ../cmd/server
$ go build -o server
$ cd ../consumers/task
$ go build -o task
$ cd ../log
$ go build -o log
$ cd ../console
$ go build -o console
$ supervisorctl restart all
# 退出容器
$ exit
$ cp /data/docker/golang/remix-api/console.conf /etc/supervisor/conf.d/console.conf
$ cp /data/docker/golang/remix-api/remix-task.conf /etc/supervisor/conf.d/remix-task.conf
$ supervisorctl restart all
```

## 体验地址
[Web Demo](https://www.trumanwl.com/tools/remix)

微信小程序（鬼畜动图生成器）

![鬼畜动图生成器](./remix-app.png)