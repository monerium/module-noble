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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/monerium/module-noble/utils"
	"github.com/monerium/module-noble/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestSendRestriction(t *testing.T) {
	keeper, ctx := mocks.FlorinKeeper()
	sender, recipient := utils.TestAccount(), utils.TestAccount()
	ONE := sdk.NewCoin("ueure", sdk.NewInt(1_000_000_000_000_000_000))

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
