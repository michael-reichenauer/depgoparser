// The Go source file parser to produce a Dependinator model in json format
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	output = flag.String("o", "", "output model path if other than default '<name>.model.json' within the project root folder")
)

// The main entry point
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "> %s [-o output] <go-folder>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		log.Fatal("Error: Need to specify one project folder path")
	}

	projectPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	modelFilePath := *output

	if *output == "" {
		modelName := filepath.Base(projectPath)
		modelFilePath = projectPath + "\\" + modelName + ".model.json"
	}

	err = parseProject(projectPath, modelFilePath)
	if err != nil {
		reportError(err)
	}
}

// projectPath := "C:\\Users\\micha\\Go\\src\\github.com\\michael-reichenauer\\depgoparser"
// projectPath := "C:\\Users\\micha\\Go\\src\\golang.org\\x\\tools\\cmd\\guru"
