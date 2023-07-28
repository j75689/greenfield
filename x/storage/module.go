package storage

import (
	"context"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/bnb-chain/greenfield/x/storage/client/cli"
	"github.com/bnb-chain/greenfield/x/storage/keeper"
	v2 "github.com/bnb-chain/greenfield/x/storage/keeper/v2"
	"github.com/bnb-chain/greenfield/x/storage/types"
	typesV2 "github.com/bnb-chain/greenfield/x/storage/types/v2"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = &AppModuleBasic{}
)

// ----------------------------------------------------------------------------
// AppModuleBasic
// ----------------------------------------------------------------------------

// AppModuleBasic implements the AppModuleBasic interface that defines the independent methods a Cosmos SDK module needs to implement.
type AppModuleBasic struct {
	cdc     codec.BinaryCodec
	version uint64

	// cached registry reference for upgrading module
	clientCtx client.Context
	mux       *runtime.ServeMux
}

func NewAppModuleBasic(cdc codec.BinaryCodec) *AppModuleBasic {
	return &AppModuleBasic{
		cdc:     cdc,
		version: types.ModuleVersion,
	}
}

// Name returns the name of the module as a string
func (*AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the amino codec for the module, which is used to marshal and unmarshal structs to/from []byte in order to persist them in the module's KVStore
func (module *AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterCodec(cdc)
}

// RegisterInterfaces registers a module's interface types and their concrete implementations as proto.Message
func (module *AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// DefaultGenesis returns a default GenesisState for the module, marshalled to json.RawMessage. The default GenesisState need to be defined by the module developer and is primarily used for testing
func (*AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis used to validate the GenesisState, given in its json.RawMessage form
func (*AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the module
func (module *AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	// nolint: errcheck
	types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	module.mux = mux
	module.clientCtx = clientCtx
}

// GetTxCmd returns the root Tx command for the module. The subcommands of this root command are used by end-users to generate new transactions containing messages defined in the module
func (module *AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns the root query command for the module. The subcommands of this root command are used by end-users to generate new queries to the subset of the state defined by the module
func (module *AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd(types.StoreKey)
}

// ----------------------------------------------------------------------------
// AppModule
// ----------------------------------------------------------------------------

// AppModule implements the AppModule interface that defines the inter-dependent methods that modules need to implement
type AppModule struct {
	*AppModuleBasic

	keeper        keeper.Keeper
	keeperV2      v2.Keeper
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	spKeeper      types.SpKeeper
}

func NewAppModule(
	cdc codec.Codec,
	keeper keeper.Keeper,
	keeperV2 v2.Keeper,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	spKeeper types.SpKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: NewAppModuleBasic(cdc),
		keeper:         keeper,
		keeperV2:       keeperV2,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		spKeeper:       spKeeper,
	}
}

// MigrateToV2 migrates the module from v1 to v2
func (am *AppModule) MigrateToV2(cfg module.Configurator, interfaceRegistry cdctypes.InterfaceRegistry) error {
	am.RegisterAminoAndInterfacesV2(interfaceRegistry)
	am.RegisterServicesV2(cfg)
	// TODO: confirm if we need to register GRPC Gateway routes for v2
	// if err := am.RegisterGRPCGatewayRoutesV2(); err != nil {
	// 	return err
	// }
	v1GroupApp := keeper.NewGroupApp(am.keeper)
	v2GroupApp := v2.NewGroupApp(v1GroupApp, &am.keeperV2)
	if err := am.keeper.GetCrossChainKeeper().MigrateChanel(types.GroupChannel, types.GroupChannelId, v2GroupApp); err != nil {
		return err
	}

	am.version = typesV2.ModuleVersion
	return nil
}

// RegisterAminoAndInterfaces registers a module's interface types and its implementations to the provided protobuf Any.
func (am *AppModule) RegisterAminoAndInterfacesV2(interfaceRegistry cdctypes.InterfaceRegistry) {
	typesV2.RegisterInterfaces(interfaceRegistry)
}

// RegisterGRPCGatewayRoutesV2 registers the gRPC Gateway routes for the module
func (am *AppModule) RegisterGRPCGatewayRoutesV2() error {
	return typesV2.RegisterQueryHandlerClient(context.Background(), am.AppModuleBasic.mux, typesV2.NewQueryClient(am.AppModuleBasic.clientCtx))
}

// RegisterServicesV2 registers a gRPC query service to respond to the module-specific gRPC queries
func (am *AppModule) RegisterServicesV2(cfg module.Configurator) {
	typesV2.RegisterMsgServer(cfg.MsgServer(), v2.NewMsgServerImpl(am.keeperV2, keeper.MsgServer{Keeper: am.keeper}))
	typesV2.RegisterQueryServer(cfg.QueryServer(), am.keeperV2)
}

// RegisterServices registers a gRPC query service to respond to the module-specific gRPC queries
func (am *AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// RegisterInvariants registers the invariants of the module. If an invariant deviates from its predicted value, the InvariantRegistry triggers appropriate logic (most often the chain will be halted)
func (am *AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs the module's genesis initialization. It returns no validator updates.
func (am *AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, gs json.RawMessage) []abci.ValidatorUpdate {
	var genState types.GenesisState
	// Initialize global index to index in genesis state
	cdc.MustUnmarshalJSON(gs, &genState)

	InitGenesis(ctx, am.keeper, genState)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis state as raw JSON bytes.
func (am *AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	genState := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(genState)
}

// ConsensusVersion is a sequence number for state-breaking change of the module. It should be incremented on each consensus-breaking change introduced by the module. To avoid wrong/empty versions, the initial version should be set to 1
func (am *AppModule) ConsensusVersion() uint64 { return am.version }

// BeginBlock contains the logic that is automatically triggered at the beginning of each block
func (am *AppModule) BeginBlock(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.keeper)
}

// EndBlock contains the logic that is automatically triggered at the end of each block
func (am *AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
