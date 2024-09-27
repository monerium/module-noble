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

package utils

import (
	"cosmossdk.io/store/rootmulti"
	"cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetKVStore retrieves the KVStore for the specified module from the context.
func GetKVStore(ctx sdk.Context, moduleName string) types.KVStore {
	return ctx.KVStore(ctx.MultiStore().(*rootmulti.Store).StoreKeysByName()[moduleName])
}
