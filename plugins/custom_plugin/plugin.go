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
	v1 "github.com/micovery/spec2proxy/pkg/apigee/v1"
	"github.com/pb33f/libopenapi"
	v2high "github.com/pb33f/libopenapi/datamodel/high/v2"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"gopkg.in/yaml.v3"
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

func (p *Plugin) ProcessOAS3SpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error {
	return nil
}

func (p *Plugin) ProcessOAS2SpecModel(specModel *libopenapi.DocumentModel[v2high.Swagger]) error {
	return nil
}

// Plugin Custom plugin for handling "x-visibility" OpenAPI extension
type Plugin struct {
}

func (p *Plugin) ProcessProxyModel(apiProxy *v1.APIProxy) error {
	if len(apiProxy.ProxyEndpoints) == 0 {
		return nil
	}

	var newFlows []*v1.ConditionalFlow

	for _, flow := range apiProxy.ProxyEndpoints[0].Flows {
		if isInternalFlow(flow) {
			fmt.Printf("Removing internal operation: %s\n", flow.Name)
			continue
		}
		newFlows = append(newFlows, flow)
	}

	apiProxy.ProxyEndpoints[0].Flows = newFlows

	//TODO: Add catch-all flow at the end

	return nil
}

func isInternalFlow(flow *v1.ConditionalFlow) bool {
	for _, extension := range flow.Extensions {
		if extension.Name == "x-visibility" && isInternal(extension.Value) {
			return true
		}
	}
	return false
}
