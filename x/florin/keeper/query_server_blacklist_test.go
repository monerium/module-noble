package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/florin/utils"
	"github.com/noble-assets/florin/utils/mocks"
	"github.com/noble-assets/florin/x/florin/keeper"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
	"github.com/stretchr/testify/require"
)

func TestBlacklistOwnerQuery(t *testing.T) {
	k, ctx := mocks.FlorinKeeper()
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewBlacklistQueryServer(k)

	// ACT: Attempt to query owner with invalid request.
	_, err := server.Owner(goCtx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetBlacklistOwner(ctx, owner.Address)

	// ACT: Attempt to query owner.
	res, err := server.Owner(goCtx, &blacklist.QueryOwner{})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetBlacklistPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to query owner.
	res, err = server.Owner(goCtx, &blacklist.QueryOwner{})
	// ASSERT: The query should've succeeded, with pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Equal(t, pendingOwner.Address, res.PendingOwner)
}
