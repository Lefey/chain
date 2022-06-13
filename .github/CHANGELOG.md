<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## Unreleased

### Features

- Unbonding time for unstaking from a pool. Protocol node runners have to keep their node running during the unbonding.
- Unbonding time for undelegating from a staker in a pool. The unbonding is performed immediately but the delegator has
  to wait until the tokens are transferred back.

### Improvements

- Bump [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) to [`v0.45.5`](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.45.5). See [CHANGELOG](https://github.com/cosmos/cosmos-sdk/blob/v0.45.5/CHANGELOG.md#v0455---2022-06-09) for more details.
- Bump [IBC](https://github.com/cosmos/ibc-go) to [`v3.0.0`](https://github.com/cosmos/ibc-go/releases/tag/v3.0.0). See [CHANGELOG](https://github.com/cosmos/ibc-go/blob/v3.0.0/CHANGELOG.md#v300---2022-03-15) for more details.

### Client Breaking Changes

- Switch vote type in `MsgVoteProposal` from `uint64` to `enum`.

## [v0.4.0](https://github.com/KYVENetwork/chain/releases/tag/v0.4.0) - 2022-06-7

### Features

- Implemented scheduled upgrades for pool versions
- Implemented `abstain` vote besides `valid` and `invalid`. Validators who don't vote 5 times in a row at all get removed with a timeout slash

### Client Breaking Changes

- The arg `vote` on `MsgVoteProposal` changed from `bool` to `uint64`. 0 = valid, 1 = invalid, 2 = abstain
- The arg `versions` on `MsgCreatePoolProposal` changed to `version`
- The arg `binaries` got added to `MsgCreatePoolProposal`

### Improvements

- Check the quorum of the bundle proposal on chain to prevent unjustified slashes
- Don't drop bundle proposals if one funder can't afford the funding cost, instead remove all of them and proceed
- If a validator submits a `NO_DATA_BUNDLE` the will just skip the upload instead of proposing an empty bundle
- Added query `QueryFunder`
- Added query `QueryStaker`
- Added query `QueryDelegator`

### Bug Fixes

### Deprecated

- Deprecated `versions` on `kyve.registry.v1beta1.Pool`
