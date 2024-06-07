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
