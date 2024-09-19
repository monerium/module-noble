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
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/v2/x/florin"
	"github.com/monerium/module-noble/v2/x/florin/keeper"
	"github.com/monerium/module-noble/v2/x/florin/types"
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

	k := keeper.NewKeeper(runtime.NewKVStoreService(key), account, bank)

	bank = bank.WithSendCoinsRestriction(k.SendRestrictionFn)
	k.SetBankKeeper(bank)

	ctx := testutil.DefaultContext(key, tkey)
	florin.InitGenesis(ctx, k, *types.DefaultGenesisState())
	ctx.KVStore(key).Delete(types.MaxMintAllowanceKey("ueure"))

	return k, ctx
}

//

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("noble", "noblepub")
	config.Seal()
}
