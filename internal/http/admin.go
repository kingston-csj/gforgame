package http

import (
	"github.com/forfun/gforgame/internal/context"

	"github.com/gin-gonic/gin"
)

// 关闭服务器
func StopServer(c *gin.Context) {
	if context.GameServer != nil {
		context.GameServer.NotifyStop()
	}
}
