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
	"github.com/stretchr/testify/assert"

	"autocctp.dev/utils"
	"autocctp.dev/utils/mocks"
)

func TestIncrementNumOfAccounts(t *testing.T) {
	// ARRANGE
	_, k, ctx := mocks.AutoCCTPKeeper(t)

	// ACT
	_, err := k.NumOfAccounts.Get(ctx, 0)
	assert.Error(t, err, "expected an error when the map is empty")
	k.IncrementNumOfAccounts(ctx, 0)
	count0, err := k.NumOfAccounts.Get(ctx, 0)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, uint64(1), count0, "expected 1 account for destination 0")

	// ACT
	k.IncrementNumOfAccounts(ctx, 1)
	k.IncrementNumOfAccounts(ctx, 1)
	count1, err := k.NumOfAccounts.Get(ctx, 1)
	assert.NoError(t, err)

	k.IncrementNumOfAccounts(ctx, 0)
	k.IncrementNumOfAccounts(ctx, 0)
	count0, err = k.NumOfAccounts.Get(ctx, 0)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, uint64(2), count1, "expected 2 accounts for destination 1")
	assert.Equal(t, uint64(3), count0, "expected 3 account for destination 0")
}

func TestIncrementNumOfTransfers(t *testing.T) {
	// ARRANGE
	_, k, ctx := mocks.AutoCCTPKeeper(t)

	// ACT
	_, err := k.NumOfTransfers.Get(ctx, 0)
	assert.Error(t, err, "expected an error when the map is empty")
	k.IncrementNumOfTransfers(ctx, 0)
	count0, err := k.NumOfTransfers.Get(ctx, 0)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, uint64(1), count0, "expected 1 transfer for destination 0")

	// ACT
	k.IncrementNumOfTransfers(ctx, 1)
	k.IncrementNumOfTransfers(ctx, 1)
	count1, err := k.NumOfTransfers.Get(ctx, 1)
	assert.NoError(t, err)

	k.IncrementNumOfTransfers(ctx, 0)
	k.IncrementNumOfTransfers(ctx, 0)
	count0, err = k.NumOfTransfers.Get(ctx, 0)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, uint64(3), count0, "expected 3 transfers for destination 0")
	assert.Equal(t, uint64(2), count1, "expected 2 transfers for destination 1")

	numPerDest, err := k.GetNumOfTransfersPerDestination(ctx)
	assert.NoError(t, err)
	assert.Len(t, numPerDest, 2)
	assert.Equal(t, uint64(3), numPerDest[0], "expected a different num of transfers for destination 0")
	assert.Equal(t, uint64(2), numPerDest[1], "expected a different num of transfers for destination 1")
}

func TestIncrementTotalTransferred(t *testing.T) {
	// ARRANGE
	_, k, ctx := mocks.AutoCCTPKeeper(t)

	// ACT
	_, err := k.TotalTransferred.Get(ctx, 0)
	assert.Error(t, err, "expected an error when the map is empty")
	k.IncrementTotalTransferred(ctx, 0, math.NewInt(1_000))
	count0, err := k.TotalTransferred.Get(ctx, 0)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, uint64(1_000), count0, "expected a different amount transferred for destination 0")

	// ACT
	k.IncrementTotalTransferred(ctx, 1, math.NewInt(1_000))
	k.IncrementTotalTransferred(ctx, 1, math.NewInt(1_000))
	count1, err := k.TotalTransferred.Get(ctx, 1)
	assert.NoError(t, err)

	k.IncrementTotalTransferred(ctx, 0, math.NewInt(1_000))
	k.IncrementTotalTransferred(ctx, 0, math.NewInt(1_000))
	count0, err = k.TotalTransferred.Get(ctx, 0)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, uint64(3_000), count0, "expected a different amount transferred for destination 0")
	assert.Equal(t, uint64(2_000), count1, "expected a different amount transferred for destination 1")

	amtPerDest, err := k.GetTotalTransferredPerDestination(ctx)
	assert.NoError(t, err)
	assert.Len(t, amtPerDest, 2)
	assert.Equal(t, uint64(3_000), amtPerDest[0], "expected a different amount transferred for destination 0")
	assert.Equal(t, uint64(2_000), amtPerDest[1], "expected a different amount transferred for destination 1")
}

func TestGetPendingTransfers(t *testing.T) {
	// ARRANGE
	_, k, ctx := mocks.AutoCCTPKeeper(t)

	// ACT: Get pending transfers when no transfers are present.
	acc, err := k.GetPendingTransfers(ctx)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, 0, len(acc), "expected no pending transfers to be returned")

	// ARRANGE
	_, err = utils.DummyPendingTransfersTest(ctx, k, 2, "", false)
	assert.NoError(t, err, "expected no error in the generation of dummy transfers")

	// ACT
	acc, err = k.GetPendingTransfers(ctx)
	assert.NoError(t, err)

	// ASSERT
	assert.Equal(t, 2, len(acc), "expected 2 pending transfers")
}
