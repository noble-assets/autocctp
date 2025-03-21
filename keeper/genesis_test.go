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

package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"autocctp.dev/testutil/mocks"
)

func TestExportGenesis(t *testing.T) {
	_, k, ctx := mocks.AutoCCTPKeeper(t)

	// Add num of accounts
	err := k.IncrementNumOfAccounts(ctx, 0)
	require.NoError(t, err)
	err = k.IncrementNumOfAccounts(ctx, 0)
	require.NoError(t, err)
	err = k.IncrementNumOfAccounts(ctx, 0)
	require.NoError(t, err)
	err = k.IncrementNumOfAccounts(ctx, 1)
	require.NoError(t, err)
	err = k.IncrementNumOfAccounts(ctx, 1)
	require.NoError(t, err)
	err = k.IncrementNumOfAccounts(ctx, 2)
	require.NoError(t, err)

	// Add num of transfers

	err = k.IncrementNumOfTransfers(ctx, 0)
	require.NoError(t, err)
	err = k.IncrementNumOfTransfers(ctx, 0)
	require.NoError(t, err)
	err = k.IncrementNumOfTransfers(ctx, 0)
	require.NoError(t, err)
	err = k.IncrementNumOfTransfers(ctx, 1)
	require.NoError(t, err)
	err = k.IncrementNumOfTransfers(ctx, 1)
	require.NoError(t, err)
	err = k.IncrementNumOfTransfers(ctx, 2)
	require.NoError(t, err)

	// Add total transferred
	err = k.IncrementTotalTransferred(ctx, 0, math.NewInt(1_000))
	require.NoError(t, err)
	err = k.IncrementTotalTransferred(ctx, 0, math.NewInt(1_000))
	require.NoError(t, err)
	err = k.IncrementTotalTransferred(ctx, 0, math.NewInt(1_000))
	require.NoError(t, err)
	err = k.IncrementTotalTransferred(ctx, 1, math.NewInt(1_000))
	require.NoError(t, err)
	err = k.IncrementTotalTransferred(ctx, 1, math.NewInt(1_000))
	require.NoError(t, err)
	err = k.IncrementTotalTransferred(ctx, 2, math.NewInt(1_000))
	require.NoError(t, err)

	genesis := k.ExportGenesis(ctx)
	require.Len(t, genesis.NumOfAccounts, 3, "expected 3 destination domain for the accounts")
	require.Len(t, genesis.NumOfTransfers, 3, "expected 3 destination domain for the num of transfers")
	require.Len(t, genesis.TotalTransferred, 3, "expected 3 destination domain for the total transferred")

	require.Equal(t, uint64(3), genesis.NumOfAccounts[0])
	require.Equal(t, uint64(3), genesis.NumOfTransfers[0])
	require.Equal(t, uint64(3_000), genesis.TotalTransferred[0])

	require.Equal(t, uint64(2), genesis.NumOfAccounts[1])
	require.Equal(t, uint64(2), genesis.NumOfTransfers[1])
	require.Equal(t, uint64(2_000), genesis.TotalTransferred[1])

	require.Equal(t, uint64(1), genesis.NumOfAccounts[2])
	require.Equal(t, uint64(1), genesis.NumOfTransfers[2])
	require.Equal(t, uint64(1_000), genesis.TotalTransferred[2])
}
