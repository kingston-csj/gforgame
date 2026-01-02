package http

import (
	"io/github/gforgame/examples/context"

	"github.com/gin-gonic/gin"
)

// 关闭服务器
func StopServer(c *gin.Context) {
	context.TcpServer.Running <- true
}
