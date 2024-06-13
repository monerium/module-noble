package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

func (k *Keeper) GetBlacklistAdmins(ctx sdk.Context) (admins []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), blacklist.AdminPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		admins = append(admins, string(itr.Key()))
	}

	return
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

func (k *Keeper) GetAdversaries(ctx sdk.Context) (adversaries []string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), blacklist.AdversaryPrefix)
	itr := store.Iterator(nil, nil)

	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		adversaries = append(adversaries, string(itr.Key()))
	}

	return
}

func (k *Keeper) IsAdversary(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(blacklist.AdversaryKey(address))
}

func (k *Keeper) SetAdversary(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(blacklist.AdversaryKey(address), []byte{})
}
