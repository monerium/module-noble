package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
	"github.com/spf13/cobra"
)

func GetBlacklistQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "blacklist",
		Short:                      "Querying commands for the blacklist submodule",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(QueryBlacklistOwner())

	return cmd
}

func QueryBlacklistOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner",
		Short: "Query the submodule's owner",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := blacklist.NewQueryClient(clientCtx)

			res, err := queryClient.Owner(context.Background(), &blacklist.QueryOwner{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
