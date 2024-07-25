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

package types

const ModuleName = "florin"

var (
	AuthorityKey           = []byte("authority")
	AllowedDenomPrefix     = []byte("allowed_denom/")
	OwnerPrefix            = []byte("owner/")
	PendingOwnerPrefix     = []byte("pending_owner/")
	SystemPrefix           = []byte("system/")
	AdminPrefix            = []byte("admin/")
	MintAllowancePrefix    = []byte("mint_allowance/")
	MaxMintAllowancePrefix = []byte("max_mint_allowance/")
)

func AllowedDenomKey(denom string) []byte {
	return append(AllowedDenomPrefix, []byte(denom)...)
}

func OwnerKey(denom string) []byte {
	return append(OwnerPrefix, []byte(denom)...)
}

func PendingOwnerKey(denom string) []byte {
	return append(PendingOwnerPrefix, []byte(denom)...)
}

func SystemKey(denom string, address string) []byte {
	return append(append(SystemPrefix, []byte(denom)...), []byte(address)...)
}

func AdminKey(denom string, address string) []byte {
	return append(append(AdminPrefix, []byte(denom)...), []byte(address)...)
}

func MintAllowanceKey(denom string, address string) []byte {
	return append(append(MintAllowancePrefix, []byte(denom)...), []byte(address)...)
}

func MaxMintAllowanceKey(denom string) []byte {
	return append(MaxMintAllowancePrefix, []byte(denom)...)
}
