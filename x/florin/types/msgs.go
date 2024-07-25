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

func (*MsgAcceptOwnership) Route() string { return ModuleName }

func (*MsgAcceptOwnership) Type() string { return "florin/AcceptOwnership" }

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

func (*MsgAddAdminAccount) Route() string { return ModuleName }

func (*MsgAddAdminAccount) Type() string { return "florin/AddAdminAccount" }

//

var _ legacytx.LegacyMsg = &MsgAddSystemAccount{}

func (msg *MsgAddSystemAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Account); err != nil {
		return fmt.Errorf("invalid account address (%s): %w", msg.Account, err)
	}

	return nil
}

func (msg *MsgAddSystemAccount) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddSystemAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAddSystemAccount) Route() string { return ModuleName }

func (*MsgAddSystemAccount) Type() string { return "florin/AddSystemAccount" }

//

var _ legacytx.LegacyMsg = &MsgAllowDenom{}

func (msg *MsgAllowDenom) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return fmt.Errorf("invalid owner address (%s): %w", msg.Owner, err)
	}

	return nil
}

func (msg *MsgAllowDenom) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgAllowDenom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgAllowDenom) Route() string { return ModuleName }

func (*MsgAllowDenom) Type() string { return "florin/AllowDenom" }

//

var _ legacytx.LegacyMsg = &MsgBurn{}

func (msg *MsgBurn) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return fmt.Errorf("invalid from address (%s): %w", msg.From, err)
	}

	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount must be positive (%s)", msg.Amount)
	}

	return nil
}

func (msg *MsgBurn) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgBurn) Route() string { return ModuleName }

func (*MsgBurn) Type() string { return "florin/Burn" }

//

var _ legacytx.LegacyMsg = &MsgMint{}

func (msg *MsgMint) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return fmt.Errorf("invalid to address (%s): %w", msg.To, err)
	}

	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount must be positive (%s)", msg.Amount)
	}

	return nil
}

func (msg *MsgMint) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgMint) Route() string { return ModuleName }

func (*MsgMint) Type() string { return "florin/Mint" }

//

var _ legacytx.LegacyMsg = &MsgRecover{}

func (msg *MsgRecover) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.From); err != nil {
		return fmt.Errorf("invalid from address (%s): %w", msg.From, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
		return fmt.Errorf("invalid to address (%s): %w", msg.To, err)
	}

	return nil
}

func (msg *MsgRecover) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRecover) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRecover) Route() string { return ModuleName }

func (*MsgRecover) Type() string { return "florin/Recover" }

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

func (*MsgRemoveAdminAccount) Route() string { return ModuleName }

func (*MsgRemoveAdminAccount) Type() string { return "florin/RemoveAdminAccount" }

//

var _ legacytx.LegacyMsg = &MsgRemoveSystemAccount{}

func (msg *MsgRemoveSystemAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Account); err != nil {
		return fmt.Errorf("invalid account address (%s): %w", msg.Account, err)
	}

	return nil
}

func (msg *MsgRemoveSystemAccount) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgRemoveSystemAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgRemoveSystemAccount) Route() string { return ModuleName }

func (*MsgRemoveSystemAccount) Type() string { return "florin/RemoveSystemAccount" }

//

var _ legacytx.LegacyMsg = &MsgSetMaxMintAllowance{}

func (msg *MsgSetMaxMintAllowance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount must be positive (%s)", msg.Amount)
	}

	return nil
}

func (msg *MsgSetMaxMintAllowance) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgSetMaxMintAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgSetMaxMintAllowance) Route() string { return ModuleName }

func (*MsgSetMaxMintAllowance) Type() string { return "florin/SetMaxMintAllowance" }

//

var _ legacytx.LegacyMsg = &MsgSetMintAllowance{}

func (msg *MsgSetMintAllowance) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return fmt.Errorf("invalid signer address (%s): %w", msg.Signer, err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Account); err != nil {
		return fmt.Errorf("invalid account address (%s): %w", msg.Account, err)
	}

	return nil
}

func (msg *MsgSetMintAllowance) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{signer}
}

func (msg *MsgSetMintAllowance) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (*MsgSetMintAllowance) Route() string { return ModuleName }

func (*MsgSetMintAllowance) Type() string { return "florin/SetMintAllowance" }

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

func (*MsgTransferOwnership) Route() string { return ModuleName }

func (*MsgTransferOwnership) Type() string { return "florin/TransferOwnership" }
