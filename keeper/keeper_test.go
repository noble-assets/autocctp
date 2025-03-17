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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"autocctp.dev/types"
	"autocctp.dev/utils"
	"autocctp.dev/utils/mocks"
)

func TestValidateAccountProperties(t *testing.T) {
	_, k, _ := mocks.AutoCCTPKeeper(t)
	validProperties := utils.ValidPropertiesTest(false)
	validPropertiesWithCaller := utils.ValidPropertiesTest(true)

	invalidFalbackRecipient := "cosmos1y5azhw4a99s4tm4kwzfwus52tjlvsaywuq3q3m"

	testCases := []struct {
		name        string
		setup       func(*types.AccountProperties)
		withCaller  bool
		errContains string
	}{
		{
			name: "fail when the mint recipient is not valid",
			setup: func(ap *types.AccountProperties) {
				ap.MintRecipient = []byte{}
			},
			withCaller:  false,
			errContains: types.ErrInvalidMintRecipient.Error(),
		},
		{
			name: "fail when the fallback recipient is empty",
			setup: func(ap *types.AccountProperties) {
				ap.FallbackRecipient = ""
			},
			withCaller:  false,
			errContains: types.ErrInvalidFallbackRecipient.Error(),
		},
		{
			name: "fail when the fallback recipient is not chain address",
			setup: func(ap *types.AccountProperties) {
				ap.FallbackRecipient = invalidFalbackRecipient
			},
			withCaller:  false,
			errContains: types.ErrInvalidFallbackRecipient.Error(),
		},
		{
			name: "fail when the destination caller is not valid",
			setup: func(ap *types.AccountProperties) {
				ap.DestinationCaller = []byte("0")
			},
			withCaller:  true,
			errContains: types.ErrInvalidDestinationCaller.Error(),
		},
		{
			name:        "success with valid properties no destination caller",
			setup:       func(ap *types.AccountProperties) {},
			withCaller:  false,
			errContains: "",
		},
		{
			name:        "success with valid properties with destination caller",
			setup:       func(ap *types.AccountProperties) {},
			withCaller:  true,
			errContains: "",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			properties := validProperties
			if tC.withCaller {
				properties = validPropertiesWithCaller
			}
			tC.setup(&properties)

			err := k.ValidateAccountProperties(properties)
			if tC.errContains == "" {
				require.NoError(t, err, "expected no error executing the validation")
			} else {
				require.Error(t, err, "expected an error executing the validation")
				require.ErrorContains(t, err, tC.errContains, "expected a different error")
			}
		})
	}
}

func TestSendRestrictionFn(t *testing.T) {
	// ARRANGE
	mocks, k, ctx := mocks.AutoCCTPKeeper(t)
	ak := mocks.AccountKeeper

	acc := utils.DummyAccountTest(false)

	// ACT
	toAddr, err := k.SendRestrictionFn(ctx, sdk.AccAddress{}, acc.GetAddress(), sdk.Coins{})

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, acc.GetAddress(), toAddr, "expected the returned address unaltered")

	_, err = k.PendingTransfers.Get(ctx, acc.GetAddress().String())
	assert.Error(t, err, "expected no registered pending transfers when receiver is not a stored account")

	// ARRANGE
	ak.Accounts[acc.GetAddress().String()] = acc.BaseAccount

	// ACT
	toAddr, err = k.SendRestrictionFn(ctx, sdk.AccAddress{}, acc.GetAddress(), sdk.Coins{})

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, acc.GetAddress(), toAddr, "expected the returned address unaltered")

	_, err = k.PendingTransfers.Get(ctx, acc.GetAddress().String())
	assert.Error(t, err, "expected no registered pending transfers when receiver is not autocctp account")

	// ARRANGE
	ak.Accounts[acc.GetAddress().String()] = &acc

	// ACT: Call function with an nil account
	toAddr, err = k.SendRestrictionFn(ctx, sdk.AccAddress{}, acc.GetAddress(), sdk.Coins{})

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, acc.GetAddress(), toAddr, "expected the returned address unaltered")

	_, err = k.PendingTransfers.Get(ctx, acc.GetAddress().String())
	assert.Error(t, err, "expected no registered pending transfers when coins is empty")

	// ARRANGE
	coins := sdk.NewCoins(sdk.NewInt64Coin("unobl", 1))

	// ACT
	toAddr, err = k.SendRestrictionFn(ctx, sdk.AccAddress{}, acc.GetAddress(), coins)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, acc.GetAddress(), toAddr, "expected the returned address unaltered")

	_, err = k.PendingTransfers.Get(ctx, acc.GetAddress().String())
	assert.Error(t, err, "expected no registered pending transfers when coins does not contain minting denom")

	// ARRANGE
	coins = coins.Add(sdk.NewInt64Coin("uusdc", 1))

	// ACT
	toAddr, err = k.SendRestrictionFn(ctx, sdk.AccAddress{}, acc.GetAddress(), coins)

	// ASSERT
	assert.NoError(t, err)
	assert.Equal(t, acc.GetAddress(), toAddr, "expected the returned address unaltered")

	_, err = k.PendingTransfers.Get(ctx, acc.GetAddress().String())
	assert.NoError(t, err, "expected one pending transfer associated with the address")
}
