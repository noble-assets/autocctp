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

package keeper

import (
	"context"
	"fmt"

	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
)

// ExecuteTransfers is an end block hook that clears all pending transfers from the transient state.
func (k *Keeper) ExecuteTransfers(ctx context.Context) {
	transfers, err := k.GetPendingTransfers(ctx)
	if len(transfers) == 0 || err != nil {
		return
	}

	k.logger.Info(fmt.Sprintf("executing %d automatic cctp transfer(s)", len(transfers)))

	mintingToken := k.ftfKeeper.GetMintingDenom(ctx)
	for _, transfer := range transfers {
		balance := k.bankKeeper.GetBalance(ctx, transfer.GetAddress(), mintingToken.Denom)
		if balance.IsZero() {
			continue
		}

		var err error
		if len(transfer.DestinationCaller) == 0 {
			_, err = k.cctpServer.DepositForBurn(ctx, &cctptypes.MsgDepositForBurn{
				From:              transfer.Address,
				Amount:            balance.Amount,
				DestinationDomain: transfer.DestinationDomain,
				MintRecipient:     transfer.MintRecipient,
				BurnToken:         balance.Denom,
			})
		} else {
			_, err = k.cctpServer.DepositForBurnWithCaller(ctx, &cctptypes.MsgDepositForBurnWithCaller{
				From:              transfer.Address,
				Amount:            balance.Amount,
				DestinationDomain: transfer.DestinationDomain,
				MintRecipient:     transfer.MintRecipient,
				BurnToken:         balance.Denom,
				DestinationCaller: transfer.DestinationCaller,
			})
		}

		if err != nil {
			k.logger.Error(
				"unable to execute automatic cctp transfer",
				"from", transfer.Address,
				"to", transfer.MintRecipient,
				"denom", balance.Denom,
				"destination_domain", transfer.DestinationDomain,
				"amount", balance.Amount,
				"err", err,
			)
		} else {
			if err := k.IncrementNumOfTransfers(ctx, transfer.DestinationDomain); err != nil {
				k.logger.Error("end block", "error", err)
			}
			if err := k.IncrementTotalTransferred(ctx, transfer.DestinationDomain, balance.Amount); err != nil {
				k.logger.Error("end block", "error", err)
			}
		}
	}
}
