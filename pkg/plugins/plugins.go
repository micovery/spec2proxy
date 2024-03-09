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

package plugins

import (
	"github.com/go-errors/errors"
	v1 "github.com/micovery/spec2proxy/pkg/apigee/v1"
	"github.com/pb33f/libopenapi"
	v2high "github.com/pb33f/libopenapi/datamodel/high/v2"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"reflect"
	"strings"
)

type Plugin interface {
	ProcessOAS3SpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error

	ProcessOAS2SpecModel(specModel *libopenapi.DocumentModel[v2high.Swagger]) error

	ProcessProxyModel(apiProxy *v1.APIProxy) error
}

var plugins = make([]Plugin, 0)

func RegisterPlugin(plug Plugin) error {
	plugins = append(plugins, plug)
	return nil
}

func InvokeFunc(pluginName string, funcName string, args ...any) error {
	var plugin Plugin
	var found bool

	registeredPlugins := GetRegisteredPlugins()

	if plugin, found = registeredPlugins[pluginName]; !found {
		return errors.Errorf("plugin %s not found", pluginName)
	}

	funcRef := reflect.ValueOf(plugin).MethodByName(funcName)

	if !funcRef.IsValid() {
		return errors.Errorf("plugin %s has no %s function", pluginName, funcName)
	}

	var in []reflect.Value
	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}

	results := funcRef.Call(in)
	if results[0].IsZero() {
		return nil
	}
	return results[0].Interface().(error)
}

func ProcessOAS3SpecModel(pluginsList string, specModel *libopenapi.DocumentModel[v3high.Document]) error {
	if pluginsList == "" {
		return nil
	}

	var err error
	pluginPackages := strings.Split(pluginsList, ",")
	for _, pluginName := range pluginPackages {
		if pluginName == "" {
			continue
		}
		if err = InvokeFunc(pluginName, "ProcessOAS3SpecModel", specModel); err != nil {
			return err
		}
	}
	return nil
}

func ProcessOAS2SpecModel(pluginsList string, specModel *libopenapi.DocumentModel[v2high.Swagger]) error {
	if pluginsList == "" {
		return nil
	}

	var err error
	pluginPackages := strings.Split(pluginsList, ",")
	for _, pluginName := range pluginPackages {
		if pluginName == "" {
			continue
		}
		if err = InvokeFunc(pluginName, "ProcessOAS2SpecModel", specModel); err != nil {
			return err
		}
	}
	return nil
}

func ProcessProxyModel(pluginsList string, apiProxy *v1.APIProxy) error {
	if pluginsList == "" {
		return nil
	}

	var err error
	pluginPackages := strings.Split(pluginsList, ",")
	for _, pluginName := range pluginPackages {
		if pluginName == "" {
			continue
		}
		if err = InvokeFunc(pluginName, "ProcessProxyModel", apiProxy); err != nil {
			return err
		}
	}
	return nil
}

func GetRegisteredPlugins() map[string]Plugin {
	pluginsByName := make(map[string]Plugin)
	for _, plugin := range plugins {
		pluginType := reflect.TypeOf(plugin).String()
		pluginPackage := strings.Split(pluginType, ".")[0]
		pluginsByName[pluginPackage[1:]] = plugin
	}

	return pluginsByName
}
