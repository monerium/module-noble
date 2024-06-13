package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

func DefaultGenesisState() *GenesisState {
	maxMintAllowance, _ := sdk.NewIntFromString("3000000000000000000000000")

	return &GenesisState{
		BlacklistState:   blacklist.DefaultGenesisState(),
		MaxMintAllowance: maxMintAllowance,
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.BlacklistState.Validate(); err != nil {
		return err
	}

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

	for _, system := range gs.Systems {
		if _, err := sdk.AccAddressFromBech32(system); err != nil {
			return fmt.Errorf("invalid system address (%s): %s", system, err)
		}
	}

	for _, admin := range gs.Admins {
		if _, err := sdk.AccAddressFromBech32(admin); err != nil {
			return fmt.Errorf("invalid admin address (%s): %s", admin, err)
		}
	}

	for address, allowance := range gs.MintAllowances {
		if _, err := sdk.AccAddressFromBech32(address); err != nil {
			return fmt.Errorf("invalid address (%s): %s", address, err)
		}

		if _, ok := sdk.NewIntFromString(allowance); !ok {
			return fmt.Errorf("invalid mint allowance (%s)", allowance)
		}
	}

	return nil
}
