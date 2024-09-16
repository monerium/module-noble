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

package mocks

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/monerium/module-noble/v2/x/florin/types"
)

var _ types.BankKeeper = BankKeeper{}

type BankKeeper struct {
	Balances    map[string]sdk.Coins
	Restriction SendRestrictionFn
}

func (k BankKeeper) BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	address := authtypes.NewModuleAddress(moduleName).String()

	balance := k.Balances[address]
	newBalance, negative := balance.SafeSub(amt)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "%s is smaller than %s", balance, amt)
	}

	k.Balances[address] = newBalance

	return nil
}

func (k BankKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	return sdk.NewCoin(denom, k.Balances[addr.String()].AmountOf(denom))
}

func (k BankKeeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	if denom == "uusdc" {
		return sdk.NewCoin(denom, sdk.NewIntFromUint64(1_000_000))
	}

	return sdk.NewCoin(denom, sdk.ZeroInt())
}

func (k BankKeeper) MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error {
	address := authtypes.NewModuleAddress(moduleName).String()
	k.Balances[address] = k.Balances[address].Add(amt...)

	return nil
}

func (k BankKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	recipientAddr := authtypes.NewModuleAddress(recipientModule)

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

func (k BankKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	senderAddr := authtypes.NewModuleAddress(senderModule)

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

//

type SendRestrictionFn func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error)

func NoOpSendRestrictionFn(_ sdk.Context, _, toAddr sdk.AccAddress, _ sdk.Coins) (sdk.AccAddress, error) {
	return toAddr, nil
}

func (k BankKeeper) WithSendCoinsRestriction(check SendRestrictionFn) BankKeeper {
	oldRestriction := k.Restriction
	k.Restriction = func(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
		newToAddr, err = check(ctx, fromAddr, toAddr, amt)
		if err != nil {
			return newToAddr, err
		}
		return oldRestriction(ctx, fromAddr, toAddr, amt)
	}
	return k
}

func (k BankKeeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	toAddr, err := k.Restriction(ctx, fromAddr, toAddr, amt)
	if err != nil {
		return err
	}

	balance := k.Balances[fromAddr.String()]
	newBalance, negative := balance.SafeSub(amt)
	if negative {
		return sdkerrors.Wrapf(errors.ErrInsufficientFunds, "%s is smaller than %s", balance, amt)
	}

	k.Balances[fromAddr.String()] = newBalance
	k.Balances[toAddr.String()] = k.Balances[toAddr.String()].Add(amt...)

	return nil
}
