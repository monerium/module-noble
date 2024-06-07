package florin

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/keeper"
	"github.com/noble-assets/florin/x/florin/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genesis types.GenesisState) {
	k.SetBlacklistOwner(ctx, genesis.BlacklistState.Owner)
	k.SetBlacklistPendingOwner(ctx, genesis.BlacklistState.PendingOwner)
	for _, admin := range genesis.BlacklistState.Admins {
		k.SetBlacklistAdmin(ctx, admin)
	}
	for _, adversary := range genesis.BlacklistState.Adversaries {
		k.SetAdversary(ctx, adversary)
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		BlacklistState: blacklist.GenesisState{
			Owner:        k.GetBlacklistOwner(ctx),
			PendingOwner: k.GetBlacklistPendingOwner(ctx),
			// TODO: Admins and adversaries
		},
	}
}
