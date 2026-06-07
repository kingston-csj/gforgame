package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	serverconfig "github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewHttpServer() *gin.Engine {
	router := gin.Default()
	// 关闭游戏服务器进
	router.POST("/api/stop", func(c *gin.Context) {
		http.StopServer(c)
	})
	// 清理数据库
	router.POST("/api/clearDb", clearDb)
	// 配置 CORS 中间
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有源，生产环境应指定具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	err := router.Run(serverconfig.ServerConfig.HttpUrl)
	if err != nil {
		panic(err)
	}

	return router
}

func clearDb(c *gin.Context) {
	dbURL, ok := serverconfig.GetExtraString("db.url")
	if !ok || strings.TrimSpace(dbURL) == "" {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "执行失败: 未配置 db.url",
		})
		return
	}

	username, password, database, err := parseMySQLDSN(dbURL)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "执行失败: " + err.Error(),
		})
		return
	}

	shellPath, err := resolveShellPath("shell", "cleardb.sh")
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "执行失败: " + err.Error(),
		})
		return
	}

	output, err := executeLocalScript(shellPath, username, password, database)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 500,
			"msg":  "执行失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": strings.TrimSpace(output),
	})
}

func parseMySQLDSN(dbURL string) (string, string, string, error) {
	parts := strings.SplitN(dbURL, "@", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("db.url 格式不正确")
	}

	userInfo := strings.SplitN(parts[0], ":", 2)
	if len(userInfo) != 2 {
		return "", "", "", fmt.Errorf("db.url 缺少用户名或密码")
	}

	dbPart := parts[1]
	slashIndex := strings.LastIndex(dbPart, "/")
	if slashIndex < 0 || slashIndex == len(dbPart)-1 {
		return "", "", "", fmt.Errorf("db.url 缺少数据库名")
	}

	databasePart := dbPart[slashIndex+1:]
	if queryIndex := strings.Index(databasePart, "?"); queryIndex >= 0 {
		databasePart = databasePart[:queryIndex]
	}
	databasePart, err := url.QueryUnescape(strings.TrimSpace(databasePart))
	if err != nil {
		return "", "", "", fmt.Errorf("解析数据库名失败: %w", err)
	}
	if databasePart == "" {
		return "", "", "", fmt.Errorf("db.url 缺少数据库名")
	}

	return strings.TrimSpace(userInfo[0]), strings.TrimSpace(userInfo[1]), databasePart, nil
}

func resolveShellPath(parts ...string) (string, error) {
	relativePath := filepath.Join(parts...)
	searchDirs := make([]string, 0, 4)

	if exePath, err := os.Executable(); err == nil {
		searchDirs = append(searchDirs, filepath.Dir(exePath))
	}
	if cwd, err := os.Getwd(); err == nil {
		searchDirs = append(searchDirs, cwd)
	}

	seen := make(map[string]struct{}, len(searchDirs))
	for _, baseDir := range searchDirs {
		dir := filepath.Clean(baseDir)
		for {
			if _, ok := seen[dir]; ok {
				break
			}
			seen[dir] = struct{}{}

			candidate := filepath.Join(dir, relativePath)
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}

			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	return "", fmt.Errorf("未找到脚本 %s", relativePath)
}

func executeLocalScript(shellPath string, args ...string) (string, error) {
	shellCmd, err := resolveShellExecutor()
	if err != nil {
		return "", err
	}

	cmdArgs := append([]string{shellPath}, args...)
	cmd := exec.Command(shellCmd, cmdArgs...)
	cmd.Dir = filepath.Dir(shellPath)
	output, err := cmd.CombinedOutput()
	result := string(output)
	if err != nil {
		if strings.TrimSpace(result) != "" {
			return "", fmt.Errorf("%v, 输出: %s", err, strings.TrimSpace(result))
		}
		return "", err
	}
	return result, nil
}

func resolveShellExecutor() (string, error) {
	candidates := []string{"sh", "bash"}
	if runtime.GOOS == "windows" {
		candidates = []string{"bash", "sh"}
	}

	for _, candidate := range candidates {
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("未找到可用的 shell 解释器")
}
