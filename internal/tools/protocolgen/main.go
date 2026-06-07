package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/forfun/gforgame/common/logger"
	_ "github.com/forfun/gforgame/internal/protos"
	protocolexporter "github.com/forfun/gforgame/internal/tools/protocol"
	"github.com/forfun/gforgame/network"
)

func main() {
	root, err := findProjectRoot()
	logger.Info(fmt.Sprintf("项目根目录: %s", root))
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(root); err != nil {
		panic(err)
	}

	protosDir := filepath.Join("internal", "protos")
	csharpOutDir := filepath.Join("internal","tools", "protocol", "output", "csharp")
	templatePath := filepath.Join("internal","tools", "protocol", "templates", "csharptemplate.tpl")
	registerFile := filepath.Join("internal", "protos", "register_gen.go")

	generator := protocolexporter.NewCSharpGenerator(
		protosDir,
		csharpOutDir,
		templatePath,
	)
	if err := generator.Generate(network.GetMsgName2IdMapper()); err != nil {
		panic(err)
	}
	if err := generator.BaseGenerator.GenerateRegisterFromTags(protosDir, registerFile, nil); err != nil {
		panic(err)
	}
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
