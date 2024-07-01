# gocron2

[![Downloads](https://img.shields.io/github/downloads/up9cloud/gocron2/total.svg)](https://github.com/up9cloud/gocron2/releases)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/up9cloud/gocron2/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/up9cloud/gocron2.svg?label=Release)](https://github.com/up9cloud/gocron2/releases)

定时任务管理系统

## 项目简介

使用Go语言开发的轻量级定时任务集中调度和管理系统, 用于替代 Linux-crontab

原有的延时任务拆分为独立项目[延迟队列](https://github.com/ouqiang/delay-queue)

原項目 [https://github.com/ouqiang/gocron](https://github.com/ouqiang/gocron) 的作者不更新了。[此為原文档](https://github.com/up9cloud/gocron2/wiki)

所以基於大家的 Folks 搞了新的 `gocron2` 縫合怪，詳細縫合訊息請看下方 Credits

## 功能特性

- Web界面管理定时任务
- crontab时间表达式, 精确到秒
- 任务执行失败可重试
- 任务执行超时, 强制结束
- 任务依赖配置, A任务完成后再执行B任务
- 账户权限控制
- 任务类型
  - shell任务
  > 在任务节点上执行shell命令, 支持任务同时在多个节点上运行
  - HTTP任务
  > 访问指定的URL地址, 由调度器直接执行, 不依赖任务节点
- 查看任务执行结果日志
- 任务执行结果通知, 支持邮件、Slack、Webhook

## 截图

![列表](https://user-images.githubusercontent.com/10205742/184531121-f5faa1a9-4d13-4132-a96d-848375765cda.jpg)
![日志](https://user-images.githubusercontent.com/10205742/184531126-0f159cda-8774-4185-9132-194e66cd5d3c.jpg)
![节点](https://user-images.githubusercontent.com/10205742/184531128-7a9a07a9-cac2-4dea-a37a-5cb57479a528.jpg)
![webhook](https://user-images.githubusercontent.com/10205742/184531159-582fd407-bed1-4ed4-a469-e8b9d5af67cb.jpg)

## 流程图

![流程图](https://raw.githubusercontent.com/up9cloud/gocron2/master/assets/screenshot/scheduler.png)

## 支持平台

> Windows、Linux、Mac OS

## 环境要求

> MySQL

## 安装

### 二进制安装

- 解压压缩包  
- `cd 解压目录`
- 启动
  - 调度器启动
    - Windows: `gocron2.exe web`
    - Linux、Mac OS:  `./gocron2 web`
  - 任务节点启动, 默认监听 0.0.0.0:5921
    - Windows:  `gocron2-node.exe`
    - Linux、Mac OS:  `./gocron2-node`
- 浏览器访问 `http://localhost:5920`

### 源码安装

- 安装Go 1.12+
- `go get -d github.com/up9cloud/gocron2`
- `export GO111MODULE=on`
- 编译 `make`
- 启动
  - gocron2 `./bin/gocron2 web`
  - gocron2-node `./bin/gocron2-node`

### docker

See: `https://github.com/docker-sstc/docker-gocron2`

配置: /app/conf/app.ini

日志: /app/log/cron.log

## 命令

- gocron2
  - -v 查看版本
- gocron2 web
  - --host 默认0.0.0.0
  - -p 端口, 指定端口, 默认5920
  - -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
  - -h 查看帮助
- gocron2-node
  - -allow-root *nix平台允许以root用户运行
  - -s ip:port 监听地址
  - -enable-tls 开启TLS
  - -ca-file   CA证书文件
  - -cert-file 证书文件
  - -key-file  私钥文件
  - -h 查看帮助
  - -v 查看版本

## 开发

- 安装 Go1.12+, Node.js (Npm), Docker
- 启动 gocron2, gocron2-node `docker compose up`
- 启动前端, `make run-vue`, API请求会转发给 gocron2
  - install 的時候，主机名填入 `db`

`其他指令請詳查 ./makefile`

## TODO

- [x] 版本升级
- [x] 批量开启、关闭、删除任务
- [x] 调度器与任务节点通信支持https
- [x] 任务分组
- [x] 多用户
- [x] 权限控制

## 程序使用的组件

- Web框架 [Macaron](http://go-macaron.com/)
- 定时任务调度 [Cron](https://github.com/robfig/cron)
- ORM [Xorm](https://github.com/go-xorm/xorm)
- UI框架 [Element Plus](https://element-plus.org)
- 依赖管理 [Govendor](https://github.com/kardianos/govendor)
- RPC框架 [gRPC](https://github.com/grpc/grpc)

## Credits

- `https://github.com/gaowei-space/gocron`
- `https://github.com/ghostong/gocron.git`
