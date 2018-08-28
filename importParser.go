package main

import (
	"go/ast"
	"path/filepath"
	"strings"
)

// Parse import declarations
func parseImport(importSpec *ast.ImportSpec, fileContext *fileContext) {
	importPath := strings.Trim(importSpec.Path.Value, "\"")

	importFullName := strings.Replace(importPath, "/", ".", -1)
	importShortName := filepath.Base(importPath)

	// Store a map from short to full name so other parsaer faunctions can get full name from short name
	fileContext.packages[importShortName] = importFullName

	if isStandardImport(importPath) {
		// Ignoring links to common standard imports
		return
	}

	fileContext.addLink(fileContext.packageName, importFullName, "NameSpace")
}

func isStandardImport(path string) bool {
	switch path {
	case "fmt", "io/ioutil", "os", "path/filepath", "strings", "log", "io", "io.ioutil", "bufio", "bytes":
		return true
	}

	return false
}
