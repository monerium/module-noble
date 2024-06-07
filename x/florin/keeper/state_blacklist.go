package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

//

func (k *Keeper) GetBlacklistOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(blacklist.OwnerKey))
}

func (k *Keeper) SetBlacklistOwner(ctx sdk.Context, owner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blacklist.OwnerKey, []byte(owner))
}

//

func (k *Keeper) DeleteBlacklistPendingOwner(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(blacklist.PendingOwnerKey)
}

func (k *Keeper) GetBlacklistPendingOwner(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(blacklist.PendingOwnerKey))
}

func (k *Keeper) SetBlacklistPendingOwner(ctx sdk.Context, pendingOwner string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blacklist.PendingOwnerKey, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteBlacklistAdmin(ctx sdk.Context, admin string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(blacklist.AdminKey(admin))
}

func (k *Keeper) IsBlacklistAdmin(ctx sdk.Context, admin string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(blacklist.AdminKey(admin))
}

func (k *Keeper) SetBlacklistAdmin(ctx sdk.Context, admin string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blacklist.AdminKey(admin), []byte{})
}

//

func (k *Keeper) DeleteAdversary(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(blacklist.AdversaryKey(address))
}

func (k *Keeper) IsAdversary(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(blacklist.AdversaryKey(address))
}

func (k *Keeper) SetAdversary(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blacklist.AdversaryKey(address), []byte{})
}
