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

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/types/blacklist"
)

type Keeper struct {
	authority    string
	storeService store.KVStoreService

	cdc          codec.Codec
	addressCodec address.Codec
	bankKeeper   types.BankKeeper
}

func NewKeeper(
	authority string,
	storeService store.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	bankKeeper types.BankKeeper,
) *Keeper {
	return &Keeper{
		authority:    authority,
		storeService: storeService,

		cdc:          cdc,
		addressCodec: addressCodec,
		bankKeeper:   bankKeeper,
	}
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

// SendRestrictionFn executes necessary checks against all EURe, GBPe, ISKe, USDe transfers.
func (k *Keeper) SendRestrictionFn(goCtx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, allowedDenom := range k.GetAllowedDenoms(ctx) {
		if amount := amt.AmountOf(allowedDenom); !amount.IsZero() {
			valid := !k.IsAdversary(ctx, fromAddr.String())
			_ = ctx.EventManager().EmitTypedEvent(&blacklist.Decision{
				From:   fromAddr.String(),
				To:     toAddr.String(),
				Amount: amount,
				Valid:  valid,
			})

			if !valid {
				return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, allowedDenom)
			}
		}
	}

	return toAddr, nil
}