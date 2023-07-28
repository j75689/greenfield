package app

import (
	"fmt"

	storagemodule "github.com/bnb-chain/greenfield/x/storage"
	storagetypesV1 "github.com/bnb-chain/greenfield/x/storage/types"
	storagetypesV2 "github.com/bnb-chain/greenfield/x/storage/types/v2"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func (app *App) RegisterUpgradeHandlers(chainID string, serverCfg *serverconfig.Config) error {
	// Register the plans from server config
	err := app.UpgradeKeeper.RegisterUpgradePlan(chainID, serverCfg.Upgrade)
	if err != nil {
		return err
	}

	// Register the upgrade handlers here
	// app.registerPublicDelegationUpgradeHandler()
	app.registerBEP1001UpgradeHandler()

	return nil
}

// registerPublicDelegationUpgradeHandler registers the upgrade handlers for the public delegation upgrade.
// func (app *App) registerPublicDelegationUpgradeHandler() {
// 	// Register the upgrade handler
// 	app.UpgradeKeeper.SetUpgradeHandler(upgradetypes.EnablePublicDelegationUpgrade,
// 		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
// 			app.Logger().Info("upgrade to ", plan.Name)
// 			return fromVM, nil
// 		})

// 	// Register the upgrade initializer
// 	app.UpgradeKeeper.SetUpgradeInitializer(upgradetypes.EnablePublicDelegationUpgrade,
// 		func() error {
// 			app.Logger().Info("Init enable public delegation upgrade")
// 			return nil
// 		},
// 	)
// }

// registerBEP1001UpgradeHandler registers the upgrade handlers for BEP1001.
func (app *App) registerBEP1001UpgradeHandler() {
	// Register the upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(upgradetypes.BEP1001,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			app.Logger().Info("processing upgrade handler", "name", plan.Name, "info", plan.Info)
			fromVM[storagetypesV1.ModuleName] = storagetypesV2.ModuleVersion
			return fromVM, nil
		})

	// Register the upgrade initializer
	app.UpgradeKeeper.SetUpgradeInitializer(upgradetypes.BEP1001,
		func() error {
			app.Logger().Info("processing upgrade initializer", "name", upgradetypes.BEP1001)
			storageModule, ok := app.mm.Modules[storagetypesV1.ModuleName].(storagemodule.AppModule)
			if !ok {
				return fmt.Errorf("storage module not found")
			}
			err := storageModule.MigrateToV2(app.configurator, app.InterfaceRegistry())
			if err != nil {
				return err
			}
			return nil
		},
	)
}
