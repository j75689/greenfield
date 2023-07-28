package v2

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bnb-chain/greenfield/x/storage/client/cli"
	"github.com/bnb-chain/greenfield/x/storage/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		cli.CmdCreateBucket(),
		cli.CmdDeleteBucket(),
		cli.CmdUpdateBucketInfo(),
		cli.CmdMirrorBucket(),
		cli.CmdDiscontinueBucket(),
	)

	cmd.AddCommand(
		cli.CmdCreateObject(),
		cli.CmdDeleteObject(),
		cli.CmdCancelCreateObject(),
		cli.CmdCopyObject(),
		cli.CmdMirrorObject(),
		cli.CmdDiscontinueObject(),
		cli.CmdUpdateObjectInfo(),
	)

	cmd.AddCommand(
		CmdCreateGroup(),
		cli.CmdDeleteGroup(),
		CmdUpdateGroupMember(),
		cli.CmdUpdateGroupExtra(),
		cli.CmdLeaveGroup(),
		cli.CmdMirrorGroup(),
	)

	cmd.AddCommand(
		cli.CmdPutPolicy(),
		cli.CmdDeletePolicy(),
	)

	cmd.AddCommand(cli.CmdCancelMigrateBucket())
	return cmd
}

func CmdCreateGroup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-group [group-name]",
		Short: "Create a new group with optional members, split member addresses by ','",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argGroupName := args[0]
			argMemberList, _ := cmd.Flags().GetString(cli.FlagMemberList)
			extra, _ := cmd.Flags().GetString(cli.FlagExtra)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var memberAddrs []sdk.AccAddress
			if argMemberList != "" {
				members := strings.Split(argMemberList, ",")
				for _, member := range members {
					memberAddr, err := sdk.AccAddressFromHexUnsafe(member)
					if err != nil {
						return err
					}
					memberAddrs = append(memberAddrs, memberAddr)
				}
			}
			msg := types.NewMsgCreateGroup(
				clientCtx.GetFromAddress(),
				argGroupName,
				memberAddrs,
				extra,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagExtra, "", "extra info for the group")
	cmd.Flags().String(cli.FlagMemberList, "", "init members of the group")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateGroupMember() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-group-member [group-name] [member-to-add] [member-to-delete]",
		Short: "Update the member of the group you own, split member addresses by ,",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argGroupName := args[0]
			argMemberToAdd := args[1]
			argMemberToDelete := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			var memberAddrsToAdd []sdk.AccAddress
			membersToAdd := strings.Split(argMemberToAdd, ",")
			for _, member := range membersToAdd {
				memberAddr, err := sdk.AccAddressFromHexUnsafe(member)
				if err != nil {
					return err
				}
				memberAddrsToAdd = append(memberAddrsToAdd, memberAddr)
			}
			var memberAddrsToDelete []sdk.AccAddress
			membersToDelete := strings.Split(argMemberToDelete, ",")
			for _, member := range membersToDelete {
				memberAddr, err := sdk.AccAddressFromHexUnsafe(member)
				if err != nil {
					return err
				}
				memberAddrsToDelete = append(memberAddrsToDelete, memberAddr)
			}
			msg := types.NewMsgUpdateGroupMember(
				clientCtx.GetFromAddress(),
				clientCtx.GetFromAddress(),
				argGroupName,
				memberAddrsToAdd,
				memberAddrsToDelete,
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
