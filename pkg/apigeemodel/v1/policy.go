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
	"github.com/go-errors/errors"
	"github.com/micovery/spec2proxy/pkg/apigeemodel/v1/policies/flowcallout"
	"gopkg.in/yaml.v3"
	"reflect"
)

type PolicyI interface {
	Enabled() bool
	Name() string
	XML() string
}

type Policy struct {
	Policy PolicyI
}

func (p *Policy) Enabled() bool {
	return p.Policy.Enabled()
}

func (p *Policy) Name() string {
	return p.Policy.Name()
}

func (p *Policy) XML() string {
	return p.Policy.XML()
}

func (p *Policy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var policy map[string]any
	var err error
	var ok bool

	if err = unmarshal(&policy); err != nil {
		return err
	}

	if len(policy) == 0 {
		return errors.Errorf("could not unmarshal empty policy")
	}

	var key string
	for key = range policy {
		break
	}

	var dataType reflect.Type
	if dataType, ok = TypesMap[key]; !ok {
		return errors.Errorf("could not find policy type %s", key)
	}

	var policyYAML []byte
	if policyYAML, err = yaml.Marshal(policy); err != nil {
		return errors.Errorf("could not marshall policy %s", key)
	}

	actualPolicy := reflect.New(dataType).Interface().(PolicyI)
	yaml.Unmarshal(policyYAML, actualPolicy)

	(*p).Policy = actualPolicy

	return nil
}

func init() {
	InitTypes()
}

var TypesMap map[string]reflect.Type

func InitTypes() error {

	policyTypes := []reflect.Type{
		reflect.TypeOf(flowcallout.FlowCallout{}),
	}

	TypesMap = make(map[string]reflect.Type)
	for _, t := range policyTypes {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if policyTag := field.Tag.Get("policy"); policyTag != "true" {
				continue
			}

			TypesMap[field.Name] = t
		}
	}

	return nil
}
