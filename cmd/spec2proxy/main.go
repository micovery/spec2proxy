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

package main

import (
	"flag"
	"fmt"
	"github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/micovery/spec2proxy/pkg/generator"
	"github.com/micovery/spec2proxy/pkg/parser"
	"github.com/micovery/spec2proxy/pkg/plugins"
	"github.com/micovery/spec2proxy/pkg/transformer"
	"github.com/micovery/spec2proxy/pkg/utils"
	_ "github.com/micovery/spec2proxy/plugins"
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

	flag.StringVar(&specFile, "oas", "", "path to OpenAPI spec file. e.g. \"./petstore.yaml\"")
	flag.StringVar(&outputDir, "out", "", "output directory. e.g \"./hello-world\"")
	flag.StringVar(&pluginsList, "plugins", "", "list of plugins. e.g. \"plugin1,plugin2,etc\"")
	flag.Parse()

	if specFile == "" {
		utils.RequireParamAndExit("oas")
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

	// call plugins to process the OAS spec
	if err = plugins.ProcessSpecModel(pluginsList, specModel); err != nil {
		utils.PrintErrorWithStackAndExit(err)
	}

	if apiModel, err = transformer.Transform(specModel); err != nil {
		utils.PrintErrorWithStackAndExit(err)
	}

	// call plugins to process the Apigee API Proxy model
	if err = plugins.ProcessProxyModel(pluginsList, apiModel); err != nil {
		utils.PrintErrorWithStackAndExit(err)
	}

	if err = generator.Generate(apiModel, outputDir); err != nil {
		utils.PrintErrorWithStackAndExit(err)
	}
}
