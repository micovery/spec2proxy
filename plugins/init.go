package plugins

import (
	"github.com/micovery/spec2proxy/pkg/plugins"
	"github.com/micovery/spec2proxy/plugins/apigee_policies"
	"github.com/micovery/spec2proxy/plugins/custom_plugin"
	"github.com/micovery/spec2proxy/plugins/example"
)

func init() {
	// register an instance of the Plugin with the generator
	plugins.RegisterPlugin(&apigee_policies.Plugin{})
	plugins.RegisterPlugin(&custom_plugin.Plugin{})
	plugins.RegisterPlugin(&example.Plugin{})
}
