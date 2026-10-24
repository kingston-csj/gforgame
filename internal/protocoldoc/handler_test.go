package protocoldoc

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/forfun/gforgame/internal/protos"
	"github.com/forfun/gforgame/network"
	"github.com/gin-gonic/gin"
)

type testRoutes struct{}

func (*testRoutes) Login(s network.Session, msg *protos.ReqPlayerLogin) *protos.ResPlayerLogin {
	return &protos.ResPlayerLogin{}
}

func TestMetadataAndAssets(t *testing.T) {
	routes := network.NewMessageRoute()
	if err := routes.RegisterMessageHandlers(&testRoutes{}); err != nil {
		t.Fatal(err)
	}
	meta := Build(Connection{Codecs: []string{"json"}}, routes)
	var login, response Message
	for _, module := range meta.Modules {
		for _, message := range module.Protocols {
			switch message.Name {
			case "ReqPlayerLogin":
				login = message
			case "ResPlayerLogin":
				response = message
			}
			if message.Module == "cross" && message.Testable {
				t.Fatal("cross protocol exposed as testable")
			}
		}
	}
	if login.Cmd == nil || *login.Cmd != 103 || !login.Testable || login.Response != response.ID {
		t.Fatalf("invalid login metadata: %+v", login)
	}
	fields := meta.Schemas[login.Schema].Fields
	if len(fields) != 2 || fields[0].Name != "PlayerId" || fields[1].Name != "serverId" {
		t.Fatalf("invalid JSON fields: %+v", fields)
	}
	if fields[1].Desc != "服务器ID" {
		t.Fatalf("missing source comment: %+v", fields[1])
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router, meta)
	for _, file := range []string{"index.html", "protocol.css", "app.mjs", "protocol-core.mjs", "struct-codec.mjs", "meta.json"} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest("GET", "/protocol-doc/"+file, nil))
		if rec.Code != 200 || rec.Header().Get("Cache-Control") != "no-store" {
			t.Fatalf("%s: %d", file, rec.Code)
		}
		if file == "meta.json" {
			var got Metadata
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if len(got.Modules) == 0 || len(got.Schemas) == 0 {
				t.Fatal("empty metadata")
			}
		}
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("GET", "/protocol-doc/README.md", nil))
	if rec.Code != 404 {
		t.Fatal("unexpected asset exposed")
	}
}
