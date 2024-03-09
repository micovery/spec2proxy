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
	"gopkg.in/yaml.v3"
	"reflect"
)

type PolicyI interface {
	Name() string
	XML() []byte
}

type Policy struct {
	name string `json:".name" yaml:".name"`
	Data *yaml.Node
}

func (p *Policy) Name() string {
	return p.name
}

func (p *Policy) XML() ([]byte, error) {
	return nil, nil
}

func (p *Policy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var policy map[string]any
	var policyYAML []byte
	var ok bool
	var err error

	if err = unmarshal(&policy); err != nil {
		return err
	}

	if len(policy) == 0 {
		return errors.Errorf("could not unmarshal empty policy")
	}

	if policyYAML, err = yaml.Marshal(policy); err != nil {
		return err
	}

	p.Data = &yaml.Node{}
	if err = yaml.Unmarshal(policyYAML, p.Data); err != nil {
		return err
	}

	var key string
	for key = range policy {
		break
	}

	var content map[string]any
	if content, ok = policy[key].(map[string]any); !ok {
		return errors.Errorf("malformed %s policy", key)
	}

	var name string
	if name, ok = content[".name"].(string); !ok {
		return errors.Errorf("malformed %s policy, missing '.name' field", key)
	}

	(*p).name = name

	return nil
}

func UnmarshalPolicy(data any, policy *Policy) error {
	switch typedData := data.(type) {
	case string:
		var err error
		if err = yaml.Unmarshal([]byte(typedData), policy); err != nil {
			return errors.New(err)
		}

		return nil
	default:
		var err error
		var policyYAML []byte
		if policyYAML, err = yaml.Marshal(typedData); err != nil {
			return errors.New(err)
		}
		return UnmarshalPolicy(string(policyYAML), policy)
	}
}

var TypesMap map[string]reflect.Type
