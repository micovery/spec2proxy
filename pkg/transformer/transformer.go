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

package transformer

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
	"net/url"
	"strings"
	"time"
)

func Transform(specModel *libopenapi.DocumentModel[v3high.Document]) (*v1.APIProxy, error) {
	var err error
	var targetEndpoint v1.TargetEndpoint
	var proxyEndpoint v1.ProxyEndpoint

	apiProxy := v1.APIProxy{}

	now := time.Now().UnixMilli()
	apiProxy.Name = slug.Make(specModel.Model.Info.Title)
	apiProxy.DisplayName = specModel.Model.Info.Title
	apiProxy.Description = specModel.Model.Info.Description
	apiProxy.Extensions = getExtensions(specModel.Model.Extensions)
	appendExtensions(apiProxy.Extensions, specModel.Model.Info.Extensions)
	apiProxy.CreatedAt = now
	apiProxy.LastModified = now

	//build proxy endpoint
	if proxyEndpoint, err = buildProxyEndpoint(specModel); err != nil {
		return nil, err
	}

	//build target endpoint
	if targetEndpoint, err = buildTargetEndpoint(specModel); err != nil {
		return nil, err
	}

	//link proxy endpoint to target endpoint with route rule
	proxyEndpoint.RouteRules = append(proxyEndpoint.RouteRules, v1.RouteRule{
		Name:           "default",
		TargetEndpoint: targetEndpoint.Name,
		Condition:      "true",
	})

	apiProxy.ProxyEndpoints = append(apiProxy.ProxyEndpoints, proxyEndpoint)
	apiProxy.TargetEndpoints = append(apiProxy.TargetEndpoints, targetEndpoint)
	apiProxy.Resources = make([]v1.Resource, 0)

	return &apiProxy, nil
}

func buildTargetEndpoint(specModel *libopenapi.DocumentModel[v3high.Document]) (v1.TargetEndpoint, error) {
	var targetEndpoint v1.TargetEndpoint
	targetEndpoint.Name = "default"
	targetEndpoint.Flows = make([]v1.ConditionalFlow, 0)
	targetEndpoint.PreFlow = v1.UnconditionalFlow{
		RequestSteps:  make([]v1.Step, 0),
		ResponseSteps: make([]v1.Step, 0),
		Extensions:    nil,
	}

	targetEndpoint.PostFlow = v1.UnconditionalFlow{
		RequestSteps:  make([]v1.Step, 0),
		ResponseSteps: make([]v1.Step, 0),
		Extensions:    nil,
	}

	targetEndpoint.HTTPTargetConnection = v1.HTTPTargetConnection{
		URL: extractBaseTargetUrl(specModel),
		SSLInfo: v1.SSLInfo{
			Enabled:                true,
			Enforce:                false,
			IgnoreValidationErrors: true,
		},
		Properties: nil,
	}

	var parsedUrl *url.URL
	var err error
	if parsedUrl, err = url.Parse(targetEndpoint.HTTPTargetConnection.URL); err != nil {

	}
	if parsedUrl.Scheme == "http" {
		targetEndpoint.HTTPTargetConnection.SSLInfo.Enabled = false
	}

	return targetEndpoint, nil
}

func buildProxyEndpoint(specModel *libopenapi.DocumentModel[v3high.Document]) (v1.ProxyEndpoint, error) {
	var proxyEndpoint v1.ProxyEndpoint

	proxyEndpoint.BasePath = extractBasePath(specModel)

	proxyEndpoint.Name = "default"
	proxyEndpoint.PreFlow = v1.UnconditionalFlow{
		RequestSteps:  make([]v1.Step, 0),
		ResponseSteps: make([]v1.Step, 0),
		Extensions:    nil,
	}

	proxyEndpoint.PostFlow = v1.UnconditionalFlow{
		RequestSteps:  make([]v1.Step, 0),
		ResponseSteps: make([]v1.Step, 0),
		Extensions:    nil,
	}

	proxyEndpoint.RouteRules = make([]v1.RouteRule, 0)

	proxyEndpoint.SecurityRequirement = specModel.Model.Security
	proxyEndpoint.Extensions = getExtensions(specModel.Model.Paths.Extensions)

	appendConditionalFlows(&proxyEndpoint, specModel.Model.Paths)

	return proxyEndpoint, nil
}

func extractBaseTargetUrl(specModel *libopenapi.DocumentModel[v3high.Document]) string {
	if len(specModel.Model.Servers) == 0 {
		return "https://mocktarget.apigee.net"
	}

	firstServer := specModel.Model.Servers[0]

	//parse the URL to make sure it's valid
	_, err := url.Parse(firstServer.URL)
	if err != nil {
		return "https://mocktarget.apigee.net"
	}

	return firstServer.URL
}

func extractBasePath(specModel *libopenapi.DocumentModel[v3high.Document]) string {
	if len(specModel.Model.Servers) == 0 {
		return "/"
	}

	firstServer := specModel.Model.Servers[0]

	url, err := url.Parse(firstServer.URL)
	if err != nil {
		return "/"
	}

	return url.Path
}

func appendConditionalFlows(endpoint *v1.ProxyEndpoint, paths *v3high.Paths) {
	endpoint.Flows = make([]v1.ConditionalFlow, 0)

	for path := paths.PathItems.First(); path != nil; path = path.Next() {
		pathSegment := path.Key()
		pathInfo := path.Value()
		operations := pathInfo.GetOperations()

		for operation := operations.First(); operation != nil; operation = operation.Next() {
			operationKey := operation.Key()
			operationInfo := operation.Value()

			conditionalFlow := v1.ConditionalFlow{
				Name:                operationInfo.OperationId,
				Description:         operationInfo.Description,
				Condition:           fmt.Sprintf("(proxy.pathsuffix MatchesPath \"%s\") and (request.verb = \"%s\")", pathSegment, strings.ToUpper(operationKey)),
				RequestSteps:        make([]v1.Step, 0),
				ResponseSteps:       make([]v1.Step, 0),
				Extensions:          getExtensions(pathInfo.Extensions),
				SecurityRequirement: operationInfo.Security,
			}

			appendExtensions(conditionalFlow.Extensions, operationInfo.Extensions)
			endpoint.Flows = append(endpoint.Flows, conditionalFlow)
		}
	}
}

func appendExtensions(target map[string]v1.Extension, source *orderedmap.Map[string, *yaml.Node]) {
	elem := source.First()
	for elem != nil {
		key := elem.Key()
		value := elem.Value()
		target[key] = v1.Extension{
			Name:  key,
			Value: value,
		}
		elem = elem.Next()
	}
}

func getExtensions(source *orderedmap.Map[string, *yaml.Node]) map[string]v1.Extension {
	result := make(map[string]v1.Extension)
	elem := source.First()

	for elem != nil {
		key := elem.Key()
		value := elem.Value()
		result[key] = v1.Extension{
			Name:  key,
			Value: value,
		}
		elem = elem.Next()
	}

	return result
}
