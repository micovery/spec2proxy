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

package v2

import (
	"fmt"
	"github.com/gosimple/slug"
	v1 "github.com/micovery/spec2proxy/pkg/apigee/v1"
	"github.com/micovery/spec2proxy/pkg/transformer"
	"github.com/pb33f/libopenapi"
	v2high "github.com/pb33f/libopenapi/datamodel/high/v2"
	"net/url"
	"strings"
	"time"
)

func Transform(specModel *libopenapi.DocumentModel[v2high.Swagger]) (*v1.APIProxy, error) {
	var err error
	var targetEndpoint *v1.TargetEndpoint
	var proxyEndpoint *v1.ProxyEndpoint

	apiProxy := v1.APIProxy{}

	now := time.Now().UnixMilli()
	apiProxy.Name = slug.Make(specModel.Model.Info.Title)
	apiProxy.DisplayName = specModel.Model.Info.Title
	apiProxy.Description = specModel.Model.Info.Description
	apiProxy.Extensions = transformer.GetExtensions(specModel.Model.Extensions)
	transformer.AppendExtensions(apiProxy.Extensions, specModel.Model.Info.Extensions)
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
	//link proxy endpoint to target endpoint with route rule
	transformer.SetupRouteRules(&apiProxy, proxyEndpoint, targetEndpoint)

	return &apiProxy, nil
}

func buildTargetEndpoint(specModel *libopenapi.DocumentModel[v2high.Swagger]) (*v1.TargetEndpoint, error) {
	var targetEndpoint v1.TargetEndpoint
	targetEndpoint.Name = "default"
	targetEndpoint.Flows = []*v1.ConditionalFlow{}
	targetEndpoint.PreFlow = &v1.UnconditionalFlow{
		Request:    []*v1.Step{},
		Response:   []*v1.Step{},
		Extensions: nil,
	}

	targetEndpoint.PostFlow = &v1.UnconditionalFlow{
		Request:    []*v1.Step{},
		Response:   []*v1.Step{},
		Extensions: nil,
	}

	targetEndpoint.HTTPTargetConnection = &v1.HTTPTargetConnection{
		URL: extractTargetEndpointUrl(specModel),
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
		return nil, err
	}

	if parsedUrl.Scheme == "http" {
		targetEndpoint.HTTPTargetConnection.SSLInfo.Enabled = false
	}

	return &targetEndpoint, nil
}

func buildProxyEndpoint(specModel *libopenapi.DocumentModel[v2high.Swagger]) (*v1.ProxyEndpoint, error) {
	var proxyEndpoint v1.ProxyEndpoint

	proxyEndpoint.BasePath = specModel.Model.BasePath

	proxyEndpoint.Name = "default"
	proxyEndpoint.PreFlow = &v1.UnconditionalFlow{
		Request:    []*v1.Step{},
		Response:   []*v1.Step{},
		Extensions: nil,
	}

	proxyEndpoint.PostFlow = &v1.UnconditionalFlow{
		Request:    []*v1.Step{},
		Response:   []*v1.Step{},
		Extensions: nil,
	}

	proxyEndpoint.RouteRules = []*v1.RouteRule{}

	proxyEndpoint.SecurityRequirement = specModel.Model.Security
	proxyEndpoint.Extensions = transformer.GetExtensions(specModel.Model.Paths.Extensions)

	appendConditionalFlows(&proxyEndpoint, specModel.Model.Paths)

	return &proxyEndpoint, nil
}

func extractTargetEndpointUrl(specModel *libopenapi.DocumentModel[v2high.Swagger]) string {

	scheme := "https"
	if len(specModel.Model.Schemes) > 0 {
		scheme = specModel.Model.Schemes[0]
	}

	host := "mocktarget.apigee.net"
	if specModel.Model.Host != "" {
		host = specModel.Model.Host
	}

	url := fmt.Sprintf("%s://%s%s", scheme, host, specModel.Model.BasePath)

	return url
}

func appendConditionalFlows(endpoint *v1.ProxyEndpoint, paths *v2high.Paths) {
	endpoint.Flows = []*v1.ConditionalFlow{}

	for path := paths.PathItems.First(); path != nil; path = path.Next() {
		pathSegment := transformer.ToApigeePath(path.Key())
		pathInfo := path.Value()
		operations := pathInfo.GetOperations()

		for operation := operations.First(); operation != nil; operation = operation.Next() {
			operationKey := operation.Key()
			operationInfo := operation.Value()

			conditionalFlow := &v1.ConditionalFlow{
				Name:                operationInfo.OperationId,
				Description:         operationInfo.Description,
				Condition:           fmt.Sprintf("(proxy.pathsuffix MatchesPath \"%s\") and (request.verb = \"%s\")", pathSegment, strings.ToUpper(operationKey)),
				Request:             []*v1.Step{},
				Response:            []*v1.Step{},
				Extensions:          transformer.GetExtensions(pathInfo.Extensions),
				SecurityRequirement: operationInfo.Security,
			}

			transformer.AppendExtensions(conditionalFlow.Extensions, operationInfo.Extensions)
			endpoint.Flows = append(endpoint.Flows, conditionalFlow)
		}
	}
}
