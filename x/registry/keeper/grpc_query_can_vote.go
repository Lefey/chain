package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kyve/x/registry/types"
)

func (k Keeper) CanVote(goCtx context.Context, req *types.QueryCanVoteRequest) (*types.QueryCanVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Load pool
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), req.PoolId)
	}

	// Check if pool is paused
	if pool.Paused {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Pool is paused",
		}, nil
	}

	// Check if sender is a staker in pool
	_, isStaker := k.GetStaker(ctx, req.Voter, req.PoolId)
	if !isStaker {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Voter is no staker",
		}, nil
	}

	// Check if empty bundle
	if pool.BundleProposal.BundleId == "" {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Can not vote on empty bundle",
		}, nil
	}

	// Check if tx matches current bundleProposal
	if req.BundleId != pool.BundleProposal.BundleId {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Provided bundleId does not match current one",
		}, nil
	}

	// Check if sender has not voted yet
	hasVotedValid, hasVotedInvalid := false, false

	for _, voter := range pool.BundleProposal.VotersValid {
		if voter == req.Voter {
			hasVotedValid = true
		}
	}

	for _, voter := range pool.BundleProposal.VotersInvalid {
		if voter == req.Voter {
			hasVotedInvalid = true
		}
	}

	if hasVotedValid || hasVotedInvalid {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Voter already voted",
		}, nil
	}

	// check if voter is not uploader
	if pool.BundleProposal.GetUploader() == req.Voter {
		return &types.QueryCanVoteResponse{
			Possible: false,
			Reason:   "Voter is uploader",
		}, nil
	}

	return &types.QueryCanVoteResponse{
		Possible: true,
		Reason:   "",
	}, nil
}
