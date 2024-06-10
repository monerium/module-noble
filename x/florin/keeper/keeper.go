package keeper

import (
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	Denom    string

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	denom string,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	return &Keeper{
		storeKey: storeKey,
		Denom:    denom,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

// SendRestrictionFn executes necessary checks against all EURe transfers.
func (k *Keeper) SendRestrictionFn(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	if amount := amt.AmountOf(k.Denom); !amount.IsZero() {
		valid := !k.IsAdversary(ctx, fromAddr.String())
		_ = ctx.EventManager().EmitTypedEvent(&blacklist.Decision{
			From:   fromAddr.String(),
			To:     toAddr.String(),
			Amount: amount,
			Valid:  valid,
		})

		if !valid {
			return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, k.Denom)
		}
	}

	return toAddr, nil
}
