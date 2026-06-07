package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	action := "all"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	root, err := findProjectRoot()
	if err != nil {
		panic(err)
	}

	var generators []string
	switch action {
	case "all":
		generators = []string{"./internal/tools/protocolgen", "./internal/tools/routedispatch"}
	case "proto":
		generators = []string{"./internal/tools/protocolgen"}
	case "route":
		generators = []string{"./internal/tools/routedispatch"}
	default:
		fmt.Printf("未知命令: %s\n", action)
		fmt.Println("用法: go run ./internal/tools/gen [all|proto|route]")
		os.Exit(2)
	}

	for _, gen := range generators {
		if err := runGenerator(root, gen); err != nil {
			panic(err)
		}
	}
}

func runGenerator(root, gen string) error {
	cmd := exec.Command("go", "run", gen)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
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
