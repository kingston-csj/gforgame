---
name: "route-req-autocomplete"
description: "根据 internal/protos/XX.go 自动生成 route 的 Req* 处理方法。在用户说“新增XX路由”或补全 internal/route 下 Req* 方法时触发。"
---

# Route Req 自动补全

使用本技能可根据 `internal/protos/XX.go` 中的 `Req*` 结构体，在 `internal/route` 下自动生成路由处理方法。

## 触发时机

- 用户输入 `新增XX路由` 或 `补充XX路由`。
- 用户要求为某个路由模块补齐全部 `Req` 处理方法。
- 你正在 `internal/route` 中补全 `Req*` 方法。

## 指令意图

输入示例：

- `新增skin路由`
- `新增equip路由`

期望行为：

1. 定位 `internal/protos/XX.go`。
2. 找出文件内所有 `Req*` 结构体。
3. 为每个 `ReqXxx` 推导对应返回 `ResXxx`。
4. 打开或创建 `internal/route/XX.go`。
5. 按项目现有风格生成缺失路由方法。
6. 不重复生成已存在的方法。

## 目标方法模板

参考风格（`signin.go`）：

```go
func (ps *SignInRoute) ReqSignIn(playerId string, index int32, msg *protos.ReqSignIn) *protos.ResSignIn {
	player := player.GetPlayerService().GetPlayer(playerId)
	err := ps.service.SignIn(player)
	if err != nil {
		return &protos.ResSignIn{
			Code: int32(err.Code()),
		}
	}
	return &protos.ResSignIn{}
}
```

生成风格（`ReqXxx` -> `ResXxx`）：

```go
func (r *XxRoute) ReqXxx(playerId string, index int32, msg *protos.ReqXxx) *protos.ResXxx {
	p := player.GetPlayerService().GetPlayer(playerId)
	code := r.service.Xxx(p, msg)
	return &protos.ResXxx{
		Code: code,
	}
}
```

## 生成规则（严格）

1. 方法名必须以 `Req` 开头。
2. 方法签名应使用：
   - `playerId string`
   - `index int32`
   - `msg *protos.ReqXxx`
   - 返回 `*protos.ResXxx`
3. 先通过 playerId 获取玩家：
   - `p := player.GetPlayerService().GetPlayer(playerId)`
4. 按模块既有命名约定调用 `r.service`。
5. 返回值映射优先级：
   - 若 service 返回 `int32`：填充 `Code`。
   - 若 service 返回 `*protos.ResXxx`：直接返回。
   - 若 service 返回 `BusinessError`：映射到 `Code`。
6. 代码风格需与模块现有文件一致（如 `equip.go`、`signin.go`、`rune.go`）。
7. import 保持最小且正确（`network`、`protos`、模块 service、player service，以及可选 event/context import）。

## 文件级行为

- 若 `internal/route/XX.go` 不存在：
  - 创建 `XXRoute`、`NewXXRoute`、`Init()` 及 `service` 字段。
- 若文件已存在：
  - 保留当前结构，仅追加缺失的 `Req*` 方法。
- 输出代码需可直接通过格式化工具处理。

## 输出检查清单

- `internal/protos/XX.go` 中所有 `Req*` 均已覆盖到路由方法。
- 不重复生成已存在的方法。
- 在当前模块 service/protos 合同下尽可能保证可编译。
