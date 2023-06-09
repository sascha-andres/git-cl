package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sascha-andres/git-cl/internal"
	"github.com/sascha-andres/reuse/flag"
)

var (
	version, configFile string
	help, printConfig   bool
)

// main you know
func main() {
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	options := make([]internal.ChangeLogGeneratorOption, 0)
	if printConfig {
		options = append(options, internal.WithPrintConfiguration())
	}
	if !printConfig && version != "" {
		options = append(options, internal.WithVersion(version))
	}
	if !printConfig && configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, fmt.Sprintf("error reading configuration file: %s", err))
			os.Exit(1)
		}
		var c internal.ChangeLogGenerator
		err = json.Unmarshal(data, &c)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, fmt.Sprintf("error parsing configuration: %s", err))
			os.Exit(1)
		}
		options = append(options, internal.WithConfiguration(&c))
	}

	clg, err := internal.NewChangeLogGenerator(os.Stdin, options...)
	if err != nil {
		os.Exit(1)
	}

	var result string

	result, err = clg.Build()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	if result != "" {
		fmt.Println(result)
	}
}

func init() {
	flag.BoolVar(&help, "help", false, "show help")
	flag.StringVar(&version, "version", "", "provide version")
	flag.StringVar(&configFile, "config-file", "", "provide path to config file")
	flag.BoolVar(&printConfig, "print-config", false, "pass to print used config")
}
