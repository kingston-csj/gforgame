package tools

import (
	"fmt"
	"strings"
)

// -------------------------- TypeScript 子类（仅实现差异化逻辑） --------------------------
// TypeScriptGenerator TS协议生成器
type TypeScriptGenerator struct {
	BaseGenerator // 嵌入基类复用通用逻辑
	typeMap       map[string]string
}

// NewTypeScriptGenerator 创建TS生成器实例
func NewTypeScriptGenerator(goDir, outputDir, tplPath string) *TypeScriptGenerator {
	return &TypeScriptGenerator{
		BaseGenerator: BaseGenerator{
			GoDir:     goDir,
			TemplatePath:  tplPath,
			OutputDir: outputDir,
		},
		typeMap: map[string]string{
			"int":     "number",
			"string":  "string",
			"bool":    "boolean",
			"int8":    "number",
			"int16":   "number",
			"int32":   "number",
			"int64":   "number",
			"uint":    "number",
			"uint8":   "number",
			"uint16":  "number",
			"uint32":  "number",
			"uint64":  "number",
			"float32": "number",
			"float64": "number",
		},
	}
}

// GetFileSuffix TS文件后缀
func (t *TypeScriptGenerator) GetFileSuffix() string {
	return ".ts"
}

// MapType Go类型 → TS类型映射（子类差异化实现）
func (t *TypeScriptGenerator) MapType(goType string) string {
	// 处理切片
	if strings.HasPrefix(goType, "slice<") {
		elemType := strings.TrimSuffix(strings.TrimPrefix(goType, "slice<"), ">")
		mappedElem := t.typeMap[elemType]
		if mappedElem == "" {
			mappedElem = elemType
		}
		return fmt.Sprintf("Array<%s>", mappedElem)
	}
	// 处理数组
	if strings.HasPrefix(goType, "array<") {
		elemType := strings.TrimSuffix(strings.TrimPrefix(goType, "array<"), ">")
		mappedElem := t.typeMap[elemType]
		if mappedElem == "" {
			mappedElem = elemType
		}
		return fmt.Sprintf("Array<%s>", mappedElem)
	}
	// 普通类型
	mappedType := t.typeMap[goType]
	if mappedType == "" {
		mappedType = goType
	}
	return mappedType
}

// TSTemplateData TS模板数据（子类专属）
type TSTemplateData struct {
	Cmd   int
	ClassName     string
	ClassComment       string
	Fields        []TSField
}

type TSField struct {
	Name     string
	FieldType   string
	Comment  string
	JsonTag  string
}

// BuildTemplateData 构建TS模板数据（子类差异化实现）
func (t *TypeScriptGenerator) BuildTemplateData(si StructInfo, msgIds map[string]int) interface{} {
	var fields []TSField
	for _, f := range si.Fields {
		fields = append(fields, TSField{
			Name:     f.Name,
			FieldType:   t.MapType(f.Type),
			Comment:  f.Comment,
			JsonTag:  f.JsonTag,
		})
	}

	return TSTemplateData{
		Cmd:   msgIds[si.Name],
		ClassName:     si.Name,
		ClassComment:       si.Comment,
		Fields:        fields,
	}
}

// Generate 暴露给外部的生成入口
func (t *TypeScriptGenerator) Generate(msgIds map[string]int) error {
	return t.BaseGenerator.Generate(t, msgIds)
}