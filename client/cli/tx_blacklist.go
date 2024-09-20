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
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/monerium/module-noble/v2/types/blacklist"
	"github.com/spf13/cobra"
)

func GetBlacklistTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "blacklist",
		Short:                      "Transactions commands for the blacklist submodule",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(TxBlacklistAcceptOwnership())
	cmd.AddCommand(TxBlacklistAddAdminAccount())
	cmd.AddCommand(TxBan())
	cmd.AddCommand(TxBlacklistRemoveAdminAccount())
	cmd.AddCommand(TxBlacklistTransferOwnership())
	cmd.AddCommand(TxUnban())

	return cmd
}

func TxBlacklistAcceptOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-ownership",
		Short: "Accept ownership of submodule",
		Long:  "Accept ownership of submodule, assuming there is an pending ownership transfer",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blacklist.MsgAcceptOwnership{
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxBlacklistAddAdminAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-admin-account [account]",
		Short: "Adds an admin account to the submodule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blacklist.MsgAddAdminAccount{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxBan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ban",
		Short: "Bans a specific adversary account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blacklist.MsgBan{
				Signer:    clientCtx.GetFromAddress().String(),
				Adversary: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxBlacklistRemoveAdminAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-admin-account [account]",
		Short: "Removes an admin account from the submodule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blacklist.MsgRemoveAdminAccount{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxBlacklistTransferOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership [new-owner]",
		Short: "Transfer ownership of submodule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blacklist.MsgTransferOwnership{
				Signer:   clientCtx.GetFromAddress().String(),
				NewOwner: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxUnban() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unban",
		Short: "Unbans a specific friend account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &blacklist.MsgUnban{
				Signer: clientCtx.GetFromAddress().String(),
				Friend: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
