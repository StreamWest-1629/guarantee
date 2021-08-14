// filter/identifier.go
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

	. "github.com/streamwest-1629/guarantee"
)

type (
	limitedDuplicate struct {
		mapping map[string]int
		limit   int
	}
)

// Make filter limitted cloning. limit argument doesn't including copy times.
func LimitedCloningFilter(limit int) Filter {
	if limit < 0 {
		limit = 0
	}
	return &limitedDuplicate{
		mapping: make(map[string]int),
		limit:   0,
	}
}

func (limited *limitedDuplicate) Initialize(chk string, initFunc GuaranteeFunc) error {
	if val, exist := limited.mapping[chk]; !exist {
		if err := initFunc(chk); err != nil {
			return err
		} else {
			limited.mapping[chk] = 0
			return nil
		}
	} else if val < limited.limit {
		if err := initFunc(chk); err != nil {
			return err
		} else {
			limited.mapping[chk] = val + 1
			return nil
		}
	} else {
		return errors.New("reached the duplicating limit")
	}
}

func (limited *limitedDuplicate) Assign(previous, chk string, assignFunc GuaranteeAssignFunc) error {
	if val, exist := limited.mapping[chk]; !exist {
		if err := assignFunc(previous, chk); err != nil {
			return err
		} else {
			limited.mapping[chk] = 0
		}
	} else if val < limited.limit {
		if err := assignFunc(previous, chk); err != nil {
			return err
		} else {
			limited.mapping[chk] = val + 1
		}
	} else {
		return errors.New("reached the duplicating limit")
	}

	if val, _ := limited.mapping[previous]; val > 0 {
		limited.mapping[previous] = val - 1
	} else {
		delete(limited.mapping, previous)
	}
	return nil
}

func (limited *limitedDuplicate) Clone(chk string, cloneFunc GuaranteeFunc) error {
	if val, _ := limited.mapping[chk]; val < limited.limit {
		if err := cloneFunc(chk); err != nil {
			return err
		} else {
			limited.mapping[chk] = val + 1
			return nil
		}
	} else {
		return errors.New("reached the duplicating limit")
	}
}
