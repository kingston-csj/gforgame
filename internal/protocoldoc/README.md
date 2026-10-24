# 在线协议调试

启动 logic 后访问管理 HTTP 端口的 `/protocol-doc/index.html`，当前默认配置为
`http://localhost:8032/protocol-doc/index.html`。

`config/config-logic.yml` 的 `game.protocol-doc` 控制是否注册页面和元数据路由；
设为 `false` 可关闭。网关模式下将 `game.protocol-doc-options.websocket-url`
设为浏览器可访问的网关地址（例如 `ws://localhost:8021/ws`）。直连模式默认
使用页面所在主机、游戏端口和 `/ws`。HTTPS 页面需要可用的 `wss://` 入口。

页面和协议源码注释嵌入可执行文件，部署时无需复制资源目录。
`/protocol-doc/meta.json` 复用 `internal/tools/protocol/basegenerator.go` 的字段解析，
并结合已注册消息、实际路由处理器和 Go 类型生成，包含
模块、命令号、请求响应配对、嵌套字段、JSON 字段名和注释。
协议结构修改后重新编译并启动服务即可更新。

在线调试使用当前 Go 服务的 JSON 文本协议，直接发送 `ReqPlayerLogin`。
输入玩家 ID 登录后可筛选请求、编辑对象/集合/Map 参数、发送请求并查看
响应耗时和推送。收发记录最多保留 200 条。登录使用真实玩家会话。
JavaScript 中的 64 位整数以文本或 BigInt 保存，避免精度丢失。

配置示例（本地 `config/*.yml` 被 Git 忽略，需要在部署环境同步配置）：

```yaml
game:
  protocol-doc: true
  protocol-doc-options:
    websocket-url: ""
```

验证：`go test ./internal/tools/protocol ./internal/protocoldoc ./cmd/logic`；
`node --test internal/protocoldoc/protocol-core.test.mjs`。
