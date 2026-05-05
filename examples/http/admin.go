package http

import (
	"github.com/forfun/gforgame/examples/context"

	"github.com/gin-gonic/gin"
)

// 关闭服务器
func StopServer(c *gin.Context) {
	context.TcpServer.Running <- true
}
