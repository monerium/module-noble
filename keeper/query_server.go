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

package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/monerium/module-noble/v2/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Authority(_ context.Context, req *types.QueryAuthority) (*types.QueryAuthorityResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryAuthorityResponse{Authority: k.authority}, nil
}

func (k queryServer) AllowedDenoms(goCtx context.Context, req *types.QueryAllowedDenoms) (*types.QueryAllowedDenomsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAllowedDenomsResponse{
		AllowedDenoms: k.GetAllowedDenoms(ctx),
	}, nil
}

func (k queryServer) Owners(goCtx context.Context, req *types.QueryOwners) (*types.QueryOwnersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryOwnersResponse{
		Owners:        k.GetOwners(ctx),
		PendingOwners: k.GetPendingOwners(ctx),
	}, nil
}

func (k queryServer) Owner(goCtx context.Context, req *types.QueryOwner) (*types.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsAllowedDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", req.Denom)
	}

	return &types.QueryOwnerResponse{
		Owner:        k.GetOwner(ctx, req.Denom),
		PendingOwner: k.GetPendingOwner(ctx, req.Denom),
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

func (k queryServer) SystemsByDenom(goCtx context.Context, req *types.QuerySystemsByDenom) (*types.QuerySystemsByDenomResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsAllowedDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", req.Denom)
	}

	return &types.QuerySystemsByDenomResponse{
		Systems: k.GetSystemsByDenom(ctx, req.Denom),
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

func (k queryServer) AdminsByDenom(goCtx context.Context, req *types.QueryAdminsByDenom) (*types.QueryAdminsByDenomResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsAllowedDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", req.Denom)
	}

	return &types.QueryAdminsByDenomResponse{
		Admins: k.GetAdminsByDenom(ctx, req.Denom),
	}, nil
}

func (k queryServer) MaxMintAllowances(goCtx context.Context, req *types.QueryMaxMintAllowances) (*types.QueryMaxMintAllowancesResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryMaxMintAllowancesResponse{
		MaxMintAllowances: k.GetMaxMintAllowances(ctx),
	}, nil
}

func (k queryServer) MaxMintAllowance(goCtx context.Context, req *types.QueryMaxMintAllowance) (*types.QueryMaxMintAllowanceResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsAllowedDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", req.Denom)
	}

	return &types.QueryMaxMintAllowanceResponse{
		MaxMintAllowance: k.GetMaxMintAllowance(ctx, req.Denom),
	}, nil
}

func (k queryServer) MintAllowances(goCtx context.Context, req *types.QueryMintAllowances) (*types.QueryMintAllowancesResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsAllowedDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", req.Denom)
	}

	allowances := make(map[string]string)
	for _, entry := range k.GetMintAllowancesByDenom(ctx, req.Denom) {
		allowances[entry.Address] = entry.Allowance.String()
	}

	return &types.QueryMintAllowancesResponse{Allowances: allowances}, nil
}

func (k queryServer) MintAllowance(goCtx context.Context, req *types.QueryMintAllowance) (*types.QueryMintAllowanceResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.IsAllowedDenom(ctx, req.Denom) {
		return nil, fmt.Errorf("%s is not an allowed denom", req.Denom)
	}

	return &types.QueryMintAllowanceResponse{
		Allowance: k.GetMintAllowance(ctx, req.Denom, req.Account),
	}, nil
}
