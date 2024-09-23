// Copyright 2024 Monerium ehf.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/monerium/module-noble/v2/types/blacklist"
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
	cmd.AddCommand(QueryBlacklistAdmins())
	cmd.AddCommand(QueryAdversaries())

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

func QueryBlacklistAdmins() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admins",
		Short: "Query the submodule's admin accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := blacklist.NewQueryClient(clientCtx)

			res, err := queryClient.Admins(context.Background(), &blacklist.QueryAdmins{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryAdversaries() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "adversaries",
		Short: "Query the banned adversary accounts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := blacklist.NewQueryClient(clientCtx)

			res, err := queryClient.Adversaries(context.Background(), &blacklist.QueryAdversaries{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
