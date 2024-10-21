package simapp

import (
	"encoding/binary"
	"encoding/json"
	"sort"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ConsensusVersion defines the current "upgrade" module consensus version.
const ConsensusVersion = 1

// UpgradeHeight defines the v2 block height.
const UpgradeHeight = 7350000

// DevnetChainID defines the Chain ID of the Florin devnet.
const DevnetChainID = "florin-devnet-1"

var (
	_ module.AppModuleBasic    = UpgradeAppModuleBasic{}
	_ module.AppModule         = UpgradeAppModule{}
	_ module.EndBlockAppModule = UpgradeAppModule{}
)

//

type UpgradeAppModuleBasic struct{}

func NewUpgradeAppModuleBasic() UpgradeAppModuleBasic {
	return UpgradeAppModuleBasic{}
}

func (UpgradeAppModuleBasic) Name() string { return types.ModuleName }

func (UpgradeAppModuleBasic) RegisterLegacyAminoCodec(_ *codec.LegacyAmino) {}

func (UpgradeAppModuleBasic) RegisterInterfaces(_ codectypes.InterfaceRegistry) {}

func (UpgradeAppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage {
	return []byte("{}")
}

func (UpgradeAppModuleBasic) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

func (UpgradeAppModuleBasic) RegisterRESTRoutes(_ client.Context, _ *mux.Router) {}

func (UpgradeAppModuleBasic) RegisterGRPCGatewayRoutes(_ client.Context, _ *runtime.ServeMux) {}

func (UpgradeAppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

func (UpgradeAppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

//

type UpgradeAppModule struct {
	UpgradeAppModuleBasic

	cdc     codec.Codec
	key     sdk.StoreKey
	manager *module.Manager
}

func NewUpgradeAppModule(cdc codec.Codec, key sdk.StoreKey, manager *module.Manager) *UpgradeAppModule {
	return &UpgradeAppModule{
		UpgradeAppModuleBasic: NewUpgradeAppModuleBasic(),
		cdc:                   cdc,
		key:                   key,
		manager:               manager,
	}
}

func (UpgradeAppModule) InitGenesis(_ sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (m UpgradeAppModule) ExportGenesis(_ sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return m.DefaultGenesis(cdc)
}

func (UpgradeAppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

func (UpgradeAppModule) Route() sdk.Route { return sdk.Route{} }

func (UpgradeAppModule) QuerierRoute() string { return types.ModuleName }

func (UpgradeAppModule) LegacyQuerierHandler(_ *codec.LegacyAmino) sdk.Querier { return nil }

func (UpgradeAppModule) RegisterServices(_ module.Configurator) {}

func (UpgradeAppModule) ConsensusVersion() uint64 { return ConsensusVersion }

// SetManager ...
func (m *UpgradeAppModule) SetManager(manager *module.Manager) {
	m.manager = manager
}

// EndBlock is a custom hook that sets the missing but expected upgrade module
// state. This has to be done as SimApp v1 didn't include the upgrade module.
func (m UpgradeAppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	if ctx.ChainID() != DevnetChainID {
		return nil
	}

	if req.Height == UpgradeHeight {
		panic("UPGRADE v2 NEEDED")
	}
	if req.Height != UpgradeHeight-1 {
		return nil
	}

	vm := m.manager.GetVersionMap()
	store := ctx.KVStore(m.key)

	// NOTE: https://github.com/cosmos/cosmos-sdk/blob/v0.45.16/x/upgrade/keeper/keeper.go#L94-L111
	if len(vm) > 0 {
		versionStore := prefix.NewStore(store, []byte{types.VersionMapByte})
		// Even though the underlying store (cachekv) store is sorted, we still
		// prefer a deterministic iteration order of the map, to avoid undesired
		// surprises if we ever change stores.
		sortedModNames := make([]string, 0, len(vm))

		for key := range vm {
			sortedModNames = append(sortedModNames, key)
		}
		sort.Strings(sortedModNames)

		for _, modName := range sortedModNames {
			ver := vm[modName]
			nameBytes := []byte(modName)
			verBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(verBytes, ver)
			versionStore.Set(nameBytes, verBytes)
		}
	}

	store.Set(types.PlanKey(), m.cdc.MustMarshal(&types.Plan{
		Name:   "v2",
		Height: UpgradeHeight,
	}))

	return nil
}
