package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"kyve/x/registry/types"
)

var _ = strconv.Itoa(0)

func CmdVoteProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-proposal [id] [bundle-id] [support]",
		Short: "Broadcast message vote-proposal",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argId, err := cast.ToUint64E(args[0])
			if err != nil {
				return err
			}
			argBundleId := args[1]
			argSupport, err := cast.ToBoolE(args[2])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgVoteProposal(
				clientCtx.GetFromAddress().String(),
				argId,
				argBundleId,
				argSupport,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
