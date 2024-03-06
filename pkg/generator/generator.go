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

package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/micovery/spec2proxy/pkg/apigeemodel/v1"
	"github.com/micovery/spec2proxy/pkg/templates"
	"os"
	"path/filepath"
	"text/template"
)

func Generate(apiProxy *v1.APIProxy, outputDir string) error {
	var manifestBytes []byte
	var proxyEndpointsBytes [][]byte
	var targetEndpointsBytes [][]byte
	var err error

	if manifestBytes, err = generateManifest(apiProxy); err != nil {
		return err
	}

	if proxyEndpointsBytes, err = generateProxyEndpoints(apiProxy); err != nil {
		return err
	}

	if targetEndpointsBytes, err = generateTargetEndpoints(apiProxy); err != nil {
		return err
	}

	var apiProxyDirPath string
	if apiProxyDirPath, err = generateDirectoryStructure(err, outputDir); err != nil {
		return err
	}

	//generate main manifest file
	if err = os.WriteFile(filepath.Join(apiProxyDirPath, fmt.Sprintf("%s.xml", apiProxy.Name)), manifestBytes, os.ModePerm); err != nil {
		return errors.New(err)
	}

	//generate proxy endpoint files
	for index, proxyEndpointBytes := range proxyEndpointsBytes {
		fileName := filepath.Join(apiProxyDirPath, "proxies", fmt.Sprintf("%s.xml", apiProxy.ProxyEndpoints[index].Name))
		if err = os.WriteFile(fileName, proxyEndpointBytes, os.ModePerm); err != nil {
			return errors.New(err)
		}
	}

	//generate target endpoint files
	for index, targetEndpointBytes := range targetEndpointsBytes {
		fileName := filepath.Join(apiProxyDirPath, "targets", fmt.Sprintf("%s.xml", apiProxy.TargetEndpoints[index].Name))
		if err = os.WriteFile(fileName, targetEndpointBytes, os.ModePerm); err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func generateDirectoryStructure(err error, outputDir string) (string, error) {
	var stat os.FileInfo
	if stat, err = os.Stat(outputDir); err != nil {
		//create output dir if it does not exist
		if err = os.Mkdir(outputDir, os.ModePerm); err != nil {
			return "", errors.New(err)
		}
	} else if !stat.IsDir() {
		return "", errors.Errorf("%s is not a directory", outputDir)
	}

	//create directory structure
	apiProxyDirPath := filepath.Join(outputDir, "apiproxy")
	targetsDirPath := filepath.Join(apiProxyDirPath, "targets")
	proxiesDirPath := filepath.Join(apiProxyDirPath, "proxies")
	policiesDirPath := filepath.Join(apiProxyDirPath, "policies")
	resourcesDirPath := filepath.Join(apiProxyDirPath, "resources")

	for _, dir := range []string{apiProxyDirPath, targetsDirPath, proxiesDirPath, policiesDirPath, resourcesDirPath} {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", errors.New(err)
		}
	}

	return apiProxyDirPath, nil
}

func generateTargetEndpoints(apiProxy *v1.APIProxy) ([][]byte, error) {
	var targetEndpointsBytes [][]byte
	for _, targetEndpoint := range apiProxy.TargetEndpoints {
		proxyEndpointBytes, err := generateTextFromTemplate(&targetEndpoint, "target-endpoint.xml.tmpl")
		if err != nil {
			return nil, err
		}
		targetEndpointsBytes = append(targetEndpointsBytes, proxyEndpointBytes)
	}

	return targetEndpointsBytes, nil
}

func generateProxyEndpoints(apiProxy *v1.APIProxy) ([][]byte, error) {
	var proxyEndpointsBytes [][]byte
	for _, proxyEndpoint := range apiProxy.ProxyEndpoints {
		proxyEndpointBytes, err := generateTextFromTemplate(&proxyEndpoint, "proxy-endpoint.xml.tmpl")
		if err != nil {
			return nil, err
		}
		proxyEndpointsBytes = append(proxyEndpointsBytes, proxyEndpointBytes)
	}

	return proxyEndpointsBytes, nil
}

func generateTextFromTemplate(source interface{}, templateFile string) ([]byte, error) {
	var err error
	var endpointBytes bytes.Buffer
	var endpointTmpl *template.Template
	if endpointTmpl, err = getEmbeddedTemplate(templateFile); err != nil {
		return nil, err
	}

	if err = endpointTmpl.Execute(&endpointBytes, source); err != nil {
		return nil, errors.New(err)
	}

	return endpointBytes.Bytes(), nil
}

func generateManifest(apiProxy *v1.APIProxy) ([]byte, error) {
	return generateTextFromTemplate(&apiProxy, "manifest.xml.tmpl")
}

func getEmbeddedTemplate(templatePath string) (*template.Template, error) {
	var err error
	//var templateBytes []byte
	//var templateText string
	var outputTemplate *template.Template

	if outputTemplate, err = template.New(templatePath).ParseFS(templates.FS, "common/*", templatePath); err != nil {
		return nil, errors.New(err)
	}

	return outputTemplate, nil
}
