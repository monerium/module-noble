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
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/types"
)

//

func (k *Keeper) GetAllowedDenoms(ctx sdk.Context) (allowedDenoms []string) {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.AllowedDenomPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		allowedDenoms = append(allowedDenoms, string(itr.Key()))
	}

	return
}

func (k *Keeper) IsAllowedDenom(ctx sdk.Context, denom string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.AllowedDenomKey(denom))
}

func (k *Keeper) SetAllowedDenom(ctx sdk.Context, denom string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.AllowedDenomKey(denom), []byte{})
}

//

func (k *Keeper) GetOwner(ctx sdk.Context, denom string) string {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return string(store.Get(types.OwnerKey(denom)))
}

func (k *Keeper) GetOwners(ctx sdk.Context) map[string]string {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.OwnerPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	owners := make(map[string]string)
	for ; itr.Valid(); itr.Next() {
		owners[string(itr.Key())] = string(itr.Value())
	}

	return owners
}

func (k *Keeper) SetOwner(ctx sdk.Context, denom string, owner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.OwnerKey(denom), []byte(owner))
}

//

func (k *Keeper) DeletePendingOwner(ctx sdk.Context, denom string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.PendingOwnerKey(denom))
}

func (k *Keeper) GetPendingOwner(ctx sdk.Context, denom string) string {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return string(store.Get(types.PendingOwnerKey(denom)))
}

func (k *Keeper) GetPendingOwners(ctx sdk.Context) map[string]string {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.PendingOwnerPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	pendingOwners := make(map[string]string)
	for ; itr.Valid(); itr.Next() {
		pendingOwners[string(itr.Key())] = string(itr.Value())
	}

	return pendingOwners
}

func (k *Keeper) SetPendingOwner(ctx sdk.Context, denom string, pendingOwner string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.PendingOwnerKey(denom), []byte(pendingOwner))
}

//

func (k *Keeper) DeleteSystem(ctx sdk.Context, denom string, address string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.SystemKey(denom, address))
}

func (k *Keeper) GetSystemsByDenom(ctx sdk.Context, denom string) (systems []string) {
	bz := append(types.SystemPrefix, []byte(denom)...)
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, bz)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		systems = append(systems, string(itr.Key()))
	}

	return
}

func (k *Keeper) GetSystems(ctx sdk.Context) (systems []types.Account) {
	allowedDenoms := k.GetAllowedDenoms(ctx)

	for _, allowedDenom := range allowedDenoms {
		for _, system := range k.GetSystemsByDenom(ctx, allowedDenom) {
			systems = append(systems, types.Account{
				Denom:   allowedDenom,
				Address: system,
			})
		}
	}

	return
}

func (k *Keeper) IsSystem(ctx sdk.Context, denom string, address string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.SystemKey(denom, address))
}

func (k *Keeper) SetSystem(ctx sdk.Context, denom string, address string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.SystemKey(denom, address), []byte{})
}

//

func (k *Keeper) DeleteAdmin(ctx sdk.Context, denom string, admin string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.AdminKey(denom, admin))
}

func (k *Keeper) GetAdminsByDenom(ctx sdk.Context, denom string) (admins []string) {
	bz := append(types.AdminPrefix, []byte(denom)...)
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, bz)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		admins = append(admins, string(itr.Key()))
	}

	return
}

func (k *Keeper) GetAdmins(ctx sdk.Context) (admins []types.Account) {
	allowedDenoms := k.GetAllowedDenoms(ctx)

	for _, allowedDenom := range allowedDenoms {
		for _, admin := range k.GetAdminsByDenom(ctx, allowedDenom) {
			admins = append(admins, types.Account{
				Denom:   allowedDenom,
				Address: admin,
			})
		}
	}

	return
}

func (k *Keeper) IsAdmin(ctx sdk.Context, denom string, admin string) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.AdminKey(denom, admin))
}

func (k *Keeper) SetAdmin(ctx sdk.Context, denom string, admin string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.AdminKey(denom, admin), []byte{})
}

//

func (k *Keeper) GetMintAllowance(ctx sdk.Context, denom string, address string) (allowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.MintAllowanceKey(denom, address))

	allowance = math.ZeroInt()
	_ = allowance.Unmarshal(bz)

	return
}

func (k *Keeper) GetMintAllowancesByDenom(ctx sdk.Context, denom string) (allowances []types.Allowance) {
	bz := append(types.MintAllowancePrefix, []byte(denom)...)
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, bz)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		var allowance math.Int
		_ = allowance.Unmarshal(itr.Value())

		allowances = append(allowances, types.Allowance{
			Denom:     denom,
			Address:   string(itr.Key()),
			Allowance: allowance,
		})
	}

	return
}

func (k *Keeper) GetMintAllowances(ctx sdk.Context) (allowances []types.Allowance) {
	allowedDenoms := k.GetAllowedDenoms(ctx)

	for _, allowedDenom := range allowedDenoms {
		allowances = append(allowances, k.GetMintAllowancesByDenom(ctx, allowedDenom)...)
	}

	return
}

func (k *Keeper) SetMintAllowance(ctx sdk.Context, denom string, address string, allowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, _ := allowance.Marshal()
	store.Set(types.MintAllowanceKey(denom, address), bz)
}

//

func (k *Keeper) GetMaxMintAllowance(ctx sdk.Context, denom string) (maxAllowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.MaxMintAllowanceKey(denom))

	maxAllowance = math.ZeroInt()
	_ = maxAllowance.Unmarshal(bz)

	return
}

func (k *Keeper) GetMaxMintAllowances(ctx sdk.Context) map[string]string {
	adapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(adapter, types.MaxMintAllowancePrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	maxAllowances := make(map[string]string)
	for ; itr.Valid(); itr.Next() {
		maxAllowance := math.ZeroInt()
		_ = maxAllowance.Unmarshal(itr.Value())

		maxAllowances[string(itr.Key())] = maxAllowance.String()
	}

	return maxAllowances
}

func (k *Keeper) SetMaxMintAllowance(ctx sdk.Context, denom string, maxAllowance math.Int) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, _ := maxAllowance.Marshal()
	store.Set(types.MaxMintAllowanceKey(denom), bz)
}
