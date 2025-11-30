package protocol

// structField 存储单个结构体字段信息
type structField struct {
	Name    string // 字段名（如 Channel）
	Type    string // Go 类型（如 int）
	Comment string // 注释（如 发送频道：1个人 2世界）
	JsonTag string // JSON Tag（如 channel）
}

// structInfo 存储单个结构体信息
type structInfo struct {
	Name    string        // 结构体名（如 ReqChat）
	Comment string        // 结构体注释（如 聊天请求）
	Fields  []structField // 字段列表
}

// 模板所需的数据结构（与模板变量对应）
type templateData struct {
	Namespace  string          // C# 命名空间
	StructName string          // 结构体名称
	StructComment    string          // 结构体注释
	Cmd        interface{}     // cmd 数字
	Fields     []templateField // 字段列表（适配模板）
}

// 单个字段的模板数据（包含C#类型转换后的值）
type templateField struct {
	Name      string // 字段名
	Comment   string // 注释
}