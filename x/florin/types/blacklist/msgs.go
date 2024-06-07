package blacklist

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

//

var _ legacytx.LegacyMsg = &MsgAcceptOwnership{}

func (msg *MsgAcceptOwnership) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	return nil
}

func (msg *MsgAcceptOwnership) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAcceptOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAcceptOwnership) Route() string { return SubmoduleName }

func (*MsgAcceptOwnership) Type() string { return "florin/blacklist/AcceptOwnership" }

//

var _ legacytx.LegacyMsg = &MsgAddAdminAccount{}

func (msg *MsgAddAdminAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Account); err != nil {
		return fmt.Errorf("invalid account address (%s): %w", msg.Account, err)
	}

	return nil
}

func (msg *MsgAddAdminAccount) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddAdminAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAddAdminAccount) Route() string { return SubmoduleName }

func (*MsgAddAdminAccount) Type() string { return "florin/blacklist/AddAdminAccount" }

//

var _ legacytx.LegacyMsg = &MsgBan{}

func (msg *MsgBan) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Adversary); err != nil {
		return fmt.Errorf("invalid adversary address (%s): %w", msg.Adversary, err)
	}

	return nil
}

func (msg *MsgBan) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgBan) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgBan) Route() string { return SubmoduleName }

func (*MsgBan) Type() string { return "florin/blacklist/Ban" }

//

var _ legacytx.LegacyMsg = &MsgRemoveAdminAccount{}

func (msg *MsgRemoveAdminAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Account); err != nil {
		return fmt.Errorf("invalid account address (%s): %w", msg.Account, err)
	}

	return nil
}

func (msg *MsgRemoveAdminAccount) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveAdminAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRemoveAdminAccount) Route() string { return SubmoduleName }

func (*MsgRemoveAdminAccount) Type() string { return "florin/blacklist/RemoveAdminAccount" }

//

var _ legacytx.LegacyMsg = &MsgTransferOwnership{}

func (msg *MsgTransferOwnership) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.NewOwner); err != nil {
		return fmt.Errorf("invalid new owner address (%s): %w", msg.NewOwner, err)
	}

	return nil
}

func (msg *MsgTransferOwnership) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgTransferOwnership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgTransferOwnership) Route() string { return SubmoduleName }

func (*MsgTransferOwnership) Type() string { return "florin/blacklist/TransferOwnership" }

//

var _ legacytx.LegacyMsg = &MsgUnban{}

func (msg *MsgUnban) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Friend); err != nil {
		return fmt.Errorf("invalid friend address (%s): %w", msg.Friend, err)
	}

	return nil
}

func (msg *MsgUnban) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgUnban) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgUnban) Route() string { return SubmoduleName }

func (*MsgUnban) Type() string { return "florin/blacklist/Unban" }
