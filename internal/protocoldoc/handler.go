package protocoldoc

import (
	"embed"
	"encoding/json"
	"net"
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
	externalURL, _ := serverconfig.GetExtraString("game.protocol-doc-options.websocket-url")
	_, port, _ := net.SplitHostPort(serverconfig.ServerConfig.ServerUrl)
	connection := Connection{URL: externalURL, Port: port, Path: "/ws", ServerID: serverconfig.ServerConfig.ServerId, TestEnabled: !serverconfig.ServerConfig.UseGateMode || externalURL != "", Codecs: []string{"json"}}
	RegisterRoutes(router, Build(connection, routes))
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
