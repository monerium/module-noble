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

import "cosmossdk.io/errors"

var (
	Codespace = "florin/blacklist"

	ErrNoOwner             = errors.Register(Codespace, 1, "there is no blacklist owner")
	ErrSameOwner           = errors.Register(Codespace, 2, "provided owner is the current owner")
	ErrInvalidOwner        = errors.Register(Codespace, 3, "signer is not blacklist owner")
	ErrNoPendingOwner      = errors.Register(Codespace, 4, "there is no blacklist pending owner")
	ErrInvalidPendingOwner = errors.Register(Codespace, 5, "signer is not blacklist pending owner")
	ErrInvalidAdmin        = errors.Register(Codespace, 6, "signer is not a blacklist admin")
)
