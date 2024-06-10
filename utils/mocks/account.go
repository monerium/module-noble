package mocks

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/noble-assets/florin/x/florin/types"
)

var _ types.AccountKeeper = AccountKeeper{}

type AccountKeeper struct {
	Accounts map[string]authtypes.AccountI
}

func (k AccountKeeper) GetAccount(_ sdk.Context, addr sdk.AccAddress) authtypes.AccountI {
	// NOTE: The bech32 prefix is already set when mocking Florin.
	return k.Accounts[addr.String()]
}
