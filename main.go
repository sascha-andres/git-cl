package main

import (
	"flag"
	"fmt"
	"github.com/sascha-andres/git-cl/internal"
	"os"
)

var (
	version string
)

// main you know
func main() {
	flag.Parse()

	options := make([]internal.ChangeLogGeneratorOption, 0)
	if version != "" {
		options = append(options, internal.WithVersion(version))
	}

	clg, err := internal.NewChangeLogGenerator(os.Stdin, options...)
	if err != nil {
		os.Exit(1)
	}
	_, err = clg.Build()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
}

func init() {
	flag.StringVar(&version, "version", lookupEnvOrString("GIT_CL_VERSION", ""), "provide version")
}

// lookupEnvOrString returns a default value or a value based on an env variable
func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
