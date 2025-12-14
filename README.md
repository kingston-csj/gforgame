## 项目介绍　　

gforgame，jforgame 的 go 语言实现。是一个轻量级高性能手游服务端框架。项目提供各种支持快速二次开发的组件，以及常用业务功能作为演示。
需要说明的是，语言只是工具，重要的是思想。 当然，由于程序语言本身语法差异，也会影响思考方式和编码方式。

## 项目特点

- 搭配框架博客栏目教程，快速理解项目模块原理和 Go 语言特性。
- 最适合 java 转 go 的开发人员，在这里，你可以找到熟悉的感觉。
- 融合 Java 游戏服务器思想，借鉴 Java 游戏服务器的设计理念，结合 Go 语言的特性进行优化和创新，充分发挥两者的优势。
- 灵活兼容 Socket 和 WebSocket 接入，通过配置参数轻松切换，无需修改业务代码，完美适配手游、页游、H5 和小游戏的服务端架构。
- 提供 protobuf /JSON/自动解析 struct 几种通信协议，满足不同客户端的需求，灵活选择最适合的协议进行通信。
- 提供独立的 HTTP 管理后台网站，支持游戏运维和运营的各类管理需求，提升运维效率和运营支持能力。 --> [后台管理系统](https://github.com/kingston-csj/gamekeeper)
- 可视化 cocos 游戏客户端，功能演示
- 可视化 unity 游戏客户端，功能演示
- 一键导出客户端协议（typescript/csharp）
- 玩家可自由选择 Cocos 或 Unity 端进入游戏，双端同步运行、实时互联，打造无缝跨平台交互玩法。
- cocos/unity 客户端，亦可作为 jforgame 的游戏客户端

## 快速入门

### 代码导入

下载代码到本地，导入项目到 vscode 或者 goland 开发工具  
下载 cocos creator3.8 编辑器
服务器入口： main.go

cocos 客户端工程(推荐)：client/cocos
unity 客户端工程：client/unity
go 客户端入口：client/go/client.go

### 私有协议栈

包括包头及包体，格式如下  
// header(12bytes) | body  
// cmd | index| len(body) | body

### 消息编解码

在 message.go 统一管理所有协议的 id
文本类型，可以使用 json 格式。  
二进制类型，可以使用基于 struct 反射（目前只支持 unity 客户端），或者 protobuf 格式。  
如果选择 json 或 struct 的方式，直接在 protos 目录下定义消息格式即可。  
如果采用 protobuf 的话，则通过 message.proto 进行注册，再通过 gen_proto.bat 脚本生成 message.pb.go 文件  
（使用 protobuf，需要下载 protobuf 编译工具，以及 go 插件，并将两者添加到系统环境变量）
新手用户推送使用 json 格式，cocos/unity 客户端默认支持。

### 消息路由

遵循“约定大于配置”思想，根据方法（非函数）签名自动扫描，满足一定格式的方法即为消息处理器

```golang
    // 消息处理器格式1： 方法第一个参数要求是session,第二个参数要求是已注册的消息; 若方法有返回值且不为空，则自动将返回值下发给客户端
    func (rs PlayerController) ReqLogin(s *network.Session, msg *protos.ReqPlayerLogin) interface{} {

    }
```

```golang
    // 消息处理器格式2： 方法第一个参数要求是session,第二个参数是一个index; 第三个参数要求是已注册的消息
    // 索引用于异步给客户端发送请求(例如另起协程)，如果是同步的话，直接通过格式1即可
    func (rs PlayerController) ReqLogin(s *network.Session, index int, msg *protos.ReqPlayerLogin)  {

    }
```

### 玩家数据读写

数据库使用 mysql, orm 使用 gorm  
当有数据发生变化时，定时全量更新，结合 cache 机制，提高读写性能

### 功能模块

每个功能以模块的形式组织业务，例如背包，任务，技能等等  
模块需继承 Module，并在 Init()方法注册该模块的所有通信协议  
新模块要通过 network.RegisterModule(player.NewPlayerController())进行注册（扫描消息路由）

### websocket

node.Startup()方法参数增加 network.WithWebsocket()代表选择 websocket  
example/cocos 为 ws 的客户端程序

### 跨服通信方式一：基于 rpc

使用 grpc 进行跨进程通信，需要先安装 protobuf 和 protoc-gen-go-grpc 编译插件  
根据不同的跨服类型，自行定义业务逻辑，参考 examples/cross 相关代码（流程待完善）

### 跨服通信方式二：基于 socket

推荐使用原生的 socket 进行跨进程通信，无论是游戏应用作为客户端，还是跨服节点作为客户端  
都使用统一的调用方式，也无须引入 grpc  
客户端消息处理支持三种方法：  
1，类似服务器消息路由（通过方法签名）  
2，客户端同步阻塞(通过方法返回值)

```golang

    import "io/github/gforgame/network/client"
    r, err := client.Request(session, &protos.ReqPlayerLogin{Id: "1001"})
	if err != nil {
		fmt.Println(err)
	}
	resPlayerLogin := r.(*protos.ResPlayerLogin)
	fmt.Println("客户端收到消息：(", resPlayerLogin, ")")
```

3，客户端异步回调(通过注册回调函数，类似 ajax)

```golang

    import "io/github/gforgame/network/client"

    // 实现 RequestCallback 接口的匿名对象
    type commonCallback struct{}

    func (rc *commonCallback) OnSuccess(result any) {
        fmt.Println("OnSuccess: ", result)
    }

    func (rc *commonCallback) OnError(err error) {
        fmt.Println("OnError: ", err)
    }

    client.Callback(session, &protos.ReqPlayerLogin{Id: "1001"}, &commonCallback{})
```

## 框架功能

- socket/websocket 网关，消息路由，消息分发链
- 客户端支持三种消息处理：异步回调/同步请求/消息路由
- 日志模块
- 多环境配置
- 玩家数据读写
- 通信协议支持 json+protobuf+struct 反射
- 使用原生 map 实现一套高效 cache 工具，直接存储原生对象引用而非 byte[]，避免频繁序列化与反序列化
- 数据异步持久化，玩家数据实时更新缓存，定时持久化到数据库
- http 管理后台
- grpc 跨服通信接入
- excel 配置文件读取, jforgame-data 实现

## 业务功能(cocos 交互界面已完成，unity 交互界面逐步完善)

同时按下键盘 alt+g+m 组合键可切换 gm 控制台

- 账号登录
- 英雄招募、升级、突破、合成、图鉴、布阵
- 背包
- 邮件
- 排行榜
- 红点系统

## 部分客户端运行效果

登录界面  
![](/screenshots/login.png "登录界面")

主界面  
![](/screenshots/main.jpg "主界面")

背包界面  
![](/screenshots/bag.jpg "背包界面")

招募界面  
![](/screenshots/recruit.jpg "招募界面")

英雄主界面  
![](/screenshots/hero.jpg "英雄主界面")

图鉴界面  
![](/screenshots/tujian.jpg "图鉴界面")

邮件界面  
![](/screenshots/mail.jpg "邮件界面")

排行榜界面  
![](/screenshots/rank.jpg "排行榜界面")
