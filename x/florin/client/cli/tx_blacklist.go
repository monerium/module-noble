package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/noble-assets/florin/x/florin/types/blacklist"
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
