package main

import (
	"path/filepath"
	"strings"
)

// Contains data related to a file, while parsing the file
type fileContext struct {
	rootPath    string
	fileName    string
	packageName string
	packages    map[string]string
	modelItems  modelItems
}

func newFileContext(rootPath, path, packageName string, modelItems []item) fileContext {
	info := fileContext{}
	info.packageName = packageName
	info.packages = make(map[string]string)

	fileName := filepath.Base(path)
	fileName = strings.Replace(fileName, ".", "*", -1)
	info.fileName = info.packageName + "." + fileName

	info.modelItems = modelItems

	return info
}

func (f *fileContext) addNode(name, nodeType, parent, description string) {
	f.modelItems = f.modelItems.addNode(name, nodeType, parent, description)
}

func (f *fileContext) addLink(source, target, targetType string) {
	f.modelItems = f.modelItems.addLink(source, target, targetType)
}
