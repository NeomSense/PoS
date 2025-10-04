package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/NeomSense/PoS/x/pos/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	// External keepers
	stakingKeeper  types.StakingKeeper
	slashingKeeper types.SlashingKeeper

	Schema         collections.Schema
	Params         collections.Item[types.Params]
	Records        collections.Map[string, types.Record]
	ValidatorStats collections.Map[string, types.ValidatorRecordStats]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	stakingKeeper types.StakingKeeper,
	slashingKeeper types.SlashingKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:   storeService,
		cdc:            cdc,
		addressCodec:   addressCodec,
		authority:      authority,
		stakingKeeper:  stakingKeeper,
		slashingKeeper: slashingKeeper,

		Params: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		Records: collections.NewMap(
			sb,
			types.RecordsKey,
			"records",
			collections.StringKey,
			codec.CollValue[types.Record](cdc),
		),
		ValidatorStats: collections.NewMap(
			sb,
			types.ValidatorStatsKey,
			"validator_stats",
			collections.StringKey,
			codec.CollValue[types.ValidatorRecordStats](cdc),
		),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
