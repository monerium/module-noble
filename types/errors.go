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

import "cosmossdk.io/errors"

var (
	ErrNoAuthority           = errors.Register(ModuleName, 1, "there is no authority")
	ErrInvalidAuthority      = errors.Register(ModuleName, 2, "signer is not authority")
	ErrInvalidDenom          = errors.Register(ModuleName, 3, "denom is already in use")
	ErrNoOwner               = errors.Register(ModuleName, 4, "there is no owner")
	ErrSameOwner             = errors.Register(ModuleName, 5, "provided owner is the current owner")
	ErrInvalidOwner          = errors.Register(ModuleName, 6, "signer is not owner")
	ErrNoPendingOwner        = errors.Register(ModuleName, 7, "there is no pending owner")
	ErrInvalidPendingOwner   = errors.Register(ModuleName, 8, "signer is not pending owner")
	ErrInvalidSystem         = errors.Register(ModuleName, 9, "signer is not a system")
	ErrInvalidAdmin          = errors.Register(ModuleName, 10, "signer is not an admin")
	ErrInvalidAllowance      = errors.Register(ModuleName, 11, "allowance cannot be negative or greater than max")
	ErrInsufficientAllowance = errors.Register(ModuleName, 12, "insufficient allowance")
	ErrInvalidPubKey         = errors.Register(ModuleName, 13, "invalid public key")
	ErrInvalidSignature      = errors.Register(ModuleName, 14, "invalid signature")
)
