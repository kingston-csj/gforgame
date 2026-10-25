package protocoldoc

import (
	"embed"
	"encoding/json"
	"net/http"

	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/network"
	"github.com/gin-gonic/gin"
)

//go:embed assets/*
var assets embed.FS

// Register adds the optional protocol documentation to the admin HTTP server.
func Register(router *gin.Engine, routes *network.MessageRoute) {
	enabled, _ := serverconfig.GetExtraBool("game.protocol-doc")
	if !enabled {
		return
	}
	// node := serverconfig.GetCurrentNode()
	// serviceAddr := node.ServiceAddr
	// // 如果使用网关的话,取网关的服务地址
	// if node.UseGateMode {
	// 	gateNodes := serverconfig.GetServersByType(uint32(serverconfig.NodeTypeGateway))
	// 	if (len(gateNodes) <= 0) {
	// 		panic("gate node not found")
	// 	}
	// 	gate := gateNodes[0]
	// 	serviceAddr = gate.Addr
	// }
	// _, port, _ := net.SplitHostPort(serviceAddr)
	// connection := Connection{ServerID: node.ServerId, Port: port, Path: "/ws", Codecs: []string{"json"}}
	// RegisterRoutes(router, Build(connection, routes))
}

// RegisterRoutes serves an immutable metadata snapshot and the bundled frontend.
func RegisterRoutes(router *gin.Engine, meta Metadata) {
	content, err := json.Marshal(meta)
	if err != nil {
		panic(err)
	}
	router.GET("/protocol-doc", func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/protocol-doc/index.html") })
	router.GET("/protocol-doc/", func(c *gin.Context) { c.Redirect(http.StatusTemporaryRedirect, "/protocol-doc/index.html") })
	router.GET("/protocol-doc/:file", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		file := c.Param("file")
		if file == "meta.json" {
			c.Data(http.StatusOK, "application/json; charset=utf-8", content)
			return
		}
		mime := map[string]string{"index.html": "text/html", "protocol.css": "text/css", "app.mjs": "application/javascript", "protocol-core.mjs": "application/javascript", "struct-codec.mjs": "application/javascript"}[file]
		if mime == "" {
			c.Status(http.StatusNotFound)
			return
		}
		data, err := assets.ReadFile("assets/" + file)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, mime+"; charset=utf-8", data)
	})
}
