package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type routeMethod struct {
	Cmd          int32
	ReceiverType string
	MethodName   string
	ReqType      string
	HasIndex     bool
	HasReturn    bool
}

func main() {
	root, err := findProjectRoot()
	if err != nil {
		panic(err)
	}

	reqCmdMap, err := parseReqCmdMap(filepath.Join(root, "internal", "protos", "register_gen.go"))
	if err != nil {
		panic(err)
	}

	methods, err := parseRouteMethods(filepath.Join(root, "internal", "route"), reqCmdMap)
	if err != nil {
		panic(err)
	}

	outFile := filepath.Join(root, "cmd", "game", "route_dispatch_gen.go")
	if err := writeGeneratedFile(outFile, methods); err != nil {
		panic(err)
	}
	
	fmt.Printf("静态路由表生成完成，输出文件：route_dispatch_gen.go\n")
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", os.ErrNotExist
		}
		wd = parent
	}
}

func parseReqCmdMap(registerFile string) (map[string]int32, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, registerFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int32)
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) != 2 {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel == nil || sel.Sel.Name != "RegisterMessage" {
			return true
		}

		cmd, ok := parseInt32Expr(call.Args[0])
		if !ok {
			return true
		}
		reqType, ok := parseReqTypeExpr(call.Args[1])
		if !ok || !strings.HasPrefix(reqType, "Req") {
			return true
		}
		result[reqType] = cmd
		return true
	})
	return result, nil
}

func parseRouteMethods(routeDir string, reqCmdMap map[string]int32) ([]routeMethod, error) {
	files, err := filepath.Glob(filepath.Join(routeDir, "*.go"))
	if err != nil {
		return nil, err
	}

	result := make([]routeMethod, 0)
	fset := token.NewFileSet()
	for _, filePath := range files {
		fileAst, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, decl := range fileAst.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv == nil || fd.Name == nil {
				continue
			}
			if !strings.HasPrefix(fd.Name.Name, "Req") {
				continue
			}
			receiverType := parseReceiverType(fd.Recv.List[0].Type)
			if receiverType == "" {
				continue
			}
			reqType, hasIndex := parseReqParam(fd.Type.Params)
			if reqType == "" {
				continue
			}
			cmd, ok := reqCmdMap[reqType]
			if !ok {
				continue
			}
			hasReturn := fd.Type.Results != nil && len(fd.Type.Results.List) > 0
			result = append(result, routeMethod{
				Cmd:          cmd,
				ReceiverType: receiverType,
				MethodName:   fd.Name.Name,
				ReqType:      reqType,
				HasIndex:     hasIndex,
				HasReturn:    hasReturn,
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Cmd != result[j].Cmd {
			return result[i].Cmd < result[j].Cmd
		}
		if result[i].ReceiverType != result[j].ReceiverType {
			return result[i].ReceiverType < result[j].ReceiverType
		}
		return result[i].MethodName < result[j].MethodName
	})
	return deduplicateByCmd(result), nil
}

func deduplicateByCmd(methods []routeMethod) []routeMethod {
	seen := make(map[int32]struct{}, len(methods))
	result := make([]routeMethod, 0, len(methods))
	for _, m := range methods {
		if _, ok := seen[m.Cmd]; ok {
			continue
		}
		seen[m.Cmd] = struct{}{}
		result = append(result, m)
	}
	return result
}

func parseReceiverType(expr ast.Expr) string {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return ""
	}
	ident, ok := star.X.(*ast.Ident)
	if !ok {
		return ""
	}
	return ident.Name
}

func parseReqParam(params *ast.FieldList) (reqType string, hasIndex bool) {
	if params == nil {
		return "", false
	}
	reqType = ""
	hasIndex = false
	for _, field := range params.List {
		// 检测 index int32 参数
		if ident, ok := field.Type.(*ast.Ident); ok && ident.Name == "int32" {
			hasIndex = true
		}
		star, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		sel, ok := star.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pkgIdent, ok := sel.X.(*ast.Ident)
		if !ok || pkgIdent.Name != "protos" {
			continue
		}
		if strings.HasPrefix(sel.Sel.Name, "Req") {
			reqType = sel.Sel.Name
		}
	}
	return reqType, hasIndex
}

func parseInt32Expr(expr ast.Expr) (int32, bool) {
	switch v := expr.(type) {
	case *ast.BasicLit:
		n, err := strconv.ParseInt(v.Value, 10, 32)
		if err != nil {
			return 0, false
		}
		return int32(n), true
	case *ast.UnaryExpr:
		if v.Op != token.SUB {
			return 0, false
		}
		n, ok := parseInt32Expr(v.X)
		if !ok {
			return 0, false
		}
		return -n, true
	default:
		return 0, false
	}
}

func parseReqTypeExpr(expr ast.Expr) (string, bool) {
	unary, ok := expr.(*ast.UnaryExpr)
	if !ok || unary.Op != token.AND {
		return "", false
	}
	comp, ok := unary.X.(*ast.CompositeLit)
	if !ok {
		return "", false
	}
	ident, ok := comp.Type.(*ast.Ident)
	if !ok {
		return "", false
	}
	return ident.Name, true
}

func writeGeneratedFile(filePath string, methods []routeMethod) error {
	var b bytes.Buffer
	b.WriteString("// Code generated by tools/routedispatch. DO NOT EDIT.\n\n")
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n\n")
	b.WriteString("\t\"github.com/forfun/gforgame/internal/protos\"\n")
	b.WriteString("\t\"github.com/forfun/gforgame/internal/route\"\n")
	b.WriteString("\t\"github.com/forfun/gforgame/network\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func init() {\n")
	b.WriteString("\tgeneratedRouteDispatchers = map[int32]generatedRouteInvoker{\n")

	for _, m := range methods {
		b.WriteString(fmt.Sprintf("\t\t%d: func(msgHandler *network.Handler, session *network.Session, index int32, msg any) (any, error) {\n", m.Cmd))
		b.WriteString(fmt.Sprintf("\t\tr, ok := msgHandler.Receiver.Interface().(*route.%s)\n", m.ReceiverType))
		b.WriteString("\t\tif !ok {\n")
		b.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"generated dispatch receiver type mismatch: cmd=%d expect=*route.%s\")\n", m.Cmd, m.ReceiverType))
		b.WriteString("\t\t}\n")
		b.WriteString(fmt.Sprintf("\t\treq, ok := msg.(*protos.%s)\n", m.ReqType))
		b.WriteString("\t\tif !ok {\n")
		b.WriteString(fmt.Sprintf("\t\t\treturn nil, fmt.Errorf(\"generated dispatch msg type mismatch: cmd=%d expect=*protos.%s\")\n", m.Cmd, m.ReqType))
		b.WriteString("\t\t}\n")
		if m.HasReturn {
			if m.HasIndex {
				b.WriteString(fmt.Sprintf("\t\treturn r.%s(session, index, req), nil\n", m.MethodName))
			} else {
				b.WriteString(fmt.Sprintf("\t\treturn r.%s(session, req), nil\n", m.MethodName))
			}
		} else {
			if m.HasIndex {
				b.WriteString(fmt.Sprintf("\t\tr.%s(session, index, req)\n", m.MethodName))
			} else {
				b.WriteString(fmt.Sprintf("\t\tr.%s(session, req)\n", m.MethodName))
			}
			b.WriteString("\t\treturn nil, nil\n")
		}
		b.WriteString("\t\t},\n")
	}

	b.WriteString("\t}\n")
	b.WriteString("}\n")
	return os.WriteFile(filePath, b.Bytes(), 0644)
}
