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
		fmt.Println("Usage: go-gen-r <module_name>||init")
		os.Exit(1)
	}

	moduleName := flag.Arg(0)
	var projectName string

	if moduleName == "init" {
		fmt.Print("Enter Project Name: ")
		fmt.Scan(&projectName)

		// validate project name must not empty and must not contain space
		if projectName == "" {
			log.Fatal("Project name must not empty")
		}

		if le := strings.Split(projectName, " "); len(le) > 1 {
			log.Fatal("Project name must not contain space")
		}

		cmd := exec.Command("go", "mod", "init", projectName)
		if errors.Is(cmd.Err, exec.ErrDot) {
			cmd.Err = nil
		}
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		cmd = exec.Command("go", "mod", "tidy")
		if errors.Is(cmd.Err, exec.ErrDot) {
			cmd.Err = nil
		}
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		initPackages()

		generators.GenerateInitialStructure()
	} else {
		generators.GenerateModules(moduleName)
	}
}

func initPackages() {
	// get fiber package
	cmd := exec.Command("go", "get", "github.com/gofiber/fiber/v2")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// get gorm package
	cmd = exec.Command("go", "get", "gorm.io/gorm")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// get validator package
	cmd = exec.Command("go", "get", "github.com/go-playground/validator/v10")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// get zaplogger package
	cmd = exec.Command("go", "get", "go.uber.org/zap")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// get viper package
	cmd = exec.Command("go", "get", "github.com/spf13/viper")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	// get postgres gorm driver
	cmd = exec.Command("go", "get", "gorm.io/driver/postgres")
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
