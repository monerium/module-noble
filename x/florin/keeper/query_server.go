package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/florin/x/florin/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Owner(goCtx context.Context, req *types.QueryOwner) (*types.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryOwnerResponse{
		Owner:        k.GetOwner(ctx),
		PendingOwner: k.GetPendingOwner(ctx),
	}, nil
}

func (k queryServer) MaxMintAllowance(goCtx context.Context, req *types.QueryMaxMintAllowance) (*types.QueryMaxMintAllowanceResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryMaxMintAllowanceResponse{
		MaxMintAllowance: k.GetMaxMintAllowance(ctx),
	}, nil
}
