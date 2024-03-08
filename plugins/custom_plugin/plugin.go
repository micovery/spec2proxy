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

package custom_plugin

import (
	"fmt"
	v1 "github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
)

type Visibility struct {
	Extent string
}

func isInternal(node *yaml.Node) bool {
	if node == nil {
		return false
	}
	var visibility Visibility
	node.Decode(&visibility)
	return strings.EqualFold(visibility.Extent, "internal")
}

func unsetOperation(pathItem *v3high.PathItem, opKey string) {
	elem := strings.ToTitle(opKey[0:1]) + opKey[1:]
	field := reflect.ValueOf(pathItem).Elem().FieldByName(elem)
	field.Set(reflect.Zero(field.Type()))
}

func (p *Plugin) ProcessSpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error {
	for path := specModel.Model.Paths.PathItems.First(); path != nil; path = path.Next() {
		pathKey := path.Key()
		pathItem := path.Value()
		operations := pathItem.GetOperations()
		for operation := operations.First(); operation != nil; operation = operation.Next() {
			opKey := operation.Key()
			opItem := operation.Value()
			if visibilityExtension, found := (opItem.Extensions.Get("x-visibility")); found && isInternal(visibilityExtension) {
				fmt.Printf("Removing internal operation: %s %s", strings.ToUpper(opKey), pathKey)
				unsetOperation(pathItem, opKey)
			}
		}
	}

	return nil
}

// Plugin Custom plugin for handling "x-visibility" OpenAPI extension
type Plugin struct {
}

func (p *Plugin) ProcessProxyModel(apiProxy *v1.APIProxy) error {
	// this is chance to modify the Apigee API Proxy model before it gets generated to a bundle on disk

	return nil
}
