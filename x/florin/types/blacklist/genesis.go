package blacklist

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func (gs *GenesisState) Validate() error {
	if gs.Owner != "" {
		if _, err := sdk.AccAddressFromBech32(gs.Owner); err != nil {
			return fmt.Errorf("invalid blacklist owner address (%s): %s", gs.Owner, err)
		}
	}

	if gs.PendingOwner != "" {
		if _, err := sdk.AccAddressFromBech32(gs.PendingOwner); err != nil {
			return fmt.Errorf("invalid pending blacklist owner address (%s): %s", gs.PendingOwner, err)
		}
	}

	for _, admin := range gs.Admins {
		if _, err := sdk.AccAddressFromBech32(admin); err != nil {
			return fmt.Errorf("invalid admin address (%s): %s", admin, err)
		}
	}

	for _, adversary := range gs.Adversaries {
		if _, err := sdk.AccAddressFromBech32(adversary); err != nil {
			return fmt.Errorf("invalid adversary address (%s): %s", adversary, err)
		}
	}

	return nil
}
