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
	"github.com/micovery/spec2proxy/pkg/apigee/v1"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
	"regexp"
)

func GetExtensions(source *orderedmap.Map[string, *yaml.Node]) map[string]*v1.Extension {
	result := make(map[string]*v1.Extension)
	elem := source.First()

	for elem != nil {
		key := elem.Key()
		value := elem.Value()
		result[key] = &v1.Extension{
			Name:  key,
			Value: value,
		}
		elem = elem.Next()
	}

	return result
}

func AppendExtensions(target map[string]*v1.Extension, source *orderedmap.Map[string, *yaml.Node]) {
	elem := source.First()
	for elem != nil {
		key := elem.Key()
		value := elem.Value()
		target[key] = &v1.Extension{
			Name:  key,
			Value: value,
		}
		elem = elem.Next()
	}
}

func SetupRouteRules(apiProxy *v1.APIProxy, proxyEndpoint *v1.ProxyEndpoint, targetEndpoint *v1.TargetEndpoint) {
	proxyEndpoint.RouteRules = append(proxyEndpoint.RouteRules, &v1.RouteRule{
		Name:           "default",
		TargetEndpoint: targetEndpoint.Name,
		Condition:      "true",
	})

	apiProxy.ProxyEndpoints = append(apiProxy.ProxyEndpoints, proxyEndpoint)
	apiProxy.TargetEndpoints = append(apiProxy.TargetEndpoints, targetEndpoint)
	apiProxy.Resources = []*v1.Resource{}
}

func ToApigeePath(oasPath string) string {
	re := regexp.MustCompile(`{[^}]*}`)
	return re.ReplaceAllString(oasPath, "*")
}
