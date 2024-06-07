package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	blacklist.RegisterLegacyAminoCodec(cdc)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	blacklist.RegisterInterfaces(registry)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
