// guarantee.go
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

package guarantee

type (

	// The function style to check string expression and assign it.
	// When chk argument is valid, this returns nil.
	// If not, it returns some error and doesn't assign it.
	GuaranteeFunc func(chk string) error

	GuaranteeAssignFunc func(previous, chk string) error

	// Interface to check string expression is valid or not.
	// When chk argument is valid, this returns nil, If not, it returns some error.
	Filter interface {
		Initialize(chk string, initFunc GuaranteeFunc) error
		Assign(previous, chk string, assignFunc GuaranteeAssignFunc) error
		Clone(original string, cloneFunc GuaranteeFunc) error
	}

	root struct {
		filter     Filter
		guaranteed string
		inited     bool
	}

	// Contains string value to be checked that string expression is valid.
	// Uses independented string buffer.
	Safety struct{ *root }

	// Contains string value to be checked that string expression is valid.
	// Uses shared buffer.
	Shared struct{ *root }
)

func init_root(filter Filter, guaranteed string) (initialized *root) {
	return &root{filter, guaranteed, true}
}

func uninit_root(filter Filter) (uninitialized *root) {
	return &root{filter, "", false}
}

func (r *root) assign(value string, deep bool) error {
	if r.guaranteed == value {
		return nil
	}
	return r.filter.Assign(r.guaranteed, value, func(previous, chk string) error {
		// when filtered successfully, assign value
		if deep {
			r.guaranteed = CloneBuiltinString(value)
		} else {
			r.guaranteed = value
		}
		r.inited = true
		return nil
	})
}

func (r *root) clone(deep bool) (dest *root, err error) {
	if err := r.filter.Clone(r.guaranteed, func(chk string) error {
		// when filtered successfully, clone value
		dest.filter, dest.inited = r.filter, r.inited
		if deep {
			dest.guaranteed = CloneBuiltinString(r.guaranteed)
		} else {
			dest.guaranteed = r.guaranteed
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return dest, nil
}
