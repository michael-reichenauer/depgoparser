package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Parses go files within the project folder and writes the corresponding json Dependinator model file
func parseProject(projectFolderPath, modelFilePath string) error {
	modelItems := newModelItems()

	modelItems, err := parseProjectGoFiles(projectFolderPath, modelItems)
	if err != nil {
		reportError(err)
		return err
	}

	return writeModelFile(modelFilePath, modelItems)
}

// Parses all go files with the project folder, ignoring test and test data files
func parseProjectGoFiles(projectFolderPath string, modelItems modelItems) ([]item, error) {
	modelItems, projectName := addProjectNode(modelItems, projectFolderPath)

	err := filepath.Walk(projectFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if isIgnoreFile(path) {
			return nil
		}

		modelItems, err = parseGoFile(path, projectName, modelItems)
		return err
	})

	return modelItems, err
}

func parseGoFile(path, workName string, modelItems modelItems) (modelItems, error) {
	goSyntaxTree, err := parseGoSyntax(path)
	if err != nil {
		reportError(err)
		return modelItems, err
	}

	fileInfo := newFileContext(path, goSyntaxTree, modelItems)

	addPackageNode(goSyntaxTree, workName, &fileInfo)
	addFileNode(path, &fileInfo)

	parseFileDeclarations(goSyntaxTree, &fileInfo)

	return fileInfo.modelItems, nil
}

func parseGoSyntax(filename string) (*ast.File, error) {
	parserMode := parser.ParseComments
	fset := token.NewFileSet()
	goSyntaxTree, err := parser.ParseFile(fset, filename, nil, parserMode)
	return goSyntaxTree, err
}

func parseFileDeclarations(goSyntaxTree *ast.File, fileContext *fileContext) {
	for _, declaration := range goSyntaxTree.Decls {
		parseDeclaration(declaration, fileContext)
	}
}

func parseDeclaration(declaration ast.Decl, fileContext *fileContext) {
	switch declarationType := declaration.(type) {
	case *ast.FuncDecl:
		parseFunction(declarationType, fileContext)

	case *ast.GenDecl:
		for _, specification := range declarationType.Specs {
			switch specificationType := specification.(type) {
			case *ast.ImportSpec:
				parseImport(specificationType, fileContext)

			case *ast.TypeSpec:
				parseType(specificationType, declarationType, fileContext)

			case *ast.ValueSpec:
				parseValues(specificationType, fileContext)

			default:
				reportError(fmt.Errorf("unknown declaration type: %s", declarationType.Tok))
			}
		}
	default:
		reportError(fmt.Errorf("unknown declaration %d", declaration.Pos()))
	}
}

// Add the root project node, which will contain all parsed nodes
func addProjectNode(modelItems modelItems, projectFolderPath string) (modelItems, string) {
	projectName := "$" + filepath.Base(projectFolderPath)
	modelItems = modelItems.addNode(projectName, "NameSpace", "", "Go project")
	return modelItems, projectName
}

// Add the package node, which will contain file nodes within a package
func addPackageNode(goSyntaxTree *ast.File, workName string, fileContext *fileContext) {
	description := goSyntaxTree.Doc.Text()
	fileContext.addNode(fileContext.packageName, "NameSpace", workName, description)
}

// Add a file node, which will contain all items within a file
func addFileNode(path string, fileContext *fileContext) {
	fileContext.addNode(fileContext.fileName, "NameSpace", fileContext.packageName, "go file")
}

// Ignoring non go files as well as testdata and test files
func isIgnoreFile(path string) bool {
	return !strings.HasSuffix(path, ".go") ||
		strings.Contains(path, "\\testdata\\") ||
		strings.Contains(path, "_test.go")
}
