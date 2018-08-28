package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

// The model, which has format version and an array of nodes and links
type model struct {
	FormatVersion string
	Items         []item
}

// The item, which can be either a node or a link
type item struct {
	Node *node `json:"Node,omitempty"`
	Link *link `json:"Link,omitempty"`
}

// A node, which e.g. can be a namespace, type member, ...
type node struct {
	Name        string
	Parent      *string `json:"Parent,omitempty"`
	Type        string
	Description *string `json:"Description,omitempty"`
}

// A link between two nodes
type link struct {
	Source     string
	Target     string
	TargetType string
}

// The node and link items to be written to json file
type modelItems []item

// Creates and link items list to be written to json file
func newModelItems() modelItems {
	return []item{}
}

// Adds a node
func (i modelItems) addNode(name, nodeType, parent, description string) modelItems {
	return append(i, item{
		Node: &node{
			Name:        name,
			Type:        nodeType,
			Parent:      optional(parent),
			Description: optional(strings.TrimSpace(description)),
		},
	})
}

// Adds a link
func (i modelItems) addLink(source, target, targetType string) modelItems {
	return append(i, item{
		Link: &link{
			Source:     source,
			Target:     target,
			TargetType: targetType,
		},
	})
}

func writeModelFile(modelFilePath string, modelItems modelItems) error {
	model := model{
		FormatVersion: "4",
		Items:         modelItems,
	}

	modelBytes, _ := json.MarshalIndent(model, "", "  ")

	err := ioutil.WriteFile(modelFilePath, modelBytes, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Model items: ", len(model.Items))
	return nil
}

// Used when empty text should be exluded froom json objects (as nil)
func optional(text string) *string {
	var textPtr *string
	if text == "" {
		textPtr = nil
	} else {
		textPtr = &text
	}

	return textPtr
}
