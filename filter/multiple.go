// filter/multiple.go
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
	. "github.com/streamwest-1629/guarantee"
)

type (
	multiFilter struct {
		filters []Filter
	}
)

// Make stacked filter instance. This checks string expression valid for all filters.
func Multiple(filters ...Filter) (stacked Filter) {
	return &multiFilter{
		filters: filters,
	}
}

func (multi *multiFilter) Initialize(chk string, initFunc GuaranteeFunc) error {
	return multi.nextInitialize(0, initFunc)(chk)
}

func (multi *multiFilter) Assign(previous, chk string, assignFunc GuaranteeAssignFunc) error {
	return multi.nextAssign(0, assignFunc)(previous, chk)
}

func (multi *multiFilter) Clone(chk string, cloneFunc GuaranteeFunc) error {
	return multi.nextClone(0, cloneFunc)(chk)
}

func (multi *multiFilter) nextInitialize(next int, finally GuaranteeFunc) GuaranteeFunc {
	if next < len(multi.filters) {
		return func(chk string) error {
			return multi.filters[next].Initialize(chk, multi.nextInitialize(next+1, finally))
		}
	} else {
		return finally
	}
}

func (multi *multiFilter) nextAssign(next int, finally GuaranteeAssignFunc) GuaranteeAssignFunc {
	if next < len(multi.filters) {
		return func(previous, chk string) error {
			return multi.filters[next].Assign(previous, chk, multi.nextAssign(next+1, finally))
		}
	} else {
		return finally
	}
}

func (multi *multiFilter) nextClone(next int, finally GuaranteeFunc) GuaranteeFunc {
	if next < len(multi.filters) {
		return func(chk string) error {
			return multi.filters[next].Clone(chk, multi.nextClone(next+1, finally))
		}
	} else {
		return finally
	}
}
