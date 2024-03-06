package plugins

import (
	"github.com/micovery/spec2proxy/pkg/plugins"
	"github.com/micovery/spec2proxy/plugins/example"
)

func init() {
	// register an instance of the Plugin with the generator
	plugins.RegisterPlugin(&example.Plugin{})
}
