package protocol

import (
	"fmt"
	"strings"
)

// CSharpGenerator C#协议生成器
type CSharpGenerator struct {
	BaseGenerator // 嵌入基类复用通用逻辑
	typeMap       map[string]string
}

// NewCSharpGenerator 创建C#生成器实例
func NewCSharpGenerator(goDir, outputDir,tplPath string) *CSharpGenerator {
	return &CSharpGenerator{
		BaseGenerator: BaseGenerator{
			GoDir:     goDir,
			TemplatePath: tplPath,
			OutputDir: outputDir,
		},
		typeMap: map[string]string{
			"int":     "int",
			"string":  "string",
			"bool":    "bool",
			"int8":    "sbyte",
			"int16":   "short",
			"int32":   "int",
			"int64":   "long",
			"uint":    "uint",
			"uint8":   "byte",
			"uint16":  "ushort",
			"uint32":  "uint",
			"uint64":  "ulong",
			"float32": "float",
			"float64": "double",
		},
	}
}

// GetFileSuffix C#文件后缀
func (c *CSharpGenerator) GetFileSuffix() string {
	return ".cs"
}

// MapType Go类型 → C#类型映射
func (c *CSharpGenerator) MapType(goType string) string {
	// 处理切片
	if strings.HasPrefix(goType, "slice<") {
		elemType := strings.TrimSuffix(strings.TrimPrefix(goType, "slice<"), ">")
		mappedElem := c.typeMap[elemType]
		if mappedElem == "" {
			mappedElem = elemType
		}
		return fmt.Sprintf("List<%s>", mappedElem)
	}
	// 处理数组
	if strings.HasPrefix(goType, "array<") {
		elemType := strings.TrimSuffix(strings.TrimPrefix(goType, "array<"), ">")
		mappedElem := c.typeMap[elemType]
		if mappedElem == "" {
			mappedElem = elemType
		}
		return fmt.Sprintf("%s[]", mappedElem)
	}
	// 普通类型
	mappedType := c.typeMap[goType]
	if mappedType == "" {
		mappedType = goType
	}
	return mappedType
}

// cSharpTemplateData C#模板数
type cSharpTemplateData struct {
	StructName    string
	StructComment string
	Cmd           interface{}
	Fields        []cSharpField
}

type cSharpField struct {
	Name       string
	FieldType string
	Comment    string
}

// BuildTemplateData 构建C#模板数据
func (c *CSharpGenerator) BuildTemplateData(si structInfo, msgIds map[string]int) interface{} {
	var fields []cSharpField
	for _, f := range si.Fields {
		fields = append(fields, cSharpField{
			Name:       f.Name,
			FieldType: c.MapType(f.Type),
			Comment:    f.Comment,
		})
	}

	data := cSharpTemplateData{
		StructName:    si.Name,
		StructComment: si.Comment,
		Fields:        fields,
	}
	if val, ok := msgIds[si.Name]; ok {
		data.Cmd = val
	}
	return data
}

// Generate 暴露给外部的生成入口
func (c *CSharpGenerator) Generate(msgIds map[string]int) error {
	return c.BaseGenerator.Generate(c, msgIds)
}