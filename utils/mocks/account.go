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
