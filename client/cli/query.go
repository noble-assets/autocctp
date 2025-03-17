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
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/gogoproto/proto"
	"github.com/spf13/cobra"

	"autocctp.dev/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Query commands for the %s module", types.ModuleName),
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(QueryStats())
	cmd.AddCommand(QueryAddress())

	return cmd
}

func QueryAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "address [destination-domain] [mint-recipient] (destination-caller)",
		Short: "Query AutoCCTP address by destination domain, a mint recipient, and a fallback recipient",
		Long: `Query AutoCCTP address by destination domain, a mint recipient, and a fallback recipient, with an optional destination caller address.
		The command creates a 32-bytes representation of the mint recipient, and optionally of the destination caller, to 
		get the cross-chain representation of the address.
`,
		Args: cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			if len(args) != 4 {
				args = append(args, "")
			}

			accountProperties, err := ValidateAndParseAccountFields(args[0], args[1], args[2], args[3])
			if err != nil {
				return types.ErrInvalidInputs.Wrap(err.Error())
			}

			res, err := queryClient.Address(context.Background(), &types.QueryAddress{
				DestinationDomain: accountProperties.DestinationDomain,
				MintRecipient:     accountProperties.MintRecipient,
				FallbackRecipient: accountProperties.FallbackRecipient,
				DestinationCaller: accountProperties.DestinationCaller,
			})
			if err != nil {
				return fmt.Errorf("error executing the query: %w", err)
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryStats() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats (destination-domain)",
		Short: "Query AutoCCTP stats",
		Long:  "Query AutoCCTP usage statistics for all or a specific destination domain",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var res proto.Message
			var err error
			if len(args) == 1 {
				destinationDomain, valError := ValidateDestinationDomain(args[0])
				if valError != nil {
					return types.ErrInvalidInputs.Wrap(valError.Error())
				}

				res, err = queryClient.StatsByDestinationDomain(context.Background(), &types.QueryStatsByDestinationDomain{
					DestinationDomain: uint32(destinationDomain),
				})
			} else {
				res, err = queryClient.Stats(context.Background(), &types.QueryStats{})
			}
			if err != nil {
				return fmt.Errorf("error executing the query: %w", err)
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
