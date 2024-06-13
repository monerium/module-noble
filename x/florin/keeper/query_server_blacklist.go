package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

var _ blacklist.QueryServer = &blacklistQueryServer{}

type blacklistQueryServer struct {
	*Keeper
}

func NewBlacklistQueryServer(keeper *Keeper) blacklist.QueryServer {
	return &blacklistQueryServer{Keeper: keeper}
}

func (k blacklistQueryServer) Owner(goCtx context.Context, req *blacklist.QueryOwner) (*blacklist.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &blacklist.QueryOwnerResponse{
		Owner:        k.GetBlacklistOwner(ctx),
		PendingOwner: k.GetBlacklistPendingOwner(ctx),
	}, nil
}
