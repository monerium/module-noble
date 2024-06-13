package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/types"
)

//

func (k *Keeper) GetOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.OwnerKey))
}

func (k *Keeper) SetOwner(ctx sdk.Context, owner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.OwnerKey, []byte(owner))
}

//

func (k *Keeper) DeletePendingOwner(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PendingOwnerKey)
}

func (k *Keeper) GetPendingOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PendingOwnerKey))
}

func (k *Keeper) SetPendingOwner(ctx sdk.Context, pendingOwner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PendingOwnerKey, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteSystem(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.SystemKey(address))
}

func (k *Keeper) GetSystems(ctx sdk.Context) (systems []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SystemPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		systems = append(systems, string(itr.Key()))
	}

	return
}

func (k *Keeper) IsSystem(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.SystemKey(address))
}

func (k *Keeper) SetSystem(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SystemKey(address), []byte{})
}

//

func (k *Keeper) DeleteAdmin(ctx sdk.Context, admin string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AdminKey(admin))
}

func (k *Keeper) GetAdmins(ctx sdk.Context) (admins []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AdminPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		admins = append(admins, string(itr.Key()))
	}

	return
}

func (k *Keeper) IsAdmin(ctx sdk.Context, admin string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.AdminKey(admin))
}

func (k *Keeper) SetAdmin(ctx sdk.Context, admin string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AdminKey(admin), []byte{})
}

//

func (k *Keeper) GetMintAllowance(ctx sdk.Context, address string) (allowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MintAllowanceKey(address))
	if bz == nil {
		return sdk.ZeroInt()
	}

	_ = allowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetMintAllowances(ctx sdk.Context) map[string]string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MintAllowancePrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	allowances := make(map[string]string)

	for ; itr.Valid(); itr.Next() {
		var allowance sdk.Int
		_ = allowance.Unmarshal(itr.Value())

		allowances[string(itr.Key())] = allowance.String()
	}

	return allowances
}

func (k *Keeper) SetMintAllowance(ctx sdk.Context, address string, allowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := allowance.Marshal()
	store.Set(types.MintAllowanceKey(address), bz)
}

//

func (k *Keeper) GetMaxMintAllowance(ctx sdk.Context) (maxAllowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MaxMintAllowanceKey)
	if bz == nil {
		return sdk.ZeroInt()
	}

	_ = maxAllowance.Unmarshal(bz)
	return
}

func (k *Keeper) SetMaxMintAllowance(ctx sdk.Context, maxAllowance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := maxAllowance.Marshal()
	store.Set(types.MaxMintAllowanceKey, bz)
}
