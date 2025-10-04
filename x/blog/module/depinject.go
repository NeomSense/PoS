package blog

import (
	"fmt"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/NeomSense/PoS/x/blog/keeper"
	"github.com/NeomSense/PoS/x/blog/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper types.AuthKeeper
	BankKeeper types.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	BlogKeeper keeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	var authorityBytes = authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		// Accept either a module name or a bech32 address from config
		authorityBytes = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	// Convert sdk.AccAddress -> bech32 string to match keeper.NewKeeper(authority string)
	authorityStr, err := in.AddressCodec.BytesToString(authorityBytes)
	if err != nil {
		panic(fmt.Sprintf("invalid authority address in depinject: %v", err))
	}

	k, err := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.AddressCodec,
		authorityStr,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create blog keeper: %v", err))
	}

	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{BlogKeeper: k, Module: m}
}
