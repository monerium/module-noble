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

func (k blacklistQueryServer) Admins(goCtx context.Context, req *blacklist.QueryAdmins) (*blacklist.QueryAdminsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &blacklist.QueryAdminsResponse{
		Admins: k.GetBlacklistAdmins(ctx),
	}, nil
}

func (k blacklistQueryServer) Adversaries(goCtx context.Context, req *blacklist.QueryAdversaries) (*blacklist.QueryAdversariesResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &blacklist.QueryAdversariesResponse{
		Adversaries: k.GetAdversaries(ctx),
	}, nil
}
