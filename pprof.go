package main

import (
	"net/http"
	"net/http/pprof"
	"strings"
)

// allowedIPs 是允许访问 pprof 的 IP 列表
var allowedIPs = []string{"127.0.0.1", "::1", "192.168.1.100"}

// ipIsAllowed 检查请求的 IP 是否在允许列表中
func ipIsAllowed(ip string) bool {
	for _, allowedIP := range allowedIPs {
		if ip == allowedIP {
			return true
		}
	}
	return false
}

// wrapPprofHandler 包装默认的 pprof 处理函数，添加 IP 检查
func wrapPprofHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if !ipIsAllowed(ip) {
			http.Error(w, "not allowed", http.StatusForbidden)
			return
		}
		handler(w, r)
	}
}

func NewHttpServeMux() *http.ServeMux {
	// 创建默认的 pprof 处理器
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", wrapPprofHandler(http.HandlerFunc(pprof.Index)))
	mux.HandleFunc("/debug/pprof/cmdline", wrapPprofHandler(http.HandlerFunc(pprof.Cmdline)))
	mux.HandleFunc("/debug/pprof/profile", wrapPprofHandler(http.HandlerFunc(pprof.Profile)))
	mux.HandleFunc("/debug/pprof/symbol", wrapPprofHandler(http.HandlerFunc(pprof.Symbol)))
	mux.HandleFunc("/debug/pprof/trace", wrapPprofHandler(http.HandlerFunc(pprof.Trace)))

	return mux
}
