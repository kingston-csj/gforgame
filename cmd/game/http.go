package main

import (
	"time"

	"github.com/forfun/gforgame/common/util/jsonutil"
	playerrepo "github.com/forfun/gforgame/internal/infra/repository/player"
	"github.com/forfun/gforgame/internal/protocoldoc"
	serverconfig "github.com/forfun/gforgame/internal/serverconfig"
	"github.com/forfun/gforgame/internal/service/recharge"
	"github.com/forfun/gforgame/network"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewHttpServer 装配管理后台 HTTP 路由并启动监听（阻塞）。
func NewHttpServer(
	messageRoutes *network.MessageRoute,
	playerRepo *playerrepo.PlayerRepository,
	rechargeService *recharge.RechargeService,
	profileService *playerrepo.PlayerProfileService,
) {
	router := gin.Default()

	// CORS 必须在注册路由之前 Use，否则对已注册路由不生效
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	protocoldoc.Register(router, messageRoutes)
	if err := router.Run(serverconfig.GetCurrentNode().HttpAddr); err != nil {
		panic(err)
	}
}

// AdminHttpResponse 管理后台统一响应结构。
type AdminHttpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Ok 包装成功响应。
// 注意：Data 当前经 jsonutil.StructToJSON 二次序列化为字符串，与 gin 自动序列化重复，
// 若前端按字符串解析则保持现状；若前端期望对象，可改为 Data: data。
func Ok(data any) *AdminHttpResponse {
	jsonStr, err := jsonutil.StructToJSON(data)
	if err != nil {
		return &AdminHttpResponse{
			Code: 500,
			Msg:  "json serialize failed: " + err.Error(),
		}
	}
	return &AdminHttpResponse{
		Code: 200,
		Msg:  "success",
		Data: jsonStr,
	}
}
