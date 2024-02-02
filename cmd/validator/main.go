package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/g4s8/openbots/pkg/spec"
)

func main() {
	var configPath string
	fset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fset.StringVar(&configPath, "config", "", "config file path")
	if err := fset.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatalf("Failed to parse config: %v", err)
	}
	if configPath == "" {
		fset.Usage()
		log.Fatal("config file path is empty")
	}

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config at: %v", err)
	}

	spec, err := spec.ParseYaml(f)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	if err := spec.Validate(); err != nil {
		fmt.Printf("There are validation errors: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Config is valid")
}
