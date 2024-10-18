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
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/types/blacklist"
)

type Keeper struct {
	authority string

	schema       collections.Schema
	storeService store.KVStoreService
	eventService event.Service

	AllowedDenoms    collections.KeySet[string]
	Owner            collections.Map[string, string]
	PendingOwner     collections.Map[string, string]
	Systems          collections.KeySet[[]byte]
	Admins           collections.KeySet[[]byte]
	MintAllowance    collections.Map[[]byte, []byte]
	MaxMintAllowance collections.Map[string, []byte]

	BlacklistOwner        collections.Item[string]
	BlacklistPendingOwner collections.Item[string]
	BlacklistAdmins       collections.KeySet[string]
	Adversaries           collections.KeySet[string]

	cdc          codec.Codec
	addressCodec address.Codec
	bankKeeper   types.BankKeeper
}

func NewKeeper(
	authority string,
	storeService store.KVStoreService,
	eventService event.Service,
	cdc codec.Codec,
	addressCodec address.Codec,
	bankKeeper types.BankKeeper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)

	keeper := &Keeper{
		authority: authority,

		storeService: storeService,
		eventService: eventService,

		AllowedDenoms:    collections.NewKeySet(builder, types.AllowedDenomPrefix, "allowedDenoms", collections.StringKey),
		Owner:            collections.NewMap(builder, types.OwnerPrefix, "owner", collections.StringKey, collections.StringValue),
		PendingOwner:     collections.NewMap(builder, types.PendingOwnerPrefix, "pendingOwner", collections.StringKey, collections.StringValue),
		Systems:          collections.NewKeySet(builder, types.SystemPrefix, "systems", collections.BytesKey),
		Admins:           collections.NewKeySet(builder, types.AdminPrefix, "admins", collections.BytesKey),
		MintAllowance:    collections.NewMap(builder, types.MintAllowancePrefix, "mintAllowance", collections.BytesKey, collections.BytesValue),
		MaxMintAllowance: collections.NewMap(builder, types.MaxMintAllowancePrefix, "maxMintAllowance", collections.StringKey, collections.BytesValue),

		BlacklistOwner:        collections.NewItem(builder, blacklist.OwnerKey, "blacklistOwner", collections.StringValue),
		BlacklistPendingOwner: collections.NewItem(builder, blacklist.PendingOwnerKey, "blacklistPendingOwner", collections.StringValue),
		BlacklistAdmins:       collections.NewKeySet(builder, blacklist.AdminPrefix, "blacklistAdmins", collections.StringKey),
		Adversaries:           collections.NewKeySet(builder, blacklist.AdversaryPrefix, "adversaries", collections.StringKey),

		cdc:          cdc,
		addressCodec: addressCodec,
		bankKeeper:   bankKeeper,
	}

	schema, err := builder.Build()
	if err != nil {
		panic(err)
	}

	keeper.schema = schema
	return keeper
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

// SendRestrictionFn executes necessary checks against all EURe, GBPe, ISKe, USDe transfers.
func (k *Keeper) SendRestrictionFn(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	for _, allowedDenom := range k.GetAllowedDenoms(ctx) {
		if amount := amt.AmountOf(allowedDenom); !amount.IsZero() {
			isAdversary := k.IsAdversary(ctx, fromAddr.String())
			_ = k.eventService.EventManager(ctx).Emit(ctx, &blacklist.Decision{
				From:   fromAddr.String(),
				To:     toAddr.String(),
				Amount: amount,
				Valid:  !isAdversary,
			})

			if isAdversary {
				return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, allowedDenom)
			}
		}
	}

	return toAddr, nil
}
