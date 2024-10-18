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

package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/monerium/module-noble/v2/types"
)

//

func (k *Keeper) GetAllowedDenoms(ctx context.Context) (allowedDenoms []string) {
	_ = k.AllowedDenoms.Walk(ctx, nil, func(denom string) (bool, error) {
		allowedDenoms = append(allowedDenoms, denom)
		return false, nil
	})
	return
}

func (k *Keeper) IsAllowedDenom(ctx context.Context, denom string) bool {
	allowed, _ := k.AllowedDenoms.Has(ctx, denom)
	return allowed
}

func (k *Keeper) SetAllowedDenom(ctx context.Context, denom string) error {
	return k.AllowedDenoms.Set(ctx, denom)
}

//

func (k *Keeper) GetOwner(ctx context.Context, denom string) string {
	owner, _ := k.Owner.Get(ctx, denom)
	return owner
}

func (k *Keeper) GetOwners(ctx context.Context) map[string]string {
	owners := make(map[string]string)
	_ = k.Owner.Walk(ctx, nil, func(key string, value string) (bool, error) {
		owners[key] = value
		return false, nil
	})
	return owners
}

func (k *Keeper) SetOwner(ctx context.Context, denom string, owner string) error {
	return k.Owner.Set(ctx, denom, owner)
}

//

func (k *Keeper) DeletePendingOwner(ctx context.Context, denom string) error {
	return k.PendingOwner.Remove(ctx, denom)
}

func (k *Keeper) GetPendingOwner(ctx context.Context, denom string) string {
	pendingOwner, _ := k.PendingOwner.Get(ctx, denom)
	return pendingOwner
}

func (k *Keeper) GetPendingOwners(ctx context.Context) map[string]string {
	pendingOwners := make(map[string]string)
	_ = k.PendingOwner.Walk(ctx, nil, func(key string, value string) (bool, error) {
		pendingOwners[key] = value
		return false, nil
	})
	return pendingOwners
}

func (k *Keeper) SetPendingOwner(ctx context.Context, denom string, pendingOwner string) error {
	return k.PendingOwner.Set(ctx, denom, pendingOwner)
}

//

func (k *Keeper) DeleteSystem(ctx context.Context, denom string, address string) error {
	return k.Systems.Remove(ctx, types.SystemKey(denom, address))
}

func (k *Keeper) GetSystemsByDenom(ctx context.Context, denom string) (systems []string) {
	prefix := []byte(denom)
	itr, _ := k.Systems.Iterate(ctx, new(collections.Range[[]byte]).Prefix(prefix))

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		systems = append(systems, string(key[len(prefix):]))
	}

	return
}

func (k *Keeper) GetSystems(ctx context.Context) (systems []types.Account) {
	for _, allowedDenom := range k.GetAllowedDenoms(ctx) {
		for _, system := range k.GetSystemsByDenom(ctx, allowedDenom) {
			systems = append(systems, types.Account{
				Denom:   allowedDenom,
				Address: system,
			})
		}
	}

	return
}

func (k *Keeper) IsSystem(ctx context.Context, denom string, address string) bool {
	system, _ := k.Systems.Has(ctx, types.SystemKey(denom, address))
	return system
}

func (k *Keeper) SetSystem(ctx context.Context, denom string, address string) error {
	return k.Systems.Set(ctx, types.SystemKey(denom, address))
}

//

func (k *Keeper) DeleteAdmin(ctx context.Context, denom string, admin string) error {
	return k.Admins.Remove(ctx, types.AdminKey(denom, admin))
}

func (k *Keeper) GetAdminsByDenom(ctx context.Context, denom string) (admins []string) {
	prefix := []byte(denom)
	itr, _ := k.Admins.Iterate(ctx, new(collections.Range[[]byte]).Prefix(prefix))

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		admins = append(admins, string(key[len(prefix):]))
	}

	return
}

func (k *Keeper) GetAdmins(ctx context.Context) (admins []types.Account) {
	for _, allowedDenom := range k.GetAllowedDenoms(ctx) {
		for _, admin := range k.GetAdminsByDenom(ctx, allowedDenom) {
			admins = append(admins, types.Account{
				Denom:   allowedDenom,
				Address: admin,
			})
		}
	}

	return
}

func (k *Keeper) IsAdmin(ctx context.Context, denom string, admin string) bool {
	isAdmin, _ := k.Admins.Has(ctx, types.AdminKey(denom, admin))
	return isAdmin
}

func (k *Keeper) SetAdmin(ctx context.Context, denom string, admin string) error {
	return k.Admins.Set(ctx, types.AdminKey(denom, admin))
}

//

func (k *Keeper) GetMintAllowance(ctx context.Context, denom string, address string) (allowance math.Int) {
	allowance = math.ZeroInt()
	bz, err := k.MintAllowance.Get(ctx, types.MintAllowanceKey(denom, address))
	if err != nil {
		return
	}
	_ = allowance.Unmarshal(bz)

	return
}

func (k *Keeper) GetMintAllowancesByDenom(ctx context.Context, denom string) (allowances []types.Allowance) {
	prefix := []byte(denom)
	itr, _ := k.MintAllowance.Iterate(ctx, new(collections.Range[[]byte]).Prefix(prefix))

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		key, _ := itr.Key()
		value, _ := itr.Value()

		var allowance math.Int
		err := allowance.Unmarshal(value)
		if err != nil {
			continue
		}

		allowances = append(allowances, types.Allowance{
			Denom:     denom,
			Address:   string(key[len(prefix):]),
			Allowance: allowance,
		})
	}

	return
}

func (k *Keeper) GetMintAllowances(ctx context.Context) (allowances []types.Allowance) {
	for _, allowedDenom := range k.GetAllowedDenoms(ctx) {
		allowances = append(allowances, k.GetMintAllowancesByDenom(ctx, allowedDenom)...)
	}

	return
}

func (k *Keeper) SetMintAllowance(ctx context.Context, denom string, address string, allowance math.Int) error {
	bz, _ := allowance.Marshal()
	return k.MintAllowance.Set(ctx, types.MintAllowanceKey(denom, address), bz)
}

//

func (k *Keeper) GetMaxMintAllowance(ctx context.Context, denom string) (maxAllowance math.Int) {
	maxAllowance = math.ZeroInt()
	bz, err := k.MaxMintAllowance.Get(ctx, denom)
	if err != nil {
		return
	}
	_ = maxAllowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetMaxMintAllowances(ctx context.Context) (maxAllowances map[string]string) {
	maxAllowances = make(map[string]string)
	_ = k.MaxMintAllowance.Walk(ctx, nil, func(key string, value []byte) (stop bool, err error) {
		maxAllowances[key] = string(value)
		return false, nil
	})
	return
}

func (k *Keeper) SetMaxMintAllowance(ctx context.Context, denom string, maxAllowance math.Int) error {
	bz, _ := maxAllowance.Marshal()
	return k.MaxMintAllowance.Set(ctx, denom, bz)
}
