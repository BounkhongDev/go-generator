package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/BounkhongDev/go-generator/generators"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: go-gen-r <module_name>|init")
		os.Exit(1)
	}

	moduleName := flag.Arg(0)

	if moduleName == "init" {
		fmt.Print("Enter Project Name: ")
		var projectName string
		fmt.Scan(&projectName)
		projectName = strings.TrimSpace(projectName)

		if projectName == "" {
			log.Fatal("Project name must not empty")
		}
		if strings.Contains(projectName, " ") {
			log.Fatal("Project name must not contain space")
		}

		runCmd("go", "mod", "init", projectName)
		runCmd("go", "mod", "tidy")
		initPackages()
		generators.GenerateInitialStructure()
	} else {
		if strings.Contains(moduleName, "-") {
			log.Fatal("Module name must not contain -")
		}
		generators.GenerateModules(moduleName)
	}
}

func initPackages() {
	packages := []string{
		"github.com/gofiber/fiber/v2",
		"gorm.io/gorm",
		"github.com/go-playground/validator/v10",
		"go.uber.org/zap",
		"github.com/spf13/viper",
		"gorm.io/driver/postgres",
	}
	for _, pkg := range packages {
		runCmd("go", "get", pkg)
	}
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatalf("Command %s %s failed: %v", name, strings.Join(args, " "), err)
	}
}
