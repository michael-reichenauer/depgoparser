// The Go source file parser to produce a Dependinator model in json format
package main

import (
	"flag"
	"log"
	"path/filepath"
)

var (
	output = flag.String("o", "", "output model path if other than default '<name>.model.json' within the project root folder")
)

// The main entry point
func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		flag.Usage()
		log.Fatal("Error: No go project folder path was specified")
		return
	}

	projectPath := args[0]
	modelFilePath := *output

	if *output == "" {
		modelName := filepath.Base(projectPath)
		modelFilePath = projectPath + "\\" + modelName + ".model.json"
	}

	err := parseProject(projectPath, modelFilePath)
	if err != nil {
		reportError(err)
	}
}

// projectPath := "C:\\Users\\micha\\Go\\src\\github.com\\michael-reichenauer\\depGoParser"
// projectPath := "C:\\Users\\micha\\Go\\src\\golang.org\\x\\tools\\cmd\\guru"
