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

package types

import (
	"fmt"
	"slices"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	"github.com/monerium/module-noble/v2/types/blacklist"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		BlacklistState: blacklist.DefaultGenesisState(),
		AllowedDenoms:  []string{"ueure"},
		MaxMintAllowances: map[string]string{
			"ueure": "3000000000000", // 3,000,000 EURe
		},
	}
}

func (gs *GenesisState) Validate(cdc address.Codec) error {
	if err := gs.BlacklistState.Validate(cdc); err != nil {
		return err
	}

	for denom, owner := range gs.Owners {
		if !slices.Contains(gs.AllowedDenoms, denom) {
			return fmt.Errorf("found an owner (%s) for a not allowed denom %s", owner, denom)
		}

		if _, err := cdc.StringToBytes(owner); err != nil {
			return fmt.Errorf("invalid owner address (%s) for denom %s: %s", owner, denom, err)
		}
	}

	for denom, pendingOwner := range gs.PendingOwners {
		if !slices.Contains(gs.AllowedDenoms, denom) {
			return fmt.Errorf("found a pending owner (%s) for a not allowed denom %s", pendingOwner, denom)
		}

		if _, err := cdc.StringToBytes(pendingOwner); err != nil {
			return fmt.Errorf("invalid pending owner address (%s) for denom %s: %s", pendingOwner, denom, err)
		}
	}

	for _, system := range gs.Systems {
		if !slices.Contains(gs.AllowedDenoms, system.Denom) {
			return fmt.Errorf("found a system account (%s) for a not allowed denom %s", system.Address, system.Denom)
		}

		if _, err := cdc.StringToBytes(system.Address); err != nil {
			return fmt.Errorf("invalid system address (%s) for denom %s: %s", system.Address, system.Denom, err)
		}
	}

	for _, admin := range gs.Admins {
		if !slices.Contains(gs.AllowedDenoms, admin.Denom) {
			return fmt.Errorf("found an admin account (%s) for a not allowed denom %s", admin.Address, admin.Denom)
		}

		if _, err := cdc.StringToBytes(admin.Address); err != nil {
			return fmt.Errorf("invalid admin address (%s) for denom %s: %s", admin.Address, admin.Denom, err)
		}
	}

	for _, entry := range gs.MintAllowances {
		if !slices.Contains(gs.AllowedDenoms, entry.Denom) {
			return fmt.Errorf("found a minter allowance (%s) for a not allowed denom %s", entry.Address, entry.Denom)
		}

		if _, err := cdc.StringToBytes(entry.Address); err != nil {
			return fmt.Errorf("invalid minter address (%s) for denom %s: %s", entry.Address, entry.Denom, err)
		}
	}

	for denom, maxAllowance := range gs.MaxMintAllowances {
		if !slices.Contains(gs.AllowedDenoms, denom) {
			return fmt.Errorf("found a max mint allowance (%s) for a not allowed denom %s", maxAllowance, denom)
		}

		if _, ok := math.NewIntFromString(maxAllowance); !ok {
			return fmt.Errorf("invalid max mint allowance (%s) for denom %s", maxAllowance, denom)
		}
	}

	return nil
}
