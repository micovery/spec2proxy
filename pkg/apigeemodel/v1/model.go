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

package v1

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"gopkg.in/yaml.v3"
)

type Extension struct {
	Name  string
	Value *yaml.Node
}

type Resource struct {
	ResourceType string
	ResourceFile string
}

type ProxyEndpoint struct {
	Name                string
	BasePath            string
	PreFlow             UnconditionalFlow
	Flows               []*ConditionalFlow
	PostFlow            UnconditionalFlow
	RouteRules          []RouteRule
	HTTPProxyConnection HTTPProxyConnection
	SecurityRequirement []*base.SecurityRequirement
	Extensions          map[string]Extension
}

type SSLInfo struct {
	Enabled                bool
	Enforce                bool
	ClientAuthEnabled      bool
	KeyStore               string
	KeyAlias               string
	TrustStore             string
	IgnoreValidationErrors bool
}

type Property struct {
	Name  string
	Value string
}

type HTTPProxyConnection struct {
	BasePath   string
	Properties []Property
}

type TargetServer struct {
	Name       string
	IsFallback bool
}

type LoadBalancer struct {
	Algorithm string
	Servers   []TargetServer
}
type HTTPTargetConnection struct {
	URL          string
	LoadBalancer LoadBalancer
	SSLInfo      SSLInfo
	Properties   []Property
}

type RouteRule struct {
	Name           string
	Condition      string
	TargetEndpoint string
}

type TargetEndpoint struct {
	Name                 string
	Description          string
	PreFlow              UnconditionalFlow
	Flows                []*ConditionalFlow
	PostFlow             UnconditionalFlow
	HTTPTargetConnection HTTPTargetConnection
	Extensions           map[string]Extension
}

type ConditionalFlow struct {
	Name                string
	Description         string
	Condition           string
	Request             []Step `json:"Request "yaml:"Request"`
	Response            []Step `json:"Response" yaml:"Response"`
	Extensions          map[string]Extension
	SecurityRequirement []*base.SecurityRequirement
}

type UnconditionalFlow struct {
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
	Request     []Step `json:"Request" yaml:"Request"`
	Response    []Step `json:"Response" yaml:"Response"`
	Extensions  map[string]Extension
}

type Step struct {
	Step struct {
		Name      string `json:"Name" yaml:"Name"`
		Condition string `json:"Condition" yaml:"Condition"`
	} `json:"Step" yaml:"Step"`
}

type APIProxy struct {
	Name            string
	Description     string
	DisplayName     string
	CreatedAt       int64
	LastModified    int64
	Policies        []Policy
	ProxyEndpoints  []ProxyEndpoint
	TargetEndpoints []TargetEndpoint
	Resources       []Resource
	Extensions      map[string]Extension
}
