package main

import (
	"os"
	"path/filepath"

	_ "io/github/gforgame/examples/protos"
	"io/github/gforgame/network"
	protocolexporter "io/github/gforgame/tools/protocol"
)

func main() {
	root, err := findProjectRoot()
	if err != nil {
		panic(err)
	}
	if err := os.Chdir(root); err != nil {
		panic(err)
	}

	protosDir := "examples\\protos"
	csharpOutDir := "tools\\protocol\\output\\csharp\\"
	templatePath := "tools\\protocol\\templates\\csharptemplate.tpl"
	registerFile := "examples\\protos\\register_gen.go"

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
