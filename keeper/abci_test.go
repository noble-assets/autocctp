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
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"autocctp.dev/utils"
	"autocctp.dev/utils/mocks"
)

func TestExecuteTransfer(t *testing.T) {
	expCall := "expected %s call(s) to the DepositForBurn endpoint"
	expCallWithCaller := "expected %s call(s) to the DepositForBurnWithCaller endpoint"

	m, k, ctx := mocks.AutoCCTPKeeper(t)
	mc := m.CCTPServer.MockCounter
	bk := m.BankKeeper
	cctps := m.CCTPServer

	// ACT: Execute transfer without pending transfers.
	k.ExecuteTransfers(ctx)

	// ASSERT: No mock state change because no pending transfers.
	assert.Equal(t, 0, mc.NumDepositForBurn)
	assert.Equal(t, 0, mc.NumDepositForBurnWithCaller)

	// ARRANGE: One pending transfer is registered but the account has not funds.
	addresses, err := utils.DummyPendingTransfersTest(ctx, k, 1, "", false)
	assert.NoError(t, err)

	// ACT
	k.ExecuteTransfers(ctx)

	// ASSERT
	assert.Equal(t, 0, mc.NumDepositForBurn, fmt.Sprintf(expCall, "no"))
	assert.Equal(t, 0, mc.NumDepositForBurnWithCaller, fmt.Sprintf(expCallWithCaller, "no"))

	// ARRANGE: Add funds to the accounts in the pending transfers but not of the minting denom.
	bk.Balances[addresses[0]] = sdk.Coins{sdk.NewInt64Coin("ubtc", 10)}

	// ACT
	k.ExecuteTransfers(ctx)

	// ASSERT
	assert.Equal(t, 0, mc.NumDepositForBurn, fmt.Sprintf(expCall, "no"))
	assert.Equal(t, 0, mc.NumDepositForBurnWithCaller, fmt.Sprintf(expCallWithCaller, "no"))

	// ARRANGE: Add correct funds to the account.
	bk.Balances[addresses[0]] = sdk.Coins{sdk.NewInt64Coin("uusdc", 10)}

	// ACT
	k.ExecuteTransfers(ctx)

	// ASSERT: CCTP server is called properly.
	assert.Equal(t, 1, mc.NumDepositForBurn, fmt.Sprintf(expCall, "one"))
	assert.Equal(t, 0, mc.NumDepositForBurnWithCaller, fmt.Sprintf(expCallWithCaller, "no"))

	// ARRANGE: Create a account with destination caller and multiple tokens in the balance.
	// Similar to previous test but the account has multiple tokens and the CCTP method called is different.
	mocks.ResetTest(t, ctx, k, m)

	addresses, err = utils.DummyPendingTransfersTest(ctx, k, 1, "0", true)
	assert.NoError(t, err)
	bk.Balances[addresses[0]] = sdk.NewCoins(
		sdk.NewInt64Coin("uusdc", 1_000_000),
		sdk.NewInt64Coin("ueurc", 10),
	)

	// ACT
	k.ExecuteTransfers(ctx)

	// ASSERT: CCTP server is called properly, with one call per tokens for the account with caller.
	assert.Equal(t, 0, mc.NumDepositForBurn, fmt.Sprintf(expCall, "no"))
	assert.Equal(t, 1, mc.NumDepositForBurnWithCaller, fmt.Sprintf(expCallWithCaller, "one"))
	numOfTransfers, err := k.NumOfTransfers.Get(ctx, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), numOfTransfers, "expected a different amount of transfers")
	totalTransferred, err := k.TotalTransferred.Get(ctx, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1_000_000), totalTransferred, "expected a different total transferred")

	// ARRANGE: Trigger an error in the CCTP server execution.
	mocks.ResetTest(t, ctx, k, m)
	cctps.Failing = true

	addresses, err = utils.DummyPendingTransfersTest(ctx, k, 1, "0", false)
	assert.NoError(t, err)
	addressesWithCaller, err := utils.DummyPendingTransfersTest(ctx, k, 1, "0", true)
	assert.NoError(t, err)

	bk.Balances[addresses[0]] = sdk.NewCoins(
		sdk.NewInt64Coin("uusdc", 1_000_000),
	)
	bk.Balances[addressesWithCaller[0]] = sdk.NewCoins(
		sdk.NewInt64Coin("uusdc", 1_000_000),
		sdk.NewInt64Coin("ueurc", 10),
	)

	// ACT
	k.ExecuteTransfers(ctx)

	// ASSERT: calls to endpoints remain unaltered
	assert.Equal(t, 0, mc.NumDepositForBurn, fmt.Sprintf(expCall, "no"))
	assert.Equal(t, 0, mc.NumDepositForBurnWithCaller, fmt.Sprintf(expCallWithCaller, "no"))
	_, err = k.NumOfTransfers.Get(ctx, 0)
	assert.Error(t, err, "expected no transfers when cctp returns an error")
	_, err = k.TotalTransferred.Get(ctx, 0)
	assert.Error(t, err, "expected zero total transferred when cctp returns an error")
}
