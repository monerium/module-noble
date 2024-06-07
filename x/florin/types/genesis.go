package types

import "github.com/noble-assets/florin/x/florin/types/blacklist"

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		BlacklistState: blacklist.DefaultGenesisState(),
	}
}

func (gs *GenesisState) Validate() error {
	if err := gs.BlacklistState.Validate(); err != nil {
		return err
	}

	return nil
}
