package protocol

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// ProtocolGenerator 协议生成器抽象基类
type ProtocolGenerator interface {
	// 通用生成入口（基类实现）
	Generate(msgIds map[string]int) error
	// 返回文件后缀（.cs/.ts）
	GetFileSuffix() string   
	// Go类型 → 目标语言类型映射           
	MapType(goType string) string
	// 构建模板数据
	BuildTemplateData(si structInfo, msgIds map[string]int) interface{} 
	// 返回模板文件路径
	GetTemplatePath() string             
}

// BaseGenerator 基类实现通用逻辑，子类嵌入该结构体复用
type BaseGenerator struct {
	GoDir      string // Go源码目录
	OutputDir  string // 生成文件输出目录
	template   *template.Template       // 解析后的模板
	TemplatePath string // 模板文件路径
}

func (b *BaseGenerator) Init(g ProtocolGenerator) error {
	// 解析模板：通过传入的 ProtocolGenerator 子类（C#/TS）获取模板路径
	tpl, err := template.ParseFiles(g.GetTemplatePath())
	if err != nil {
		return fmt.Errorf("解析模板失败：%w", err)
	}
	b.template = tpl // 保存解析后的模板到基类属性
	// 创建输出目录：确保生成文件的目录存在
	if err := os.MkdirAll(b.OutputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败：%w", err)
	}
	return nil
}

func (b *BaseGenerator) GetTemplatePath() string {
	return b.TemplatePath
}

// Generate 通用生成逻辑（所有语言共享）
func (b *BaseGenerator) Generate(g ProtocolGenerator, msgIds map[string]int) error {
	// 初始化模板
	if err := b.Init(g); err != nil {
		return err
	}

	// 读取Go目录下所有文件
	files, err := os.ReadDir(b.GoDir)
	if err != nil {
		return fmt.Errorf("读取Go目录失败：%w", err)
	}

	// 遍历文件解析并生成
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		// 通用AST解析
		filePath := b.GoDir + "\\" + file.Name()
		structInfos, err := b.parseGoFile(filePath)
		if err != nil {
			fmt.Printf("解析文件 %s 失败：%v\n", filePath, err)
			continue
		}

		// 为每个结构体生成文件
		for _, si := range structInfos {
			if err := b.generateStructFile(g, si, msgIds); err != nil {
				fmt.Printf("生成结构体 %s 失败：%v\n", si.Name, err)
				continue
			}
			// fmt.Printf("已生成：%s%s\n", b.OutputDir, si.Name+g.GetFileSuffix())
		}
	}

	fmt.Printf("%s协议生成完成，输出目录：%s\n", g.GetFileSuffix()[1:], b.OutputDir)
	return nil
}

// generateStructFile 生成单个结构体文件（通用逻辑）
func (b *BaseGenerator) generateStructFile(g ProtocolGenerator, si structInfo, msgIds map[string]int) error {
	// 子类构建模板数据
	data := g.BuildTemplateData(si, msgIds)
	// 渲染模板
	var buf bytes.Buffer
	if err := b.template.Execute(&buf, data); err != nil {
		panic(fmt.Sprintf("渲染模板失败：%v", err))
	}
	// 写入文件
	outputPath := b.OutputDir + si.Name + g.GetFileSuffix()
	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入文件失败：%w", err)
	}
	return nil
}

// parseGoFile 通用AST解析
func (b *BaseGenerator) parseGoFile(filePath string) ([]structInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	cm := ast.NewCommentMap(fset, node, node.Comments)
	var structInfos []structInfo

	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structName := ts.Name.Name
		structLine := fset.Position(ts.Pos()).Line
		var fields []structField

		if structType.Fields != nil {
			for _, field := range structType.Fields.List {
				var fieldName string
                if len(field.Names) > 0 {
                    fieldName = field.Names[0].Name
                } else {
                    continue
                }
                if fieldName == "_" {
                    continue
                }

				fieldType := b.getFieldTypeStr(field.Type)
				fieldComment := b.getCommentText(field.Comment)
				jsonTag := ""
				if field.Tag != nil {
					tagStr := strings.Trim(field.Tag.Value, "`")
					jsonTag = b.extractJsonTag(tagStr)
				}

				fields = append(fields, structField{
					Name:    fieldName,
					Type:    fieldType,
					Comment: fieldComment,
					JsonTag: jsonTag,
				})
			}
		}

		// 提取结构体注释（按行号过滤）
		var structComment string
		commentGroups := cm.Filter(ts).Comments()
		for _, cg := range commentGroups {
			for _, c := range cg.List {
				commentLine := fset.Position(c.Pos()).Line
				if commentLine > structLine {
					continue
				}

				text := strings.TrimSpace(c.Text)
				if strings.HasPrefix(text, "//go:") {
					continue
				}
				text = strings.TrimPrefix(text, "//")
				text = strings.TrimPrefix(text, "/*")
				text = strings.TrimSuffix(text, "*/")
				text = strings.TrimSpace(text)
				if text != "" {
					if structComment != "" {
						structComment += "\n"
					}
					structComment += text
				}
			}
		}

		if structComment == "" {
			structComment = b.getCommentText(ts.Doc)
			if structComment == "" {
				structComment = b.getCommentText(ts.Comment)
			}
		}

		structInfos = append(structInfos, structInfo{
			Name:    structName,
			Comment: structComment,
			Fields:  fields,
		})

		return true
	})

	return structInfos, nil
}

