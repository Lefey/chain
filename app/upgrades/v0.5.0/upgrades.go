package v0_5_0

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"strconv"
	"time"

	registrykeeper "github.com/KYVENetwork/chain/x/registry/keeper"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

func migratePools(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {
	registryKeeper.ParamStore().Set(ctx, types.KeyStorageCost, uint64(50))

	for _, pool := range registryKeeper.GetAllPool(ctx) {
		// deprecate pool versions
		pool.Versions = ""

		// deprecate pool from_height
		pool.BundleProposal.FromHeight = 0

		// set 2.5 $KYVE as operating cost
		pool.OperatingCost = 2_500_000_000

		// migrate height to custom keys
		pool.StartKey = strconv.FormatUint(pool.CurrentHeight, 10)

		// schedule upgrades for each runtime
		switch pool.Runtime {
		case "@kyve/evm":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "1.2.0",
				Binaries: "{\"macos\":\"https://cdn.discordapp.com/attachments/889827445132374036/989097802112041041/kyve-macos.zip?checksum=e2d99b4b6631e0f5350637bf51a5e3712a30f7aadedeb4f2975ebaa95880c861\"}",
			}
		case "@kyve/stacks":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/stacks/releases/download/v0.2.0/kyve-linux.zip?checksum=1c8b5ec983bccb70cbe435ded2edd7c2e65e9d12c7b032973f4adea6fa68281b\",\"macos\":\"https://github.com/kyve-org/stacks/releases/download/v0.2.0/kyve-macos.zip?checksum=acea13fd91281d545b500c8be5c02f92c0b296fee921516b072ed57b2607ca54\"}",
			}
		case "@kyve/bitcoin":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.2.0/kyve-linux.zip?checksum=0fbeaa64b22ab2ba34b6fa4571d2ff1fe7e573885abba48c15d878c19853bd79\",\"macos\":\"https://github.com/kyve-org/bitcoin/releases/download/v0.2.0/kyve-macos.zip?checksum=de901be709378f6bec81ed3c31d54fdfd67f26215e2614ffc14e80f6f98a29d1\"}",
			}
		case "@kyve/solana":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/solana/releases/download/v0.2.0/kyve-linux.zip?checksum=0e94e710624fb42947f538892160179050f8e99b286ae71668b7f5a48c285af0\",\"macos\":\"https://github.com/kyve-org/solana/releases/download/v0.2.0/kyve-macos.zip?checksum=bf25242be4a99c7b6181b835af7d90969bcdcd9a0473e850aa0ba065d723b038\"}",
			}
		case "@kyve/zilliqa":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.2.0/kyve-linux.zip?checksum=895b090bb2746a29fd4de168280b224353eb33bf7dd3f72fca8a60c250cfef2a\",\"macos\":\"https://github.com/kyve-org/zilliqa/releases/download/v0.2.0/kyve-macos.zip?checksum=5ca1330c922bfefacf1b94c8933930789cd314624aa29a22aaa55d2ab64e7d83\"}",
			}
		case "@kyve/near":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/near/releases/download/v0.2.0/kyve-linux.zip?checksum=7d310651e19aedfff9e3360b5e6108a6271314eedf52dd86e8945fd1bb1d3793\",\"macos\":\"https://github.com/kyve-org/near/releases/download/v0.2.0/kyve-macos.zip?checksum=34082976222fd7e8feac57eaf9af828508398e0688b7a3af9e0caf02777ab51d\"}",
			}
		case "@kyve/celo":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/celo/releases/download/v0.2.0/kyve-linux.zip?checksum=686c3bd436a3322f6bf09d2b3df465186360981aabd4da25e8576f0a3a867d66\",\"macos\":\"https://github.com/kyve-org/celo/releases/download/v0.2.0/kyve-macos.zip?checksum=1ddf112ea5e161108bfd3678fa7bb9a94be5cda8795ffb1cf8be904116518f05\"}",
			}
		case "@kyve/cosmos":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/cosmos/releases/download/v0.2.0/kyve-linux.zip?checksum=be96d9befd3a1084af6d9de2b49d81834edbf5888fb7dbdcbcaafbe22c1975a5\",\"macos\":\"https://github.com/kyve-org/cosmos/releases/download/v0.2.0/kyve-macos.zip?checksum=fcdd4c561096f7f7ccbf7eeb625c2317c5166b615fa208f92f73a2c704e1966a\"}",
			}
		case "@kyve/substrate":
			pool.UpgradePlan = &types.UpgradePlan{
				Version:  "0.2.0",
				Binaries: "{\"linux\":\"https://github.com/kyve-org/substrate/releases/download/v0.2.0/kyve-linux.zip?checksum=be96d9befd3a1084af6d9de2b49d81834edbf5888fb7dbdcbcaafbe22c1975a5\",\"macos\":\"https://github.com/kyve-org/substrate/releases/download/v0.2.0/kyve-macos.zip?checksum=fcdd4c561096f7f7ccbf7eeb625c2317c5166b615fa208f92f73a2c704e1966a\"}",
			}
		default:
			pool.UpgradePlan = &types.UpgradePlan{}
		}

		pool.UpgradePlan.ScheduledAt = uint64(ctx.BlockTime().Unix())
		pool.UpgradePlan.Duration = 300 // 30min

		// save changes
		registryKeeper.SetPool(ctx, pool)
	}
}

