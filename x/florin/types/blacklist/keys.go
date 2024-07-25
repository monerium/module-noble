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

package blacklist

const SubmoduleName = "florin-blacklist"

var (
	OwnerKey        = []byte("blacklist/owner")
	PendingOwnerKey = []byte("blacklist/pending_owner")
	AdminPrefix     = []byte("blacklist/admin/")
	AdversaryPrefix = []byte("blacklist/adversary/")
)

func AdminKey(address string) []byte {
	return append(AdminPrefix, []byte(address)...)
}

func AdversaryKey(address string) []byte {
	return append(AdversaryPrefix, []byte(address)...)
}
