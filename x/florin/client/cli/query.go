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
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gogo/protobuf/proto"
	"github.com/monerium/module-noble/v2/x/florin/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetBlacklistQueryCmd())

	cmd.AddCommand(QueryAllowedDenoms())
	cmd.AddCommand(QueryOwners())
	cmd.AddCommand(QuerySystems())
	cmd.AddCommand(QueryAdmins())
	cmd.AddCommand(QueryMaxMintAllowances())
	cmd.AddCommand(QueryMintAllowances())
	cmd.AddCommand(QueryMintAllowance())

	return cmd
}

func QueryAllowedDenoms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowed-denoms",
		Short: "Query the allowed denoms of this module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AllowedDenoms(context.Background(), &types.QueryAllowedDenoms{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryOwners() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "owner (denom)",
		Short: "Query the owner of a specific or all denoms",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				res, err = queryClient.Owner(context.Background(), &types.QueryOwner{Denom: args[0]})
			} else {
				res, err = queryClient.Owners(context.Background(), &types.QueryOwners{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QuerySystems() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "systems (denom)",
		Short: "Query the system accounts of a specific or all denoms",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				res, err = queryClient.SystemsByDenom(context.Background(), &types.QuerySystemsByDenom{Denom: args[0]})
			} else {
				res, err = queryClient.Systems(context.Background(), &types.QuerySystems{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryAdmins() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "admins (denom)",
		Short: "Query the admin accounts of a specific or all denoms",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				res, err = queryClient.AdminsByDenom(context.Background(), &types.QueryAdminsByDenom{Denom: args[0]})
			} else {
				res, err = queryClient.Admins(context.Background(), &types.QueryAdmins{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryMaxMintAllowances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "max-mint-allowance (denom)",
		Short: "Query the max mint allowance of a specific or all denoms",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error

			if len(args) == 1 {
				res, err = queryClient.MaxMintAllowance(context.Background(), &types.QueryMaxMintAllowance{Denom: args[0]})
			} else {
				res, err = queryClient.MaxMintAllowances(context.Background(), &types.QueryMaxMintAllowances{})
			}

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryMintAllowances() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-allowances [denom]",
		Short: "Query the mint allowances of a specific denom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MintAllowances(context.Background(), &types.QueryMintAllowances{
				Denom: args[0],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryMintAllowance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint-allowance [denom] [account]",
		Short: "Query the mint allowance of a specific system account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.MintAllowance(context.Background(), &types.QueryMintAllowance{
				Denom:   args[0],
				Account: args[1],
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
