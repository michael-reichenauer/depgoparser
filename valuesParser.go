package main

import "go/ast"

// Parses file value declarations
func parseValues(spec *ast.ValueSpec, fileContext *fileContext) {
	for _, id := range spec.Names {
		parseValue(id, fileContext)
	}
}

func parseValue(id *ast.Ident, fileContext *fileContext) {
	valueName := fileContext.packageName + "." + id.Name
	fileContext.addNode(valueName, "Member", fileContext.fileName, "")
}
