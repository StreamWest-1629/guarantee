// make.go
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

import "errors"

// Make new guarantee.Safety instance.
// When chk argument's expression is valid, this returns guaranteed value and nil error.
// If not, it returns some error
func MakeSafety(filter Filter, chk string) (guaranteed *Safety, err error) {

	if err := filter.Initialize(chk, func(chk string) error {

		// initialize instance
		guaranteed = &Safety{
			init_root(filter, chk),
		}
		return nil

	}); err != nil {
		// when failed initializing
		return nil, err
	} else {
		// when success
		return guaranteed, nil
	}
}

// Make new guarantee.Shared instance.
// When chk argument's expression is valid, this returns guaranteed value and nil error.
// If not, it returns some error.
func MakeShared(filter Filter, chk string) (guaranteed *Shared, err error) {

	if err := filter.Initialize(chk, func(chk string) error {

		// initialize instance
		guaranteed = &Shared{
			init_root(filter, chk),
		}
		return nil

	}); err != nil {
		// when failed initializing
		return nil, err
	} else {
		// when success
		return guaranteed, nil
	}
}

// Make new guarantee.Safety instance.
// When chk argument's expression is valid, this returns guaranteed value.
// If not, it returns un-initialized instance, it has empty string expression.
// To check guarantee.Safety instance is initialized,
// call guarantee.(*Safety).IsInitialized() function.
//
// String expression uses cloned value not to be changed by other variables sharing buffer.
func WrapSafety(filter Filter, chk string) (unsafe *Safety) {
	if str, err := MakeSafety(filter, chk); err != nil {
		return &Safety{
			uninit_root(filter),
		}
	} else {
		return str
	}
}

// Make new guarantee.Shared instance.
// When chk argument's expression is valid, this returns guaranteed value.
// If not, it returns un-initialized instance, it has empty string expression.
// To check guarantee.Shared instance is initialized,
// call guarantee.(*Shared).IsInitialized() function.
func WrapShared(filter Filter, chk string) (unsafe *Shared) {
	if str, err := MakeShared(filter, chk); err != nil {
		return &Shared{
			uninit_root(filter),
		}
	} else {
		return str
	}
}

// Check filter passing or not.
// Call Initialize function in filter,
// initFunc argument is given function returning any error.
func Check(filter Filter, chk string) (passed bool) {
	passed = false
	err := errors.New("this is check function.")
	test := func(str string) error {
		passed = true
		return err
	}
	filter.Initialize(chk, test)
	return
}
