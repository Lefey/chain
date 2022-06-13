package v0_5_0

import (
	"fmt"
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func createUnbondingParameters(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	// init param
	registryKeeper.ParamStore().Set(ctx, types.KeyUnbondingStakingTime, types.DefaultUnbondingStakingTime)

	// init param
	registryKeeper.ParamStore().Set(ctx, types.KeyUnbondingDelegationTime, types.DefaultUnbondingDelegationTime)
}

func createProposalIndex(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	fmt.Printf("%sCreating thrid proposal index\n", MigrationLoggerPrefix)

	// Set all delegators again to create the index
	proposals := registryKeeper.GetAllProposal(ctx)
	for index, proposal := range proposals {

		registryKeeper.SetProposal(ctx, proposal)

		if index%1000 == 0 {
			fmt.Printf("%sProposals processed: %d\n", MigrationLoggerPrefix, index)
		}
	}

	fmt.Printf("%sFinished index creation\n", MigrationLoggerPrefix)
}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		createUnbondingParameters(registryKeeper, ctx)

		createProposalIndex(registryKeeper, ctx)

		// Return.
		return vm, nil
	}
}
