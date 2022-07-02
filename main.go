package main

import (
	"fmt"
	"github.com/sascha-andres/git-cl/internal"
	"os"
)

func main() {
	clg, err := internal.NewChangeLogGenerator(os.Stdin)
	if err != nil {
		os.Exit(1)
	}
	_, err = clg.Build()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
}
