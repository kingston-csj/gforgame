package http

import (
	serverpkg "github.com/forfun/gforgame/network/server"
	"github.com/gin-gonic/gin"
)

// 关闭服务器
func StopServer(c *gin.Context, gameServer serverpkg.Server) {
	if gameServer != nil {
		gameServer.NotifyStop()
	}
}
