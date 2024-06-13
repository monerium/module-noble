package keeper_test

import (
	"testing"

	"github.com/noble-assets/florin/utils"
	"github.com/noble-assets/florin/utils/mocks"
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
	keeper.SetMintAllowance(ctx, minter1.Address, One)
	keeper.SetMintAllowance(ctx, minter2.Address, One.MulRaw(2))

	// ACT: Attempt to get mint allowances.
	res = keeper.GetMintAllowances(ctx)
	// ASSERT: The action should've succeeded.
	require.Len(t, res, 2)
	require.Equal(t, One.String(), res[minter1.Address])
	require.Equal(t, One.MulRaw(2).String(), res[minter2.Address])
}
