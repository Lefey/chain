package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"kyve/x/registry/types"
)

// ClaimUploaderRole handles the logic of an SDK message that allows protocol nodes to claim the uploader role.
// Note that this function can only be called while the specified pool is in "genesis state".
// This function obeys "first come, first serve" mentality.
func (k msgServer) ClaimUploaderRole(
	goCtx context.Context, msg *types.MsgClaimUploaderRole,
) (*types.MsgClaimUploaderRoleResponse, error) {
	// Unwrap context and attempt to fetch the pool.
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.Id)

	// Error if the pool isn't found.
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), msg.Id)
	}

	// Error if the pool isn't in "genesis state".
	if pool.BundleProposal.NextUploader != "" {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrUploaderAlreadyClaimed.Error())
	}

	// Check if the sender is a protocol node (aka has staked into this pool).
	_, isStaker := k.GetStaker(ctx, msg.Creator, msg.Id)
	if !isStaker {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrNoStaker.Error())
	}

	// Check if enough nodes are online
	if len(pool.Stakers) < 2 {
		return nil, sdkErrors.Wrap(types.ErrNotEnoughNodesOnline, types.ErrNoStaker.Error())
	}

	// Check if pool is funded
	if pool.TotalFunds == 0 {
		return nil, sdkErrors.Wrap(types.ErrFundsTooLow, types.ErrNoStaker.Error())
	}

	// Error if the pool is paused.
	if pool.Paused {
		return nil, sdkErrors.Wrap(sdkErrors.ErrUnauthorized, types.ErrPoolPaused.Error())
	}

	// Update and return.
	pool.BundleProposal.NextUploader = msg.Creator
	pool.BundleProposal.CreatedAt = uint64(ctx.BlockTime().Unix())
	k.SetPool(ctx, pool)

	return &types.MsgClaimUploaderRoleResponse{}, nil
}
