package plugins

import (
	"github.com/go-errors/errors"
	v1 "github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"reflect"
	"strings"
)

type Plugin interface {
	ProcessSpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error
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
		return errors.Errorf("plugin %s has no Process function", pluginName)
	}

	var in []reflect.Value
	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}

	funcRef.Call(in)
	return nil
}

func ProcessSpecModel(pluginsList string, specModel *libopenapi.DocumentModel[v3high.Document]) error {
	if pluginsList == "" {
		return nil
	}

	var err error
	pluginPackages := strings.Split(pluginsList, ",")
	for _, pluginName := range pluginPackages {
		if pluginName == "" {
			continue
		}
		if err = InvokeFunc(pluginName, "ProcessSpecModel", specModel); err != nil {
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
