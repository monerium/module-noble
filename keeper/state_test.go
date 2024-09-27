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

package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/monerium/module-noble/v2/types"
	"github.com/monerium/module-noble/v2/utils"
	"github.com/monerium/module-noble/v2/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetMintAllowances(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()

	// ACT: Attempt to get mint allowances with no state.
	res := k.GetMintAllowances(ctx)
	// ASSERT: The action should've succeeded, returns empty.
	require.Empty(t, res)

	// ARRANGE: Set mint allowances in state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	err := k.SetMintAllowance(ctx, "ueure", minter1.Address, One)
	require.NoError(t, err)
	err = k.SetMintAllowance(ctx, "ueure", minter2.Address, One.MulRaw(2))
	require.NoError(t, err)

	// ACT: Attempt to get mint allowances.
	res = k.GetMintAllowances(ctx)
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Len(t, res, 2)
	require.Contains(t, res, types.Allowance{
		Denom:     "ueure",
		Address:   minter1.Address,
		Allowance: One,
	})
	require.Contains(t, res, types.Allowance{
		Denom:     "ueure",
		Address:   minter2.Address,
		Allowance: One.MulRaw(2),
	})

	// ARRANGE: Set invalid mint allowance
	key := collections.Join("ueure", "address")
	_ = k.MintAllowance.Set(ctx, key, []byte("panic"))

	// ACT: Attempt to get mint allowances.
	allowances := k.GetMintAllowancesByDenom(ctx, "ueure")
	// ASSERT: The action should've succeeded, returns empty.
	require.Empty(t, allowances)
}
