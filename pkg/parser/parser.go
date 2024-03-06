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

package parser

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/pb33f/libopenapi"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"os"
)

func Parse(specFile string) (*libopenapi.DocumentModel[v3high.Document], []error) {
	var specBytes []byte
	var err error
	if specBytes, err = os.ReadFile(specFile); err != nil {
		return nil, []error{errors.New(err)}
	}

	var specDoc libopenapi.Document
	if specDoc, err = libopenapi.NewDocument(specBytes); err != nil {
		return nil, []error{errors.New(err)}
	}

	var model *libopenapi.DocumentModel[v3high.Document]
	var errs []error
	if model, errs = specDoc.BuildV3Model(); len(errs) != 0 {
		var index int
		for index, err = range errs {
			errs[index] = errors.New(err)
		}
		return nil, errs
	}

	if model == nil {
		return nil, []error{errors.New(fmt.Errorf("could not parse spec file %s", specFile))}
	}

	return model, nil
}
