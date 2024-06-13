package cli

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/noble-assets/florin/x/florin/types"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transactions commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetBlacklistTxCmd())

	cmd.AddCommand(TxAcceptOwnership())
	cmd.AddCommand(TxAddAdminAccount())
	cmd.AddCommand(TxAddSystemAccount())
	cmd.AddCommand(TxBurn())
	cmd.AddCommand(TxMint())
	cmd.AddCommand(TxRemoveAdminAccount())
	cmd.AddCommand(TxRemoveSystemAccount())
	cmd.AddCommand(TxSetMaxMintAllowance())
	cmd.AddCommand(TxSetMintAllowance())
	cmd.AddCommand(TxTransferOwnership())

	return cmd
}

func TxAcceptOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept-ownership",
		Short: "Accept ownership of module",
		Long:  "Accept ownership of module, assuming there is an pending ownership transfer",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgAcceptOwnership{
				Signer: clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAddAdminAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-admin-account [account]",
		Short: "Adds an admin account to the module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgAddAdminAccount{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxAddSystemAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-system-account [account]",
		Short: "Adds a system account to the module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgAddSystemAccount{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxBurn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn [from] [amount]",
		Short: "Transaction that burns tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := &types.MsgBurn{
				Signer: clientCtx.GetFromAddress().String(),
				From:   args[0],
				Amount: amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxMint() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [to] [amount]",
		Short: "Transaction that mints tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := &types.MsgMint{
				Signer: clientCtx.GetFromAddress().String(),
				To:     args[0],
				Amount: amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRemoveAdminAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-admin-account [account]",
		Short: "Removes an admin account from the module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRemoveAdminAccount{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRemoveSystemAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove-system-account [account]",
		Short: "Removes a system account from the module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgRemoveSystemAccount{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxSetMaxMintAllowance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-max-mint-allowance [amount]",
		Short: "Sets the max mint allowance a minter can have",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[0])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := &types.MsgSetMaxMintAllowance{
				Signer: clientCtx.GetFromAddress().String(),
				Amount: amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxSetMintAllowance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-mint-allowance [account] [amount]",
		Short: "Sets the mint allowance of a minter",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdk.NewIntFromString(args[1])
			if !ok {
				return errors.New("invalid amount")
			}

			msg := &types.MsgSetMintAllowance{
				Signer:  clientCtx.GetFromAddress().String(),
				Account: args[0],
				Amount:  amount,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxTransferOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer-ownership [new-owner]",
		Short: "Transfer ownership of module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgTransferOwnership{
				Signer:   clientCtx.GetFromAddress().String(),
				NewOwner: args[0],
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
