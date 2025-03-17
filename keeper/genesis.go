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

	"autocctp.dev/types"
)

func (k *Keeper) InitGenesis(ctx context.Context, genesis types.GenesisState) {
	for key, value := range genesis.NumOfAccounts {
		if err := k.NumOfAccounts.Set(ctx, key, value); err != nil {
			panic(err)
		}
	}
	for key, value := range genesis.NumOfTransfers {
		if err := k.NumOfTransfers.Set(ctx, key, value); err != nil {
			panic(err)
		}
	}
	for key, value := range genesis.TotalTransferred {
		if err := k.TotalTransferred.Set(ctx, key, value); err != nil {
			panic(err)
		}
	}
}

func (k *Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	// NOTE: errors are intentionally not handled to unconditionally allow genesis export.
	numOfAccount, _ := k.GetNumOfAccountPerDestination(ctx)
	numOfTransfers, _ := k.GetNumOfTransfersPerDestination(ctx)
	totTransferred, _ := k.GetTotalTransferredPerDestination(ctx)

	return &types.GenesisState{
		NumOfAccounts:    numOfAccount,
		NumOfTransfers:   numOfTransfers,
		TotalTransferred: totTransferred,
	}
}
