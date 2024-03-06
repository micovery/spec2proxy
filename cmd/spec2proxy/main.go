package main

import (
	"flag"
	"fmt"
	"github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/micovery/spec2proxy/pkg/generator"
	"github.com/micovery/spec2proxy/pkg/parser"
	"github.com/micovery/spec2proxy/pkg/transformer"
	"github.com/micovery/spec2proxy/pkg/utils"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func main() {
	var specFile string
	var outputDir string
	var pluginsList string

	var err error
	var specModel *libopenapi.DocumentModel[v3high.Document]
	var apiModel *v1.APIProxy

	flag.StringVar(&specFile, "spec", "", "path to OpenAPI spec file. e.g. \"./petstore.yaml\"")
	flag.StringVar(&outputDir, "out", "", "output directory. e.g \"./hello-world\"")
	flag.StringVar(&pluginsList, "plugins", "", "list of plugins. e.g. \"plugin1,plugin2,etc\"")
	flag.Parse()

	if specFile == "" {
		utils.RequireParamAndExit("spec")
	}

	if outputDir == "" {
		utils.RequireParamAndExit("out")
	}

	var errs []error
	if specModel, errs = parser.Parse(specFile); len(errs) != 0 {
		for _, err = range errs {
			fmt.Println(err)
		}

		utils.PrintErrorWithStackAndExit(errs[0])
	}

	if apiModel, err = transformer.Transform(specModel); err != nil {
		panic(err)
	}

	//TODO: add plugin logic

	if err = generator.Generate(apiModel, outputDir); err != nil {
		utils.PrintErrorWithStackAndExit(err)
	}
}
