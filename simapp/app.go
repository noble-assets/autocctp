package simapp

import (
	_ "embed"
	"io"
	"os"
	"path/filepath"

	_ "github.com/circlefin/noble-cctp/x/cctp"
	cctpkeeper "github.com/circlefin/noble-cctp/x/cctp/keeper"
	_ "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory"
	ftfkeeper "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"

	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	_ "cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	_ "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	_ "github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	_ "github.com/cosmos/cosmos-sdk/x/consensus"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	_ "github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	_ "github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	transferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	autocctp "autocctp.dev"
	autocctpkeeper "autocctp.dev/keeper"
)

var DefaultNodeHome string

//go:embed app.yaml
var AppConfigYAML []byte

var (
	_ runtime.AppI            = (*SimApp)(nil)
	_ servertypes.Application = (*SimApp)(nil)
)

// SimApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions.
type SimApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// Cosmos Modules
	AccountKeeper   authkeeper.AccountKeeper
	BankKeeper      bankkeeper.Keeper
	ConsensusKeeper consensuskeeper.Keeper
	ParamsKeeper    paramskeeper.Keeper
	StakingKeeper   *stakingkeeper.Keeper
	UpgradeKeeper   *upgradekeeper.Keeper

	// IBC Modules
	CapabilityKeeper *capabilitykeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper
	TransferKeeper   transferkeeper.Keeper

	// Circle Modules
	CCTPKeeper *cctpkeeper.Keeper
	FTFKeeper  *ftfkeeper.Keeper

	// Custom Modules
	AutoCCTPKeeper *autocctpkeeper.Keeper
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".simapp")
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appconfig.LoadYAML(AppConfigYAML),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			},
		),
	)
}

// NewSimApp returns a reference to an initialized SimApp.
func NewSimApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*SimApp, error) {
	var (
		app        = &SimApp{}
		appBuilder *runtime.AppBuilder
	)

	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Supply(
				logger,
				appOpts,
			),
		),
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		// Cosmos Modules
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.ConsensusKeeper,
		&app.ParamsKeeper,
		&app.StakingKeeper,
		&app.UpgradeKeeper,
		// Circle Modules
		&app.CCTPKeeper,
		&app.FTFKeeper,
		// Custom Modules
		&app.AutoCCTPKeeper,
	); err != nil {
		return nil, err
	}

	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	app.RegisterCCTPService()
	if err := app.RegisterIBCModules(); err != nil {
		return nil, err
	}

	// When initializing the upgrade keeper via dependency injection, the
	// initial module version map is created using only the modules that are
	// wired through dependency injection. As a result, any "legacy" modules
	// (those that don't support dependency injection) are excluded. The line
	// below updates the version map to ensure that all modules are included.
	app.UpgradeKeeper.SetInitVersionMap(app.ModuleManager.GetVersionMap())

	anteHandler, err := NewAnteHandler(HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			SignModeHandler: app.txConfig.SignModeHandler(),
			SigGasConsumer:  autocctp.SigVerificationGasConsumer,
		},
		FTFKeeper:  app.FTFKeeper,
		BankKeeper: app.BankKeeper,
	})
	if err != nil {
		return nil, err
	}

	app.SetAnteHandler(anteHandler)
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *SimApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *SimApp) SimulationManager() *module.SimulationManager {
	return nil
}

//

func (app *SimApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	key, _ := app.UnsafeFindStoreKey(storeKey).(*storetypes.KVStoreKey)
	return key
}

func (app *SimApp) GetMemKey(memKey string) *storetypes.MemoryStoreKey {
	key, _ := app.UnsafeFindStoreKey(memKey).(*storetypes.MemoryStoreKey)
	return key
}

func (app *SimApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *SimApp) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// RegisterCCTPService is a method used to register the CCTP message and query servers into the
// AutoCCTP keeper after building the app.
func (app *SimApp) RegisterCCTPService() {
	cctpMsgServer := cctpkeeper.NewMsgServerImpl(app.CCTPKeeper)

	app.AutoCCTPKeeper.SetCCTPService(cctpMsgServer, app.CCTPKeeper)
}
