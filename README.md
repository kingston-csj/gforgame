
## 项目介绍　　
gforgame，jforgame的go语言实现。是一个轻量级高性能手游服务端框架。项目提供各种支持快速二次开发的组件，以及常用业务功能作为演示。
需要说明的是，语言只是工具，重要的是思想。 当然，由于程序语言本身语法差异，也会影响思考方式和编码方式。 

## 项目特点
* 搭配框架博客栏目教程，快速理解项目模块原理和Go语言特性。
* 融合Java游戏服务器思想: 借鉴Java游戏服务器的设计理念，结合Go语言的特性，优化和创新，充分发挥两者的优势。
* 灵活接入，兼容Socket和WebSocket接入，通过配置参数轻松切换，无需修改业务代码，完美适配手游、页游、H5和小游戏的服务端架构。
* 提供protobuf和JSON两种通信协议，满足不同客户端的需求，灵活选择最适合的协议进行通信。
* 提供独立的HTTP管理后台网站，支持游戏运维和运营的各类管理需求，提升运维效率和运营支持能力。  --> [后台管理系统](https://github.com/kingston-csj/gamekeeper)  

## 快速入门
### 代码导入
下载代码到本地，导入项目到vscode或者goland开发工具
项目自带多个模块案例代码，如player_service,chat_service
服务器入口： main.go  (参数增加network.WithWebsocket()代表选择websocket，默认为tcpsocket)
客户端入口： client.go

### 私有协议栈
包括包头及包体，格式如下
//      header(8bytes)     | body
// msgLength = 8+len(body) | body
//  cmd | msgLength        | body

### 消息编解码
在message.go统一管理所有协议的id
如果采用json的话，可直接在message.go定义协议结构
如果采用protobuf的话，则通过message.proto进行注册，再通过gen_proto.bat脚本生成message.pb.go文件
（使用protobuf，需要下载protobuf编译工具，以及go插件，并将两者添加到系统环境变量）

### 消息路由
遵循“约定大于配置”思想，根据方法签名自动扫描，满足一定格式的方法即为消息处理器，
```golang
    // 消息处理器格式： 第一个参数要求是session,第二个参数要求是已注册的消息; 若方法有返回值且不为空，则自动将返回值下发给客户端
    func (rs PlayerService) ReqLogin(s *network.Session, msg *protos.ReqPlayerLogin) interface{} {
    
    }
```

### 玩家数据读写
数据库使用mysql, orm使用gorm，
当有数据发生变化时，定时全量更新，结合cache机制(待实现)，提高读写性能

### 功能模块
每个功能以模块的形式组织业务，例如背包，任务，技能等等
模块需继承Module，并在Init()方法注册该模块的所有通信协议
新模块要通过network.RegisterModule(player.NewPlayerService())进行注册（扫描消息路由）

## 已实现功能
* tcp网关，消息路由，消息分发链  
* 日志模块
* 事件驱动
* 玩家数据读写
* 通信协议支持json+protobuf
* websocket接入

## 近期功能
* csv配置文件读取, jforgame-data实现
* 数据缓存与异步持久化
* grpc接入
