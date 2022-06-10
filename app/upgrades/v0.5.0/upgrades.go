package v0_5_0

import (
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		// init param
		registryKeeper.ParamStore().Set(ctx, types.KeyUnbondingStakingTime, types.DefaultUnbondingStakingTime)

		// init param
		registryKeeper.ParamStore().Set(ctx, types.KeyUnbondingDelegationTime, types.DefaultUnbondingDelegationTime)

		// Return.
		return vm, nil
	}
}
