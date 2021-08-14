// string.go
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

// Assign string expression to this instance.
//
// Check string expression is valid or not.
// When it is valid, this assigns new expression and returns nil error.
// If not, this doesn't assigns new expression and returns some error.
//
// String expression uses cloned value not to be changed by other variables sharing buffer.
func (dest *Safety) Assign(chk string) error {
	return dest.assign(chk, true)
}

// Assign string expression to this instance.
//
// Check string expression is valid or not.
// When it is valid, this assigns new expression and returns nil error.
// If not, this doesn't assigns new expression and returns some error.
func (dest *Shared) Assign(chk string) error {
	return dest.assign(chk, false)
}

// Clone instance.
// This holding string expression's buffer is shared with cloned instance.
func (src *Safety) Clone() (cloned *Safety, err error) {
	if cloned, err := src.clone(false); err != nil {
		return nil, err
	} else {
		return &Safety{cloned}, nil
	}
}

// Deep copy instance.
func (src *Shared) Clone() (cloned *Shared, err error) {
	if cloned, err := src.clone(true); err != nil {
		return nil, err
	} else {
		return &Shared{cloned}, nil
	}
}

// Copy instance and change filter.
//
// Check string expression is valid or not.
// When it is valid, this asigns new expression and returns nil error.
// If not, it returns some error.
//
// This holding string expression's buffer is shared with cloned instance.
func (src *Safety) Changed(filter Filter) (changed *Safety, err error) {
	if err := filter.Initialize(src.guaranteed, func(chk string) error {
		// when filtered successfully, assign value
		changed = &Safety{init_root(filter, chk)}
		return nil
	}); err != nil {
		return nil, err
	} else {
		return changed, nil
	}
}

// Copy instance and change filter.
//
// Check string expression is valid or not.
// When it is valid, this asigns new expression and returns nil error.
// If not, it returns some error.
func (src *Shared) Changed(filter Filter) (changed *Shared, err error) {
	return MakeShared(filter, src.guaranteed)
}

// Convert type to builtin string expression.
func (origin *Safety) String() string {
	return CloneBuiltinString(origin.guaranteed)
}

// Convert type to built in string expression.
func (origin *Shared) String() string {
	return CloneBuiltinString(origin.guaranteed)
}
