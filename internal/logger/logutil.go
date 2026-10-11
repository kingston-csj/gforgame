package logger

// 日志工具函数 统一对外接口
import (
	"github.com/forfun/gforgame/common/logger"
	"github.com/forfun/gforgame/config"
	"github.com/forfun/gforgame/internal/constants"
)

// topLevelLogKeys 会放在 JSON 顶层；其他字段全部合并到 ext 嵌套对象。
// {"time":...,"date":...,"sid":...,"playerId":...,"module":...,"ext":{...}}
var topLevelLogKeys = map[string]struct{}{
	"sid":        {},
	"serverId":   {},
	"playerId":   {},
	"module":     {},
	"playerName": {},
}

// topFields 日志 JSON 顶层固定字段（time/date 由 common/logger/log.go 注入）。
type topFields struct {
	Sid        uint32
	PlayerID   string
	PlayerName string
	Module     string
}

func buildTopFields(actorID, actorName string, category *constants.LoggerCategory) topFields {
	t := topFields{
		Sid:    config.ServerConfig.ServerId,
		Module: category.Code,
	}
	if actorID != "" {
		t.PlayerID = actorID
		t.PlayerName = actorName
	}
	return t
}

func (t topFields) toKVs() []any {
	kvs := make([]any, 0, 8)
	if t.Sid != 0 {
		kvs = append(kvs, "sid", t.Sid)
	}
	if t.PlayerID != "" {
		kvs = append(kvs, "playerId", t.PlayerID)
	}
	if t.PlayerName != "" {
		kvs = append(kvs, "playerName", t.PlayerName)
	}
	if t.Module != "" {
		kvs = append(kvs, "module", t.Module)
	}
	return kvs
}

// splitArgs 将扁平 KV 分成顶层追加 + ext 对象。
func splitArgs(args []any) (topExtra []any, ext map[string]any) {
	ext = make(map[string]any, len(args)/2)
	for i := 0; i+1 < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			continue
		}
		v := args[i+1]
		if _, isTop := topLevelLogKeys[k]; isTop {
			if k == "serverId" {
				topExtra = append(topExtra, "sid", v)
			} else {
				topExtra = append(topExtra, k, v)
			}
			continue
		}
		ext[k] = v
	}
	return topExtra, ext
}

// mergeTopKVs 合并固定字段 + 调用方显式传入的顶层字段；后者优先级更高（可覆盖默认值）。
func mergeTopKVs(base []any, extra []any) []any {
	if len(extra) == 0 {
		return base
	}
	uniq := make(map[string]any, len(base)/2+len(extra)/2)
	for i := 0; i+1 < len(base); i += 2 {
		k, _ := base[i].(string)
		if k == "" {
			continue
		}
		uniq[k] = base[i+1]
	}
	for i := 0; i+1 < len(extra); i += 2 {
		k, _ := extra[i].(string)
		if k == "" {
			continue
		}
		uniq[k] = extra[i+1]
	}
	out := make([]any, 0, len(uniq)*2)
	order := []string{"sid", "playerId", "playerName", "module"}
	for _, k := range order {
		if v, ok := uniq[k]; ok {
			out = append(out, k, v)
			delete(uniq, k)
		}
	}
	for k, v := range uniq {
		out = append(out, k, v)
	}
	return out
}

// packSlogAttrs 组装传给 common 层的 attrs。
// 注意：必须保持扁平 KV（string key, value 交替）
// ext 作为嵌套 map 值直接按 flat 对传："ext", map[xxx]，slog JSON handler 会输出嵌套对象。
func packSlogAttrs(category constants.LoggerCategory, actorID, actorName string, args ...any) []any {
	t := buildTopFields(actorID, actorName, &category)
	baseTop := t.toKVs()
	extraTop, ext := splitArgs(args)
	topKVs := mergeTopKVs(baseTop, extraTop)

	result := make([]any, 0, len(topKVs)+2)
	result = append(result, topKVs...)
	if len(ext) > 0 {
		result = append(result, "ext", ext)
	}
	return result
}

// Info 记录一条普通日志（仅控制台 + app 文件）
func Info(v string) {
	logger.Info(v)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

// Error 记录错误日志（带自定义描述 + 调用栈），同时写入 err 文件和控制台
func Error(customMsg string, err error) {
	logger.Error(customMsg, err)
}

// ErrorNoStack 记录单行重要错误（不带调用栈）
func ErrorNoStack(customMsg any) {
	logger.ErrorNoStack(customMsg)
}

// PanicToError 将 recover 到的 panic 值包装为带原始堆栈的 error
func PanicToError(p any) error {
	return logger.PanicToError(p)
}

// ============================================================
// 系统日志（按天滚动）
// ============================================================

func LogSystem(category constants.LoggerCategory, args ...any) {
	attrs := packSlogAttrs(category, "", "", args...)
	logger.LogSystem(category.Code, attrs...)
}

// ============================================================
// 运营/业务日志（按小时滚动，统一文件 business.%Y%m%d%H）
// ============================================================

func LogBusinessWithActor(actorID, actorName string, category constants.LoggerCategory, args ...any) {
	attrs := packSlogAttrs(category, actorID, actorName, args...)
	logger.LogBusinessTo(category.Code, attrs...)
}

func LogBusinessNoActor(category constants.LoggerCategory, args ...any) {
	attrs := packSlogAttrs(category, "", "", args...)
	logger.LogBusinessTo(category.Code, attrs...)
}

func LogBusinessDefault(args ...any) {
	extraTop, ext := splitArgs(args)
	attrs := make([]any, 0, len(extraTop)+2)
	sid := config.ServerConfig.ServerId
	if sid != 0 {
		attrs = append(attrs, "sid", sid)
	}
	attrs = append(attrs, extraTop...)
	if len(ext) > 0 {
		attrs = append(attrs, "ext", ext)
	}
	logger.LogBusiness(attrs...)
}
