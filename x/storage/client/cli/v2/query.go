package v2

import (
	"fmt"

	"github.com/bnb-chain/greenfield/x/storage/client/cli"
	"github.com/bnb-chain/greenfield/x/storage/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group storage queries under a subcommand
	storageQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	storageQueryCmd.AddCommand(
		cli.CmdQueryParams(),
		cli.CmdHeadBucket(),
		cli.CmdHeadObject(),
		cli.CmdListBuckets(),
		cli.CmdListObjects(),
		cli.CmdVerifyPermission(),
		cli.CmdHeadGroup(),
		cli.CmdListGroup(),
		// request body did not change, so we can reuse the same command
		cli.CmdHeadGroupMember())

	return storageQueryCmd
}
