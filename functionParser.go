package main

import "go/ast"

// Parses a function declaration
func parseFunction(function *ast.FuncDecl, fileContext *fileContext) {
	functionName := addFunctionNode(function, fileContext)

	parseParameterTypes(function, functionName, fileContext)

	parseResultTypes(function, functionName, fileContext)
}

func addFunctionNode(function *ast.FuncDecl, fileContext *fileContext) string {
	functionName := fileContext.packageName + "." + function.Name.Name
	fileContext.addNode(functionName, "Member", fileContext.fileName, function.Doc.Text())

	return functionName
}

func parseParameterTypes(function *ast.FuncDecl, functionName string, fileContext *fileContext) {
	if function.Type.Params != nil {
		parameters := function.Type.Params.List
		addFieldsLinks(parameters, functionName, fileContext)
	}
}

func parseResultTypes(function *ast.FuncDecl, functionName string, fileContext *fileContext) {
	if function.Type.Results != nil {
		results := function.Type.Results.List
		addFieldsLinks(results, functionName, fileContext)
	}
}
