package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/g4s8/openbots-go/pkg/bot"
	"github.com/g4s8/openbots-go/pkg/spec"
	"github.com/pkg/errors"
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
	if errs := spec.Validate(); len(errs) > 0 {
		fmt.Printf("There are %d validation errors:\n", len(errs))
		for i, err := range errs {
			fmt.Printf("  [%d] ERROR: %v\n", i, err)
		}
		os.Exit(1)
	}

	bot, err := bot.NewFromSpec(spec.Bot)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}
	if err := bot.Start(); err != nil {
		log.Fatal("Failed to start bot: ", err)
	}

	done := make(chan struct{})
	<-done

	if err := bot.Stop(); err != nil {
		log.Fatal("Failed to stop bot: ", err)
	}
}
