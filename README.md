# Develop Applet

做个小程序后端服务

## 简介

「做个小程序」 是一款致力于利用AI能力快速创造工具的 Web 平台(https://ai.moyu.chat)
- 此项目为社区共建基于AI的 ```做个小程序``` 项目，它基于 ```摸鱼星球``` 改造而来，故同时包含 <摸鱼星球> 的功能和 API 文档，并共享部分功能
- 此项目基于自研的开源框架 ```Gone``` 打造 
 https://github.com/gone-io/gone

## API

  [做个小程序 API 文档](docs/Program-Developer-API.md)

## 本地运行

### 1.启动依赖服务

- ***确保已正确安装 `docker` 与 `docker-compose`***
- ***确保本地没有占用 `docker-compose.yml` 声明的端口***

```sh
$ docker-compose up
```

### 2.运行项目

```sh
$ make run
```

也可以

```sh
$ make gone
```

生成依赖注入代码然后运行 `cmd/server/main.go`

### 3.注意

由于项目本身不包含部分模块所需的敏感信息文件，所以本地运行的状态下以下功能不可用

- 微信支付相关功能
- 微信登录相关功能
- 阿里云 OSS 相关功能

涉及以上内容的开发应当本地使用单测打桩跑通业务后在 `dev` 环境进行进一步联调

## 二、项目结构

```sh
.
├── asserts   其他资源文件
├── cmd    程序入口
├── config    配置文件目录
├── docs    文档目录
├── internal   
│     ├── controller    Controller
│     ├── interface    
│     │     ├── entity    实体对象定义
│     │     └── service   模块接口定义
│     ├── middleware  gin 中间件
│     ├── module    模块实现
│     ├── pkg    通用工具代码
│     └── router    路由
├── k8s    k8s 配置目录
├── scripts   SQL更新脚本目录 
```shell
