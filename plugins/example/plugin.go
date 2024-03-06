package example

import (
	v1 "github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Plugin Sample plugin
type Plugin struct {
}

func (p *Plugin) ProcessSpecModel(specModel *libopenapi.DocumentModel[v3high.Document]) error {
	// this is your chance to modify the OpenAPI spec after it has been parsed
	return nil
}

func (p *Plugin) ProcessProxyModel(apiProxy *v1.APIProxy) error {
	// this is chance to modify the Apigee API Proxy model before it gets generated to a bundle on disk

	return nil
}
