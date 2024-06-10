package mocks

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/keeper"
	"github.com/noble-assets/florin/x/florin/types"
)

func FlorinKeeper() (*keeper.Keeper, sdk.Context) {
	return FlorinWithKeepers(AccountKeeper{}, BankKeeper{})
}

func FlorinWithKeepers(account types.AccountKeeper, bank BankKeeper) (*keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey("transient_florin")

	reg := codectypes.NewInterfaceRegistry()
	types.RegisterInterfaces(reg)
	_ = codec.NewProtoCodec(reg)

	k := keeper.NewKeeper(key, "aeure", account, bank)

	bank = bank.WithSendCoinsRestriction(k.SendRestrictionFn)
	k.SetBankKeeper(bank)

	return k, testutil.DefaultContext(key, tkey)
}

//

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("noble", "noblepub")
	config.Seal()
}
