package logger

import (
	"io/github/gforgame/domain"
)

func LogPlayer(player domain.Player, logType Type, args ...interface{}) {
	// 构建基础参数
	baseArgs := []interface{}{
		"playerId", player.GetId(),
		"name", player.GetName(),
	}
	// 合并所有参数
	allArgs := append(baseArgs, args...)
	Log(logType, allArgs...)
}
