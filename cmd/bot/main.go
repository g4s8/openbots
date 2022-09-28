package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/g4s8/openbots/pkg/bot"
	"github.com/g4s8/openbots/pkg/spec"
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

	doneCh := make(chan struct{})
	exitCh := make(chan os.Signal, 1)
	go func() {
		defer close(doneCh)
		<-exitCh
		log.Println("Shutting down")
		if err := bot.Stop(); err != nil {
			log.Fatal("Failed to stop bot: ", err)
		}
		log.Println("Shutdown completed")
	}()
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	<-doneCh
	log.Println("Exit")
	os.Exit(0)
}
