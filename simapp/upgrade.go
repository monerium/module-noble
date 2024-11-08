package simapp

import (
	"context"
	"cosmossdk.io/errors"
	"cosmossdk.io/store/rootmulti"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (app *SimApp) RegisterUpgradeHandler() {
	app.UpgradeKeeper.SetUpgradeHandler(
		"v2",
		func(ctx context.Context, _ types.Plan, vm module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			for _, subspace := range app.ParamsKeeper.GetSubspaces() {
				var keyTable paramstypes.KeyTable
				switch subspace.Name() {
				case authtypes.ModuleName:
					keyTable = authtypes.ParamKeyTable() //nolint:staticcheck
				case banktypes.ModuleName:
					keyTable = banktypes.ParamKeyTable() //nolint:staticcheck
				case stakingtypes.ModuleName:
					keyTable = stakingtypes.ParamKeyTable() //nolint:staticcheck
				}

				if !subspace.HasKeyTable() {
					subspace.WithKeyTable(keyTable)
				}
			}

			vm, err := app.ModuleManager.RunMigrations(ctx, app.Configurator(), vm)
			if err != nil {
				return vm, err
			}

			subspace := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()) //nolint:staticcheck
			err = baseapp.MigrateParams(sdkCtx, subspace, app.ConsensusKeeper.ParamsStore)
			if err != nil {
				return nil, errors.Wrap(err, "failed to migrate consensus params")
			}

			return vm, nil
		},
	)

	app.SetStoreLoader(func(ms storetypes.CommitMultiStore) error {
		tmp := ms.(*rootmulti.Store)
		_ = tmp.LoadLatestVersion()

		store := tmp.GetStoreByName(consensustypes.ModuleName)
		if store != nil {
			return baseapp.DefaultStoreLoader(ms)
		}

		return ms.LoadLatestVersionAndUpgrade(&storetypes.StoreUpgrades{
			Added: []string{consensustypes.StoreKey},
		})
	})
}
