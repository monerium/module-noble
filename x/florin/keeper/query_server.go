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

func (k queryServer) Systems(goCtx context.Context, req *types.QuerySystems) (*types.QuerySystemsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QuerySystemsResponse{
		Systems: k.GetSystems(ctx),
	}, nil
}

func (k queryServer) Admins(goCtx context.Context, req *types.QueryAdmins) (*types.QueryAdminsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAdminsResponse{
		Admins: k.GetAdmins(ctx),
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

func (k queryServer) MintAllowance(goCtx context.Context, req *types.QueryMintAllowance) (*types.QueryMintAllowanceResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryMintAllowanceResponse{
		Allowance: k.GetMintAllowance(ctx, req.Account),
	}, nil
}
