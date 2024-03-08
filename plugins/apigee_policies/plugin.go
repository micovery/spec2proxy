// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apigee_policies

import (
	"encoding/xml"
	"github.com/go-errors/errors"
	v1 "github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/index"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func (p *Plugin) ProcessSpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error {
	return nil
}

// Plugin Custom plugin for handling "x-visibility" OpenAPI extension
type Plugin struct {
}

func ResolveReferences(specModel *libopenapi.DocumentModel[v3high.Document]) {
	basePath := "."

	// create an index config
	config := index.CreateOpenAPIIndexConfig()

	// the rolodex will automatically try and check for circular references, you don't want to do this
	// if you're resolving the spec, as the node tree is marked as 'seen' and you won't be able to resolve
	// correctly.
	config.AvoidCircularReferenceCheck = true
	config.AllowFileLookup = true

	// new in 0.13+ is the ability to add remote and local file systems to the index
	// requires a new part, the rolodex. It holds all the indexes and knows where to find
	// every reference across local and remote files.
	rolodex := index.NewRolodex(config)

	// create a local file system config, tell it where to look from and the index config to pay attention to.
	fsCfg := &index.LocalFSConfig{
		BaseDirectory: basePath,
		IndexConfig:   config,
	}

	// create a local file system using config.
	fileFS, err := index.NewLocalFSWithConfig(fsCfg)
	if err != nil {
		panic(err)
	}

	// unmarshal the spec into a yaml node
	var rootNode = specModel.Index.GetRootNode()

	// set the root node of the rolodex (this is the root of the spec)
	rolodex.SetRootNode(rootNode)

	// add local file system to rolodex
	rolodex.AddLocalFS(basePath, fileFS)

	// add a new remote file system.
	remoteFS, _ := index.NewRemoteFSWithConfig(config)

	// add the remote file system to the rolodex
	rolodex.AddRemoteFS("", remoteFS)

	// set the root node of the rolodex, this is your spec.
	rolodex.SetRootNode(rootNode)

	// index the rolodex
	indexingError := rolodex.IndexTheRolodex()
	if indexingError != nil {
		panic(indexingError)
	}

	// resolve the rolodex (if you want to)
	rolodex.Resolve()

	// there should be no errors at this point
	resolvingErrors := rolodex.GetCaughtErrors()
	if resolvingErrors != nil {
		panic(resolvingErrors)
	}
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func (p *Plugin) ProcessProxyModel(apiProxy *v1.APIProxy) error {

	var err error

	// handle policies
	if err = UnmarshalExtension("x-Apigee-Policies", apiProxy.Extensions, &apiProxy.Policies); err != nil {
		return err
	}

	// handle PostFlow
	if err = UnmarshalExtension("x-Apigee-PostFlow", apiProxy.Extensions, &apiProxy.ProxyEndpoints[0].PostFlow); err != nil {
		return err
	}

	// handle PreFlow
	if err = UnmarshalExtension("x-Apigee-PreFlow", apiProxy.Extensions, &apiProxy.ProxyEndpoints[0].PreFlow); err != nil {
		return err
	}

	// handle conditional flows
	for _, proxyEndpoint := range apiProxy.ProxyEndpoints {
		for _, conditionalFlow := range proxyEndpoint.Flows {
			if err = UnmarshalExtension("x-Apigee-Flow", conditionalFlow.Extensions, conditionalFlow); err != nil {
				return err
			}
		}
	}

	return nil
}

func UnmarshalExtension(extensionName string, extensions map[string]v1.Extension, target any) error {
	var rawExtension v1.Extension
	var ok bool
	var err error

	if rawExtension, ok = extensions[extensionName]; !ok {
		return nil
	}

	if rawExtension.Value, err = ResolveYAMLRefs(rawExtension.Value); err != nil {
		return errors.New(err)
	}

	var yamlText []byte
	if yamlText, err = yaml.Marshal(rawExtension.Value); err != nil {
		return errors.New(err)
	}

	yaml.Unmarshal(yamlText, target)

	return err
}

func isYAMLRef(node *yaml.Node) bool {
	if node == nil {
		return false
	}

	return node.Kind == yaml.MappingNode &&
		len(node.Content) == 2 &&
		node.Content[0].Value == "$ref"
}

func getYAMLRefLocation(node *yaml.Node) string {
	return node.Content[1].Value
}

var ParsedYAMLFiles map[string]*yaml.Node

func ParseYAMLFile(filePath string) (*yaml.Node, error) {
	var fileBytes []byte
	var err error
	var ok bool
	var rootNode *yaml.Node

	if rootNode, ok = ParsedYAMLFiles[filePath]; ok {
		return rootNode, nil
	}

	if fileBytes, err = os.ReadFile(filePath); err != nil {
		return nil, err
	}

	rootNode = &yaml.Node{}

	if err = yaml.Unmarshal(fileBytes, rootNode); err != nil {
		return nil, err
	}

	var resolvedNode *yaml.Node
	if resolvedNode, err = ResolveYAMLRefs(rootNode); err != nil {
		return nil, err
	}

	ParsedYAMLFiles[filePath] = resolvedNode
	return rootNode, nil
}

func ConvertJSONRef2YAMLPath(jsonRef string) string {

	yamlPath := strings.ReplaceAll(jsonRef, "/", ".")

	return "$" + yamlPath
}

func ResolveYAMLRef(location string) (*yaml.Node, error) {
	var err error

	locationParts := strings.Split(location, "#")

	if len(locationParts) != 2 {
		return nil, errors.Errorf("JSONRef '%s' is not valid", location)
	}

	filePath := locationParts[0]
	docPath := locationParts[1]

	if filePath == "" {
		return nil, errors.Errorf("self referncing JSONRef '%s' is not supported", location)
	}

	var fileRootNode *yaml.Node
	if fileRootNode, err = ParseYAMLFile(filePath); err != nil {
		return nil, err
	}

	convertedPath := ConvertJSONRef2YAMLPath(docPath)
	var yamlPath *yamlpath.Path
	if yamlPath, err = yamlpath.NewPath(convertedPath); err != nil {
		return nil, err
	}

	var yamlNodes []*yaml.Node
	if yamlNodes, err = yamlPath.Find(fileRootNode); err != nil {
		return nil, err
	}

	if len(yamlNodes) == 0 {
		return nil, errors.Errorf("no node found at JSONRef '%s'", location)
	}

	if len(yamlNodes) > 1 {
		return nil, errors.Errorf("more than one node found at JSONRef '%s'", location)
	}

	return yamlNodes[0], nil
}

func ResolveYAMLRefs(node *yaml.Node) (*yaml.Node, error) {
	if node == nil {
		return nil, nil
	}

	var resolvedNode *yaml.Node
	var err error

	if node.Kind == yaml.MappingNode && isYAMLRef(node) {

		location := getYAMLRefLocation(node)
		if resolvedNode, err = ResolveYAMLRef(location); err != nil {
			return nil, err
		}
		return resolvedNode, nil
	} else if node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			if resolvedNode, err = ResolveYAMLRefs(node.Content[i+1]); err != nil {
				return nil, err
			}
			node.Content[i+1] = resolvedNode
		}
		return node, nil
	} else if node.Kind == yaml.SequenceNode {
		for i := 0; i+1 < len(node.Content); i += 1 {
			if resolvedNode, err = ResolveYAMLRefs(node.Content[i]); err != nil {
				return nil, err
			}
			node.Content[i] = resolvedNode
		}
		return node, nil
	}

	return node, nil

}

func init() {
	ParsedYAMLFiles = make(map[string]*yaml.Node)
}
