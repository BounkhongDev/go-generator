package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/BounkhongDev/go-generator/pkg/generator"
)

func usage() {
	fmt.Fprintf(os.Stderr, "go-gen-r - simple Go project generator\n\n")
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r init\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r <module_name>\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r test <module_name>\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r test <module_name> --force\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r auto-test <module_name>\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r auto-test <module_name> --force\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r run-test\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r init\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r user\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r test users\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r test users --force\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r auto-test users\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r auto-test users --force\n")
	fmt.Fprintf(os.Stderr, "  go-gen-r run-test\n\n")
	fmt.Fprintf(os.Stderr, "Notes:\n")
	fmt.Fprintf(os.Stderr, "  - module_name must not contain '-'\n")
	fmt.Fprintf(os.Stderr, "  - use underscores for multi-word names (e.g. user_account)\n")
	fmt.Fprintf(os.Stderr, "  - 'test <module>' generates test files\n")
	fmt.Fprintf(os.Stderr, "  - add '--force' to regenerate existing test files\n")
	fmt.Fprintf(os.Stderr, "  - 'auto-test <module>' regenerates service tests from service methods\n")
	fmt.Fprintf(os.Stderr, "  - 'run-test' runs: go test ./...\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	if command == "init" {
		fmt.Print("Enter Project Name: ")
		var projectName string
		fmt.Scan(&projectName)

		if err := generator.Init(projectName); err != nil {
			log.Fatal(err)
		}
		return
	}

	if command == "test" {
		if flag.NArg() < 2 {
			log.Fatal("Usage: go-gen-r test <module_name>")
		}

		moduleName := strings.TrimSpace(flag.Arg(1))
		if moduleName == "" {
			log.Fatal("Module name must not be empty")
		}
		if strings.Contains(moduleName, "-") {
			log.Fatal("Module name must not contain -")
		}

		projectName, err := generator.ResolveProjectName()
		if err != nil {
			log.Fatal(err)
		}

		force := false
		for i := 2; i < flag.NArg(); i++ {
			if flag.Arg(i) == "--force" {
				force = true
			}
		}

		if force {
			if err := generator.GenerateTestFilesForce(moduleName, projectName); err != nil {
				log.Fatal(err)
			}
			return
		}

		if err := generator.GenerateTestFiles(moduleName, projectName); err != nil {
			log.Fatal(err)
		}
		return
	}

	if command == "run-test" {
		cmd := exec.Command("go", "test", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if command == "auto-test" {
		if flag.NArg() < 2 {
			log.Fatal("Usage: go-gen-r auto-test <module_name>")
		}

		moduleName := strings.TrimSpace(flag.Arg(1))
		if moduleName == "" {
			log.Fatal("Module name must not be empty")
		}
		if strings.Contains(moduleName, "-") {
			log.Fatal("Module name must not contain -")
		}

		projectName, err := generator.ResolveProjectName()
		if err != nil {
			log.Fatal(err)
		}

		force := false
		for i := 2; i < flag.NArg(); i++ {
			if flag.Arg(i) == "--force" {
				force = true
			}
		}

		if err := generator.GenerateAutoServiceTests(moduleName, projectName, force); err != nil {
			log.Fatal(err)
		}
		return
	}

	{
		moduleName := command
		if strings.Contains(moduleName, "-") {
			log.Fatal("Module name must not contain -")
		}
		if err := generator.GenerateModule(moduleName); err != nil {
			log.Fatal(err)
		}
	}
}
