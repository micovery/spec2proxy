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
	"github.com/go-errors/errors"
	"github.com/micovery/spec2proxy/pkg/apigee/v1"
	"github.com/micovery/spec2proxy/pkg/generator"
	"github.com/micovery/spec2proxy/pkg/parser"
	"github.com/micovery/spec2proxy/pkg/plugins"
	v2 "github.com/micovery/spec2proxy/pkg/transformer/v2"
	"github.com/micovery/spec2proxy/pkg/transformer/v3"
	"github.com/micovery/spec2proxy/pkg/utils"
	_ "github.com/micovery/spec2proxy/plugins"
	"github.com/pb33f/libopenapi"
	v2high "github.com/pb33f/libopenapi/datamodel/high/v2"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func main() {
	var specFile string
	var outputDir string
	var pluginsList string

	var errs []error
	var err error
	var specModelV2 *libopenapi.DocumentModel[v2high.Swagger]
	var specModelV3 *libopenapi.DocumentModel[v3high.Document]
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

	var spec libopenapi.Document
	if spec, err = parser.Parse(specFile); err != nil {
		fmt.Println(err)
		utils.PrintErrorWithStackAndExit(err)
	}
	specVersion := spec.GetSpecInfo().VersionNumeric

	if specVersion == 2.0 {
		if specModelV2, errs = parser.BuildOAS2Model(spec); len(errs) != 0 {
			for _, err = range errs {
				fmt.Println(err)
			}
			utils.PrintErrorWithStackAndExit(errs[0])
		}
		// call plugins to process the OAS spec
		if err = plugins.ProcessOAS2SpecModel(pluginsList, specModelV2); err != nil {
			utils.PrintErrorWithStackAndExit(err)
		}

		if apiModel, err = v2.Transform(specModelV2); err != nil {
			utils.PrintErrorWithStackAndExit(err)
		}
	} else if specVersion >= 3 {
		if specModelV3, errs = parser.BuildOAS3Model(spec); len(errs) != 0 {
			for _, err = range errs {
				fmt.Println(err)
			}
			utils.PrintErrorWithStackAndExit(errs[0])
		}
		// call plugins to process the OAS spec
		if err = plugins.ProcessOAS3SpecModel(pluginsList, specModelV3); err != nil {
			utils.PrintErrorWithStackAndExit(err)
		}

		if apiModel, err = v3.Transform(specModelV3); err != nil {
			utils.PrintErrorWithStackAndExit(err)
		}
	} else {
		err = errors.Errorf("OpenAPI spec version %s is not supported", spec.GetVersion())
		fmt.Println(err)
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
