package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BounkhongDev/go-generator/pkg/generator"
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

		if err := generator.Init(projectName); err != nil {
			log.Fatal(err)
		}
	} else {
		if strings.Contains(moduleName, "-") {
			log.Fatal("Module name must not contain -")
		}
		if err := generator.GenerateModule(moduleName); err != nil {
			log.Fatal(err)
		}
	}
}
