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

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAcceptOwnership{}, "florin/blacklist/AcceptOwnership", nil)
	cdc.RegisterConcrete(&MsgAddAdminAccount{}, "florin/blacklist/AddAdminAccount", nil)
	cdc.RegisterConcrete(&MsgBan{}, "florin/blacklist/Ban", nil)
	cdc.RegisterConcrete(&MsgRemoveAdminAccount{}, "florin/blacklist/RemoveAdminAccount", nil)
	cdc.RegisterConcrete(&MsgTransferOwnership{}, "florin/blacklist/TransferOwnership", nil)
	cdc.RegisterConcrete(&MsgUnban{}, "florin/blacklist/Unban", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAcceptOwnership{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddAdminAccount{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgBan{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveAdminAccount{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransferOwnership{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUnban{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
