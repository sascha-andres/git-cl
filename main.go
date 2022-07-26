package main

import (
	"encoding/json"
	"fmt"
	"github.com/sascha-andres/flag"
	"github.com/sascha-andres/git-cl/internal"
	"io/ioutil"
	"os"
)

var (
	version, configFile string
	printConfig         bool
)

// main you know
func main() {
	flag.Parse()

	options := make([]internal.ChangeLogGeneratorOption, 0)
	if printConfig {
		options = append(options, internal.WithPrintConfiguration())
	}
	if !printConfig && version != "" {
		options = append(options, internal.WithVersion(version))
	}
	if !printConfig && configFile != "" {
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("error reading configuration file: %s", err))
			os.Exit(1)
		}
		var c internal.ChangeLogGenerator
		err = json.Unmarshal(data, &c)
		if err != nil {
			fmt.Fprint(os.Stderr, fmt.Sprintf("error parsing configuration: %s", err))
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
	flag.StringVar(&version, "version", lookupEnvOrString("GIT_CL_VERSION", ""), "provide version")
	flag.StringVar(&configFile, "config-file", lookupEnvOrString("GIT_CL_CONFIG_FILE", ""), "provide path to config file")
	flag.BoolVar(&printConfig, "print-config", false, "pass to print used config")
}

// lookupEnvOrString returns a default value or a value based on an env variable
func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
