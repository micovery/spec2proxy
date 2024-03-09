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

package example

import (
	v1 "github.com/micovery/spec2proxy/pkg/apigee/v1"
	"github.com/pb33f/libopenapi"
	v2high "github.com/pb33f/libopenapi/datamodel/high/v2"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Plugin Sample plugin
type Plugin struct {
}

func (p *Plugin) ProcessOAS3SpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error {
	// this is your chance to modify the OpenAPI spec after it has been parsed
	return nil
}

func (p *Plugin) ProcessOAS2SpecModel(specModel *libopenapi.DocumentModel[v2high.Swagger]) error {
	// this is your chance to modify the OpenAPI spec after it has been parsed
	return nil
}
func (p *Plugin) ProcessProxyModel(apiProxy *v1.APIProxy) error {
	// this is chance to modify the Apigee API Proxy model before it gets generated to a bundle on disk

	return nil
}
