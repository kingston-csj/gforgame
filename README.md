
## 项目介绍　　
gforgame，jforgame的go语言实现。是一个轻量级高性能手游服务端框架。项目提供各种支持快速二次开发的组件，以及常用业务功能作为演示。
需要说明的是，语言只是工具，重要的是思想。 当然，由于程序语言本身语法差异，也会影响思考方式和编码方式。 

## 项目特点
* 搭配框架博客栏目教程，快速理解项目模块原理
* 支持socket/webSocket接入，完美适配手游/页游/H5/小游戏服务端架构
* 通信协议支持protobuf或json，为客户端提供多种选择
* 有独立http管理后台网站，为游戏运维/运营提供支持  --> [后台管理系统](https://github.com/kingston-csj/gamekeeper)  

## 快速入门
下载代码到本地，导入项目到vscode或者goland开发工具
服务器入口： main.go
客户端入口： client.go

## 已实现功能
* tcp网关，消息路由，消息分发链  
* 日志模块
* 事件驱动
* 玩家数据读写
* 通信协议支持protobuf

## 近期功能
* websocket接入
* csv配置文件读取, jforgame-data实现
* 数据缓存与异步持久化
