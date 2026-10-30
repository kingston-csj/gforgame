package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/forfun/gforgame/internal/logger"
	serverconfig "github.com/forfun/gforgame/internal/serverconfig"
	"github.com/gin-gonic/gin"
)

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

	logger.Info(fmt.Sprintf("exec clear db, user: %s, db: %s", username, database))
	output, err := executeLocalScript(shellPath, username, password, database)
	if err != nil {
		logger.ErrorNoStack(err)
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

// parseMySQLDSN 从 db.url（user:pass@host/db?params 形式）解析出用户名、密码、库名。
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

// resolveShellPath 在可执行文件目录及其祖先目录中查找脚本，返回首个命中路径。
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

// executeLocalScript 用本地 shell 解释器执行脚本并返回合并输出。
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

// resolveShellExecutor 查找可用的 shell 解释器，优先 bash（兼容 [[ ]] 语法）。
func resolveShellExecutor() (string, error) {
	for _, candidate := range []string{"bash", "sh"} {
		if _, err := exec.LookPath(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("未找到可用的 shell 解释器")
}
