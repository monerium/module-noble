// Copyright 2024 Monerium ehf.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blacklist

import (
	"fmt"

	"cosmossdk.io/core/address"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs *GenesisState) Validate(cdc address.Codec) error {
	if gs.Owner != "" {
		if _, err := cdc.StringToBytes(gs.Owner); err != nil {
			return fmt.Errorf("invalid blacklist owner address (%s): %s", gs.Owner, err)
		}
	}

	if gs.PendingOwner != "" {
		if _, err := cdc.StringToBytes(gs.PendingOwner); err != nil {
			return fmt.Errorf("invalid pending blacklist owner address (%s): %s", gs.PendingOwner, err)
		}
	}

	for _, admin := range gs.Admins {
		if _, err := cdc.StringToBytes(admin); err != nil {
			return fmt.Errorf("invalid admin address (%s): %s", admin, err)
		}
	}

	for _, adversary := range gs.Adversaries {
		if _, err := cdc.StringToBytes(adversary); err != nil {
			return fmt.Errorf("invalid adversary address (%s): %s", adversary, err)
		}
	}

	return nil
}
