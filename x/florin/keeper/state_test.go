package keeper_test

import (
	"testing"

	"github.com/noble-assets/florin/utils"
	"github.com/noble-assets/florin/utils/mocks"
	"github.com/noble-assets/florin/x/florin/types"
	"github.com/stretchr/testify/require"
)

func TestGetMintAllowances(t *testing.T) {
	keeper, ctx := mocks.FlorinKeeper()

	// ACT: Attempt to get mint allowances with no state.
	res := keeper.GetMintAllowances(ctx)
	// ASSERT: The action should've succeeded, returns empty.
	require.Empty(t, res)

	// ARRANGE: Set mint allowances in state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	keeper.SetMintAllowance(ctx, "ueure", minter1.Address, One)
	keeper.SetMintAllowance(ctx, "ueure", minter2.Address, One.MulRaw(2))

	// ACT: Attempt to get mint allowances.
	res = keeper.GetMintAllowances(ctx)
	// ASSERT: The action should've succeeded.
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
}
