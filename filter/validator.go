// filter/validator.go
// Copyright (C) 2021 Kasai Koji

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"errors"
	"regexp"

	. "github.com/streamwest-1629/guarantee"
)

type (
	Validator func(chk string) error
)

func (v Validator) Initialize(chk string, initFunc GuaranteeFunc) error {
	if err := v(chk); err != nil {
		return err
	} else {
		return initFunc(chk)
	}
}

func (v Validator) Assign(previous, chk string, assignFunc GuaranteeAssignFunc) error {
	if err := v(chk); err != nil {
		return err
	} else {
		return assignFunc(previous, chk)
	}
}

func (v Validator) Clone(chk string, cloneFunc GuaranteeFunc) error {
	return cloneFunc(chk)
}

// Make filter using regular expression.
func RegexpFilter(r *regexp.Regexp) Filter {
	return Validator(
		func(chk string) error {
			if locates := r.FindStringIndex(chk); len(locates) != 2 {
				// when cannot found matches.
				return errors.New("cannot found what regular expression matches.")
			} else if locates[0] != 0 || locates[1] != len(chk) {
				// when partially found matches.
				return errors.New("found what of regular expression matches, but it isn't fully in string expression.")
			}
			return nil
		},
	)
}

// Make filter using list. Internal logic uses converted builtin map.
func ListMatchingFilter(list ...string) Filter {
	mapping := make(map[string]interface{})
	for _, str := range list {
		mapping[str] = nil
	}
	return Validator(
		func(chk string) error {
			if _, ok := mapping[chk]; !ok {
				return errors.New("cannot found expression in list")
			} else {
				return nil
			}
		},
	)
}
