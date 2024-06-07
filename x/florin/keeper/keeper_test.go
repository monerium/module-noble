package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/utils"
	"github.com/noble-assets/florin/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestSendRestriction(t *testing.T) {
	keeper, ctx := mocks.FlorinKeeper()
	sender, recipient := utils.TestAccount(), utils.TestAccount()
	ONE := sdk.NewCoin(keeper.Denom, sdk.NewInt(1_000_000_000_000_000_000))

	// ACT: Attempt transfer with non $EURe coin.
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	_, err := keeper.SendRestrictionFn(
		ctx, sender.Bytes, recipient.Bytes,
		sdk.NewCoins(sdk.NewCoin("uusdc", sdk.NewInt(1_000_000))),
	)
	// ASSERT: The transfer should've succeeded.
	require.NoError(t, err)
	events := ctx.EventManager().Events()
	require.Empty(t, events)

	// ACT: Attempt transfer with friendly sender.
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	_, err = keeper.SendRestrictionFn(
		ctx, sender.Bytes, recipient.Bytes,
		sdk.NewCoins(ONE),
	)
	// ASSERT: The transfer should've succeeded.
	require.NoError(t, err)
	events = ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.Decision", events[0].Type)

	// ARRANGE: Set sender as adversary.
	keeper.SetAdversary(ctx, sender.Address)

	// ACT: Attempt transfer with adversarial sender.
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	_, err = keeper.SendRestrictionFn(
		ctx, sender.Bytes, recipient.Bytes,
		sdk.NewCoins(ONE),
	)
	// ASSERT: The transfer should've failed.
	require.ErrorContains(t, err, "blocked from sending")
	events = ctx.EventManager().Events()
	require.Len(t, events, 1)
	require.Equal(t, "florin.blacklist.v1.Decision", events[0].Type)
}
