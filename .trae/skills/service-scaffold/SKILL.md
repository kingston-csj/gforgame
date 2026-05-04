---
name: service-scaffold
description: Generates Go service, route, and protos module skeletons based on vipservice.go/vip.go. Invoke when user asks to add a new module or XXService boilerplate.
---

# Service Scaffold

用于在项目中快速新增符合现有风格的 Go Service + Route + Protos 文件骨架，参考 `examples/service/vip/vipservice.go` 与 `examples/route/vip.go` 的结构与命名。

## 何时使用

- 用户提出“新增一个 XXService”
- 用户提出“生成 service 模板/骨架”
- 用户希望按现有项目规范快速创建新服务文件

## 目标输出

在以下三个位置生成骨架：

1. `examples/service/<servicefolder>/<servicefile>.go`
2. `examples/route/<routefile>.go`
3. `protos/<protofile>.go`

### Service 文件

1. `package` 与目录名一致
2. 标准 import 区（先保留基础依赖，业务依赖按需补充）
3. `type <Name>Service struct{}`
4. 单例变量：
   - `instance *<Name>Service`
   - `once sync.Once`
5. `Get<Name>Service() *<Name>Service`，使用 `once.Do(...)` 初始化

### Route 文件

1. `type <Name>Route struct`
   - 嵌入 `network.Base`
   - 持有 `service *<pkg>.<Name>Service`
2. `New<Name>Route() *<Name>Route`
3. `Init()` 中初始化 service：`r.service = <pkg>.Get<Name>Service()`

### Protos 文件

1. `package protos`
2. 文件名与模块名一致（例如 `equip.go`、`vip.go`）
3. 若文件不存在则创建占位骨架，后续消息结构体由业务再补充

## 命名规则

- 用户给定 `XXService` 时：
  - 结构体：`XXService`
  - 获取函数：`GetXXService`
  - Route 结构体：`XXRoute`
  - Route 构造函数：`NewXXRoute`
  - 包名：小写业务名（例如 `vip`, `mail`, `arena`）
  - 文件名：`xxservice.go`（全小写）
  - Route 文件名：`xx.go`（放在 `examples/route`）
  - Protos 文件名：`xx.go`（放在 `protos`）
- 如果用户只给业务名（如 `vip`）：
  - 自动推导结构体为 `VipService`
  - 自动推导函数为 `GetVipService`
  - 自动推导 Route 为 `VipRoute`
  - 文件名为 `vipservice.go`
  - Route 文件名为 `vip.go`
  - Protos 文件名为 `vip.go`

## 生成模板

### Service 模板

```go
package <packageName>

import (
	"sync"
  configdomain "io/github/gforgame/examples/domain/config"
	playerdomain "io/github/gforgame/examples/domain/player"
)

type <ServiceName> struct{}

var (
	instance *<ServiceName>
	once     sync.Once
)

func Get<ServiceName>() *<ServiceName> {
	once.Do(func() {
		instance = &<ServiceName>{}
	})
	return instance
}
```

### Route 模板

```go
package route

import (
	"io/github/gforgame/examples/service/<packageName>"
	"io/github/gforgame/network"
)

type <RouteName> struct {
	network.Base
	service *<packageName>.<ServiceName>
}

func New<RouteName>() *<RouteName> {
	return &<RouteName>{}
}

func (r *<RouteName>) Init() {

}
```

### Protos 模板

```go
package protos
```

## 执行步骤

1. 先解析用户输入中的服务名（`XXService` 或业务名）。
2. 推导包名、service 文件名、route 文件名、protos 文件名、结构体名、Get/New 函数名。
3. 若目标文件已存在，优先读取并在保持风格前提下最小改动；若不存在则创建新文件。
4. 先生成 service 文件，再生成 route 文件，再生成 `protos/<module>.go`。
5. 写入模板代码，并保持 gofmt 风格（制表符缩进）。
6. 在main.go文件中注册新的路由，例如：route.NewVipRoute(), route.NewQiRiRoute()等。
7. 返回变更路径和关键代码片段。

## 示例触发

- “新增一个 MailService”
- “帮我加一个 arena service”
- “按 vipservice 的结构新建一个背包服务”
