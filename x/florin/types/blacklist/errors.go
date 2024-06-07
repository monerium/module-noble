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
