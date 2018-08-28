package main

import (
	"fmt"
	"go/ast"
	"strings"
)

// Parses type declarations
func parseType(typeSpec *ast.TypeSpec, typeDeclaration *ast.GenDecl, fileContext *fileContext) {
	typeName := addTypeNode(typeSpec, typeDeclaration, fileContext)

	switch typeType := typeSpec.Type.(type) {
	case *ast.StructType:
		parseStructFields(typeType, typeName, fileContext)

	case *ast.ArrayType, ast.Expr:
		addTypeLinks(typeType, typeName, fileContext)

	default:
		reportError(fmt.Errorf("unexpected type %d", typeDeclaration.Pos()))
	}
}

func parseStructFields(structType *ast.StructType, structName string, fileContext *fileContext) {
	for _, field := range structType.Fields.List {
		fieldName := structName + "." + field.Names[0].Name
		comment := field.Doc.Text()
		fileContext.addNode(fieldName, "Member", "", comment)

		addTypeLinks(field.Type, fieldName, fileContext)
	}
}

func addTypeLinks(typeExpression ast.Expr, sourceName string, fileContext *fileContext) {
	switch typeType := typeExpression.(type) {
	case *ast.MapType:
		addMapTypeLinks(typeType, sourceName, fileContext)

	case *ast.FuncType:
		addFieldsLinks(typeType.Params.List, sourceName, fileContext)

	default:
		filedTypeName := getTypeFullName(typeExpression, fileContext)
		if filedTypeName != "" {
			fileContext.addLink(sourceName, filedTypeName, "Type")
		}
	}
}

func addMapTypeLinks(mapType *ast.MapType, sourceName string, fileContext *fileContext) {
	keyTypeName := getTypeFullName(mapType.Key, fileContext)
	if keyTypeName != "" {
		fileContext.addLink(sourceName, keyTypeName, "Type")
	}

	valueTypeName := getTypeFullName(mapType.Value, fileContext)
	if valueTypeName != "" {
		fileContext.addLink(sourceName, valueTypeName, "Type")
	}
}

func addFieldsLinks(fields []*ast.Field, sourceName string, fileContext *fileContext) {
	for _, field := range fields {
		addTypeLinks(field.Type, sourceName, fileContext)
	}
}

func getTypeFullName(fieldExpression ast.Expr, fileContext *fileContext) string {
	typeName := getTypeName(fieldExpression)

	switch typeName {
	case "", "string", "int", "error", "bool", "byte", "rune",
		"int8", "int16", "int32", "int64", "uin8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex64", "complex128":
		return ""
	default:
		parts := strings.Split(typeName, ".")
		partsCount := len(parts)
		if partsCount > 1 {
			packageName := parts[partsCount-2]
			packageFullName, ok := fileContext.packages[packageName]
			if !ok {
				reportError(fmt.Errorf("unknown field type type %d", fieldExpression.Pos()))
				return typeName
			}

			if isStandardImport(packageFullName) {
				// Ignoring common standard types
				return ""
			}

			return packageFullName + "." + parts[partsCount-1]
		}

		return fileContext.packageName + "." + typeName
	}
}

func getTypeName(typeExpression ast.Expr) string {
	switch fieldType := typeExpression.(type) {
	case *ast.Ident:
		return fieldType.Name
	case *ast.ArrayType:
		return getTypeName(fieldType.Elt)
	case *ast.StarExpr:
		return getTypeName(fieldType.X)
	case *ast.SelectorExpr:
		return getTypeName(fieldType.X) + "." + getTypeName(fieldType.Sel)
	case *ast.InterfaceType:
		return ""
	case *ast.Ellipsis:
		return getTypeName(fieldType.Elt)
	}

	reportError(fmt.Errorf("unexpected type %d", typeExpression.Pos()))
	return "__unknown__"
}

func addTypeNode(typeSpec *ast.TypeSpec, declaration *ast.GenDecl, fileContext *fileContext) string {
	typeName := fileContext.packageName + "." + typeSpec.Name.String()
	comment := declaration.Doc.Text()

	fileContext.addNode(typeName, "Type", fileContext.fileName, comment)

	return typeName
}
