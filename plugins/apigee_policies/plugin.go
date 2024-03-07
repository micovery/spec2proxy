package apigee_policies

import (
	"fmt"
	"github.com/go-errors/errors"
	v1 "github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"gopkg.in/yaml.v3"
)

func (p *Plugin) ProcessSpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error {

	return nil
}

// Plugin Custom plugin for handling "x-visibility" OpenAPI extension
type Plugin struct {
}

func (p *Plugin) ProcessProxyModel(apiProxy *v1.APIProxy) error {

	var err error

	if err = ProcessPoliciesExtension(apiProxy); err != nil {
		return err
	}

	if err = ProcessPreFlowExtension(apiProxy); err != nil {
		return err
	}

	return nil
}

func ProcessPreFlowExtension(apiProxy *v1.APIProxy) error {

	var rawPreFlow v1.Extension
	var ok bool
	var err error

	if rawPreFlow, ok = apiProxy.Extensions["x-Apigee-PreFlow"]; !ok {
		return nil
	}

	var yamlText []byte
	if yamlText, err = yaml.Marshal(rawPreFlow.Value); err != nil {
		return errors.New(err)
	}

	yaml.Unmarshal(yamlText, &apiProxy.ProxyEndpoints[0].PreFlow)

	fmt.Printf("%v\n", yamlText)

	return err
}

func ProcessPoliciesExtension(apiProxy *v1.APIProxy) error {
	var rawPolicies v1.Extension
	var ok bool
	var err error

	if rawPolicies, ok = apiProxy.Extensions["x-Apigee-Policies"]; !ok {
		return nil
	}

	var yamlText []byte
	if yamlText, err = yaml.Marshal(rawPolicies.Value); err != nil {
		return errors.New(err)
	}

	yaml.Unmarshal(yamlText, &apiProxy.Policies)
	return nil
}
