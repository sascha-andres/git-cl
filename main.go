package main

import (
	"fmt"
	"github.com/sascha-andres/git-cl/internal"
	"os"
	"strings"
)

func main() {

	someString := `feat: abcde
fix(#1): qwerty
feat: poiuy
`
	myReader := strings.NewReader(someString)

	clg, err := internal.NewChangeLogGenerator(myReader)
	if err != nil {
		os.Exit(1)
	}
	_, err = clg.Build()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
}
