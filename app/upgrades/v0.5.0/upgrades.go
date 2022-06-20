package v0_5_0

import (
	"fmt"
	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
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

func migrateIBCDenoms(ctx sdk.Context, transferKeeper *ibctransferkeeper.Keeper) {
	var newTraces []ibctransfertypes.DenomTrace

	transferKeeper.IterateDenomTraces(ctx,
		func(dt ibctransfertypes.DenomTrace) bool {
			newTrace := ibctransfertypes.ParseDenomTrace(dt.GetFullDenomPath())

			if err := newTrace.Validate(); err == nil && !equalTraces(newTrace, dt) {
				newTraces = append(newTraces, newTrace)
			}

			return false
		},
	)

	for _, nt := range newTraces {
		transferKeeper.SetDenomTrace(ctx, nt)
	}
}

func CreateUpgradeHandler(
	registryKeeper *registrykeeper.Keeper,
	transferKeeper *ibctransferkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		createUnbondingParameters(registryKeeper, ctx)

		createProposalIndex(registryKeeper, ctx)

		migrateIBCDenoms(ctx, transferKeeper)

		// Return.
		return vm, nil
	}
}

func equalTraces(dtA, dtB ibctransfertypes.DenomTrace) bool {
	return dtA.BaseDenom == dtB.BaseDenom && dtA.Path == dtB.Path
}