func migrateProposals(registryKeeper *registrykeeper.Keeper, ctx sdk.Context) {

	fmt.Printf("%sMigration Proposals to new key system\n", MigrationLoggerPrefix)

	for _, pool := range registryKeeper.GetAllPool(ctx) {
		proposalPrefixBuilder := types.KeyPrefixBuilder{Key: types.ProposalKeyPrefixIndex3}.AInt(pool.Id)
		store := prefix.NewStore(ctx.KVStore(registryKeeper.StoreKey()), proposalPrefixBuilder.Key)
		iterator := sdk.KVStorePrefixIterator(store, []byte{})

		defer iterator.Close()
		id := uint64(0)

		for ; iterator.Valid(); iterator.Next() {
			bundleId := string(iterator.Value())
			proposal, _ := registryKeeper.GetProposal(ctx, bundleId)
			proposal.Id = id
			id += 1
			proposal.Key = strconv.FormatUint(proposal.ToHeight-1, 10)
			registryKeeper.SetProposal(ctx, proposal)

			if id%100 == 0 {
				fmt.Printf("%sPool %d : Proposals processed: %d\n", MigrationLoggerPrefix, pool.Id, id)
			}

		}
		pool.TotalBundles = id + 1

		registryKeeper.SetPool(ctx, pool)
	}

	fmt.Printf("%sFinished proposal migration\n", MigrationLoggerPrefix)
}

func updateGovParams(ctx sdk.Context, govKeeper *govkeeper.Keeper) {
	govKeeper.SetDepositParams(ctx, govtypes.DepositParams{
		// 20,000 $KYVE
		MinDeposit: sdk.NewCoins(sdk.NewInt64Coin("tkyve", 20_000_000_000_000)),
		// 5 minutes
		MaxDepositPeriod: time.Minute * 5,
		// 100,000 $KYVE
		MinExpeditedDeposit: sdk.NewCoins(sdk.NewInt64Coin("tkyve", 100_000_000_000_000)),
	})

	govKeeper.SetVotingParams(ctx, govtypes.VotingParams{
		// 1 day
		VotingPeriod: time.Hour * 24,
		ProposalVotingPeriods: []govtypes.ProposalVotingPeriod{
			{
				ProposalType: "kyve.registry.v1beta1.CreatePoolProposal",
				// 2 hours
				VotingPeriod: time.Hour * 2,
			},
			{
				ProposalType: "kyve.registry.v1beta1.PausePoolProposal",
				// 5 minutes
				VotingPeriod: time.Minute * 5,
			},
		},
		// 30 minutes
		ExpeditedVotingPeriod: time.Minute * 30,
	})
}

func CreateUpgradeHandler(
	govKeeper *govkeeper.Keeper,
	registryKeeper *registrykeeper.Keeper,
	transferKeeper *ibctransferkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		createUnbondingParameters(registryKeeper, ctx)

		createProposalIndex(registryKeeper, ctx)

		migrateIBCDenoms(ctx, transferKeeper)

		updateGovParams(ctx, govKeeper)

		migratePools(registryKeeper, ctx)

		migrateProposals(registryKeeper, ctx)

		// Return.
		return vm, nil
	}
}

func equalTraces(dtA, dtB ibctransfertypes.DenomTrace) bool {
	return dtA.BaseDenom == dtB.BaseDenom && dtA.Path == dtB.Path
}
