// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, NASD Inc. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"

	"autocctp.dev/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Transaction commands for the %s module", types.ModuleName),
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(TxRegisterAccount())
	cmd.AddCommand(TxRegisterAccountSignerlessly())

	return cmd
}

func TxRegisterAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-account [destination-domain] [mint-recipient] [fallback-recipient] (destination-caller)",
		Short: "Register an AutoCCTP account for a destination domain, a mint recipient, and a fallback recipient",
		Long:  "Register an AutoCCTP account for a destination domain, a mint recipient, and a fallback recipient, with an optional destination caller",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args) != 4 {
				args = append(args, "")
			}
			accountProperties, err := ValidateAndParseAccountFields(args[0], args[1], args[2], args[3])
			if err != nil {
				return types.ErrInvalidInputs.Wrap(err.Error())
			}

			msg := &types.MsgRegisterAccount{
				Signer:            clientCtx.GetFromAddress().String(),
				DestinationDomain: accountProperties.DestinationDomain,
				MintRecipient:     accountProperties.MintRecipient,
				FallbackRecipient: accountProperties.FallbackRecipient,
				DestinationCaller: accountProperties.DestinationCaller,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func TxRegisterAccountSignerlessly() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-account-signerlessly [destination-domain] [mint-recipient] (destination-caller)",
		Short: "Signerlessly register an AutoCCTP account for a destination domain, a mint recipient, and a fallback recipient",
		Long: `Signerlessly register an AutoCCTP account for a destination domain, a mint recipient, and a fallback recipient, with an optional destination caller.
		A signerless registration does not require an existing wallet because no signature is required.`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args) != 4 {
				args = append(args, "")
			}
			accountProperties, err := ValidateAndParseAccountFields(args[0], args[1], args[2], args[3])
			if err != nil {
				return types.ErrInvalidInputs.Wrap(err.Error())
			}

			address := types.GenerateAddress(*accountProperties)

			msg := &types.MsgRegisterAccountSignerlessly{
				Signer:            address.String(),
				DestinationDomain: accountProperties.DestinationDomain,
				MintRecipient:     accountProperties.MintRecipient,
				FallbackRecipient: accountProperties.FallbackRecipient,
				DestinationCaller: accountProperties.DestinationCaller,
			}

			factory, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			builder, err := factory.BuildUnsignedTx(msg)
			if err != nil {
				return err
			}

			// Create an empty signature with the custom PubKey to allow non existent
			// account to send `MsgRegisterAccountSignerlessly` messages.
			err = builder.SetSignatures(signingtypes.SignatureV2{
				PubKey: &types.PubKey{Key: address},
				Data: &signingtypes.SingleSignatureData{
					SignMode:  signingtypes.SignMode_SIGN_MODE_DIRECT,
					Signature: []byte(""),
				},
			})
			if err != nil {
				return nil
			}

			if clientCtx.GenerateOnly {
				bz, err := clientCtx.TxConfig.TxJSONEncoder()(builder.GetTx())
				if err != nil {
					return err
				}

				return clientCtx.PrintString(fmt.Sprintf("%s\n", bz))
			}

			bz, err := clientCtx.TxConfig.TxEncoder()(builder.GetTx())
			if err != nil {
				return err
			}
			res, err := clientCtx.BroadcastTx(bz)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