// -------------------------- 通用工具方法--------------------------
func (b *BaseGenerator) getFieldTypeStr(expr ast.Expr) string {
    if starExpr, ok := expr.(*ast.StarExpr); ok {
        return b.getFieldTypeStr(starExpr.X)
    }

	if arrayExpr, ok := expr.(*ast.ArrayType); ok {
		elemType := b.getFieldTypeStr(arrayExpr.Elt)
		if arrayExpr.Len != nil {
			return fmt.Sprintf("array<%s>", elemType)
		}
		return fmt.Sprintf("slice<%s>", elemType)
	}

    if ident, ok := expr.(*ast.Ident); ok {
        return ident.Name
    }

    if _, ok := expr.(*ast.StructType); ok {
        return "struct"
    }

    if selExpr, ok := expr.(*ast.SelectorExpr); ok {
        return selExpr.Sel.Name
    }

	fmt.Printf("警告：未处理的类型节点 %T，字段类型可能解析错误", expr)
	return ""
}

func (b *BaseGenerator) getCommentText(comment *ast.CommentGroup) string {
	if comment == nil {
		return ""
	}
	var buf bytes.Buffer
	for _, c := range comment.List {
		text := strings.TrimSpace(c.Text)
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)
		if text != "" {
			buf.WriteString(text)
			buf.WriteString("\n")
		}
	}
	return strings.TrimSuffix(buf.String(), "\n")
}

func (b *BaseGenerator) extractJsonTag(tagStr string) string {
    re := regexp.MustCompile(`json:"([^"]+)"`)
    matches := re.FindStringSubmatch(tagStr)
    if len(matches) >= 2 {
        return matches[1]
    }
    return ""
}

func (b *BaseGenerator) parseAllTags(tagStr string) map[string]string {
    re := regexp.MustCompile(`(\w+):"([^"]*)"`)
    out := make(map[string]string)
    matches := re.FindAllStringSubmatch(tagStr, -1)
    for _, m := range matches {
        if len(m) >= 3 {
            out[m[1]] = m[2]
        }
    }
    return out
}

func (b *BaseGenerator) GenerateRegisterFromTags(goDir string, outputFile string, msgConsts map[string]int) error {
    if msgConsts == nil || len(msgConsts) == 0 {
        msgConsts = b.buildMsgConstMap(goDir + "\\" + "message.go")
    }
    files, err := os.ReadDir(goDir)
    if err != nil {
        return fmt.Errorf("读取Go目录失败：%w", err)
    }

    type entry struct{ Type string; Cmd int; FileName string ;}
    entries := make([]entry, 0, 64)
    var buf bytes.Buffer
    buf.WriteString("package protos\n\n")
    buf.WriteString("import (\n\t\"io/github/gforgame/network\"\n)\n\n")
    buf.WriteString("func init() {\n")

    fset := token.NewFileSet()
    for _, file := range files {
        if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
            continue
        }
        filePath := goDir + "\\" + file.Name()
        node, err := parser.ParseFile(fset, filePath, nil, 0)
        if err != nil {
            continue
        }

		
        ast.Inspect(node, func(n ast.Node) bool {
            ts, ok := n.(*ast.TypeSpec)
            if !ok {
                return true
            }
            st, ok := ts.Type.(*ast.StructType)
            if !ok {
                return true
            }
            typeName := ts.Name.Name
            if st.Fields != nil {
                for _, fld := range st.Fields.List {
                    if fld.Tag == nil {
                        continue
                    }
                    tagStr := strings.Trim(fld.Tag.Value, "`")
                    tags := b.parseAllTags(tagStr)
                    cmd := 0
                    found := false
                    if ref, ok := tags["cmd_ref"]; ok {
                        if v, ok2 := msgConsts[ref]; ok2 {
                            cmd = v
                            found = true
                        }
                    } else if s, ok := tags["cmd"]; ok {
                        if v, err := strconv.Atoi(s); err == nil {
                            cmd = v
                            found = true
                        }
                    }
                    if found {
                        entries = append(entries, entry{Type: typeName, Cmd: cmd, FileName: file.Name()})
                        break
                    }
                }
            }
            return true
        })
    }


    grouped := make(map[string][]entry)
    for _, e := range entries {
        grouped[e.FileName] = append(grouped[e.FileName], e)
    }
    for fileName, list := range grouped {
        buf.WriteString(fmt.Sprintf("\t// ----from %s----\n", fileName))
        for _, e := range list {
            buf.WriteString(fmt.Sprintf("\tnetwork.RegisterMessage(%d, &%s{})\n", e.Cmd, e.Type))
        }
        buf.WriteString("\n")
    }
    buf.WriteString("}\n")

    if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
        return fmt.Errorf("写入文件失败：%w", err)
    }
    return nil
}

func (b *BaseGenerator) buildMsgConstMap(filePath string) map[string]int {
    out := make(map[string]int)
    data, err := os.ReadFile(filePath)
    if err != nil {
        return out
    }
    re := regexp.MustCompile(`(?m)^\s*([A-Za-z_][A-Za-z0-9_]*)\s*=\s*([0-9]+)\s*$`)
    matches := re.FindAllStringSubmatch(string(data), -1)
    for _, m := range matches {
        if len(m) >= 3 {
            v, err := strconv.Atoi(m[2])
            if err == nil {
                out[m[1]] = v
            }
        }
    }
    return out
}
