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

import "context"

//

func (k *Keeper) GetBlacklistOwner(ctx context.Context) string {
	owner, _ := k.BlacklistOwner.Get(ctx)
	return owner
}

func (k *Keeper) SetBlacklistOwner(ctx context.Context, owner string) error {
	return k.BlacklistOwner.Set(ctx, owner)
}

//

func (k *Keeper) DeleteBlacklistPendingOwner(ctx context.Context) error {
	return k.BlacklistPendingOwner.Remove(ctx)
}

func (k *Keeper) GetBlacklistPendingOwner(ctx context.Context) string {
	pendingOwner, _ := k.BlacklistPendingOwner.Get(ctx)
	return pendingOwner
}

func (k *Keeper) SetBlacklistPendingOwner(ctx context.Context, pendingOwner string) error {
	return k.BlacklistPendingOwner.Set(ctx, pendingOwner)
}

//

func (k *Keeper) DeleteBlacklistAdmin(ctx context.Context, admin string) error {
	return k.BlacklistAdmins.Remove(ctx, admin)
}

func (k *Keeper) GetBlacklistAdmins(ctx context.Context) (admins []string) {
	_ = k.BlacklistAdmins.Walk(ctx, nil, func(admin string) (bool, error) {
		admins = append(admins, admin)
		return false, nil
	})
	return
}

func (k *Keeper) IsBlacklistAdmin(ctx context.Context, admin string) bool {
	isAdmin, _ := k.BlacklistAdmins.Has(ctx, admin)
	return isAdmin
}

func (k *Keeper) SetBlacklistAdmin(ctx context.Context, admin string) error {
	return k.BlacklistAdmins.Set(ctx, admin)
}

//

func (k *Keeper) DeleteAdversary(ctx context.Context, address string) error {
	return k.Adversaries.Remove(ctx, address)
}

func (k *Keeper) GetAdversaries(ctx context.Context) (adversaries []string) {
	_ = k.Adversaries.Walk(ctx, nil, func(adversary string) (bool, error) {
		adversaries = append(adversaries, adversary)
		return false, nil
	})
	return
}

func (k *Keeper) IsAdversary(ctx context.Context, address string) bool {
	isAdversary, _ := k.Adversaries.Has(ctx, address)
	return isAdversary
}

func (k *Keeper) SetAdversary(ctx context.Context, address string) error {
	return k.Adversaries.Set(ctx, address)
}
