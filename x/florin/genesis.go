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

	k.SetOwner(ctx, genesis.Owner)
	k.SetPendingOwner(ctx, genesis.PendingOwner)
	for _, system := range genesis.Systems {
		k.SetSystem(ctx, system)
	}
	for _, admin := range genesis.Admins {
		k.SetAdmin(ctx, admin)
	}
	for address, rawAllowance := range genesis.MintAllowances {
		allowance, _ := sdk.NewIntFromString(rawAllowance)
		k.SetMintAllowance(ctx, address, allowance)
	}
	k.SetMaxMintAllowance(ctx, genesis.MaxMintAllowance)
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		BlacklistState: blacklist.GenesisState{
			Owner:        k.GetBlacklistOwner(ctx),
			PendingOwner: k.GetBlacklistPendingOwner(ctx),
			Admins:       k.GetBlacklistAdmins(ctx),
			Adversaries:  k.GetAdversaries(ctx),
		},
		Owner:            k.GetOwner(ctx),
		PendingOwner:     k.GetPendingOwner(ctx),
		Systems:          k.GetSystems(ctx),
		Admins:           k.GetAdmins(ctx),
		MintAllowances:   k.GetMintAllowances(ctx),
		MaxMintAllowance: k.GetMaxMintAllowance(ctx),
	}
}
