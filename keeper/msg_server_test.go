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
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"autocctp.dev/keeper"
	"autocctp.dev/types"
	"autocctp.dev/utils"
	"autocctp.dev/utils/mocks"
)

func TestRegisterAccount_NewAccount(t *testing.T) {
	// ARRANGE
	signer := utils.AddressTest()
	accountProperties := utils.ValidPropertiesTest(false)
	customAddress := types.GenerateAddress(accountProperties)

	invalidFallbackRecipient := "cosmos1y5azhw4a99s4tm4kwzfwus52tjlvsaywuq3q3m"

	testCases := []struct {
		mode       string
		msg        func(types.AccountProperties) interface{}
		serverCall func(types.MsgServer, context.Context, interface{}) (string, error)
	}{
		{
			mode: "standard",
			msg: func(ap types.AccountProperties) interface{} {
				return types.MsgRegisterAccount{
					Signer:            signer,
					DestinationDomain: ap.DestinationDomain,
					MintRecipient:     ap.MintRecipient,
					FallbackRecipient: ap.FallbackRecipient,
				}
			},
			serverCall: func(s types.MsgServer, ctx context.Context, msgI interface{}) (string, error) {
				msg := msgI.(types.MsgRegisterAccount)
				resp, err := s.RegisterAccount(ctx, &msg)
				if err != nil {
					return "", err
				}
				return resp.Address, nil
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.mode, func(t *testing.T) {
			validMsg := tC.msg(accountProperties)

			mocks, k, ctx := mocks.AutoCCTPKeeper(t)
			server := keeper.NewMsgServer(k)

			// ACT
			invalidProp := accountProperties
			invalidProp.FallbackRecipient = invalidFallbackRecipient

			invalidMsg := tC.msg(invalidProp)
			_, err := tC.serverCall(server, ctx, invalidMsg)

			// ASSERT
			assert.Error(t, err, "expected error when account properties validation fails")
			assert.ErrorContains(t, err, types.ErrInvalidAccountProperties.Error())

			// ARRANGE

			// ACT: Register a new account succeed.
			resp, err := tC.serverCall(server, ctx, validMsg)

			// ASSERT: One account has been added but no pending transfers.
			assert.NoError(t, err, "expected no error during account registration")
			assert.Equal(t, customAddress.String(), resp, "expected a different address returned")

			nAccount, _ := k.NumOfAccounts.Get(ctx, accountProperties.DestinationDomain)
			assert.Equal(t, uint64(1), nAccount, "expected only one account registered")

			_, err = k.PendingTransfers.Get(ctx, customAddress.String())
			assert.Error(t, err, "expected no registered pending transfers")

			acc, found := mocks.AccountKeeper.Accounts[customAddress.String()]
			assert.True(t, found, "the account should be registered")

			_, ok := acc.(*authtypes.BaseAccount)
			assert.False(t, ok, "expected the account to not be a base account")

			_, ok = acc.(*types.Account)
			assert.True(t, ok, "expected the account to be of custom type")
		})
	}
}

func TestRegisterAccount_ExistingAccount(t *testing.T) {
	// ARRANGE
	signer := utils.AddressTest()

	accountProperties := utils.ValidPropertiesTest(false)
	customAddress := types.GenerateAddress(accountProperties)

	testCases := []struct {
		mode       string
		msg        func() interface{}
		serverCall func(types.MsgServer, context.Context, interface{}) (interface{}, error)
	}{
		{
			mode: "standard",
			msg: func() interface{} {
				return types.MsgRegisterAccount{
					Signer:            signer,
					DestinationDomain: accountProperties.DestinationDomain,
					MintRecipient:     accountProperties.MintRecipient,
					FallbackRecipient: accountProperties.FallbackRecipient,
				}
			},
			serverCall: func(s types.MsgServer, ctx context.Context, msgI interface{}) (interface{}, error) {
				msg := msgI.(types.MsgRegisterAccount)
				return s.RegisterAccount(ctx, &msg)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.mode, func(t *testing.T) {
			msg := tC.msg()
			m, k, ctx := mocks.AutoCCTPKeeper(t)
			server := keeper.NewMsgServer(k)

			// Simulate an account that has been registered after receiving funds via x/bank with
			// an invalid sequence. If the account has only received funds, the sequence should be 0.
			registeredAcc := m.AccountKeeper.NewAccountWithAddress(ctx, customAddress)
			err := registeredAcc.SetSequence(1)
			require.NoError(t, err, "expected no error setting the sequence")

			m.AccountKeeper.Accounts[registeredAcc.GetAddress().String()] = registeredAcc
			m.BankKeeper.Balances[registeredAcc.GetAddress().String()] = sdk.Coins{sdk.NewInt64Coin("uusdc", 1)}

			// ACT
			_, err = tC.serverCall(server, ctx, msg)

			// ASSERT
			assert.Error(t, err, "expected error when new account has a sequence different than zero")
			assert.ErrorContains(t, err, "error validating existing account")

			// ARRANGE: Set the sequence of the previous account to zero.
			err = registeredAcc.SetSequence(0)
			require.NoError(t, err, "expected no error setting the sequence")

			acc := m.AccountKeeper.Accounts[customAddress.String()]
			_, ok := acc.(*authtypes.BaseAccount)
			assert.True(t, ok, "the account should be initially a base account")
			_, ok = acc.(*types.Account)
			assert.False(t, ok, "the account sould NOT be initially of custom type")

			// ACT
			_, err = tC.serverCall(server, ctx, msg)

			// ASSERT: One account has been added along with the pending transfer.
			assert.NoError(t, err, "expected no error during account registration")

			nAccount, _ := k.NumOfAccounts.Get(ctx, accountProperties.DestinationDomain)
			assert.Equal(t, uint64(1), nAccount, "expected only one account registered")

			_, err = k.PendingTransfers.Get(ctx, customAddress.String())
			assert.NoError(t, err, "expected new account added to pending transfers")

			// Verify correct account type update
			acc = m.AccountKeeper.Accounts[customAddress.String()]
			_, ok = acc.(*authtypes.BaseAccount)
			assert.False(t, ok, "the updated account should NOT be a base account")

			_, ok = acc.(*types.Account)
			assert.True(t, ok, "expected the account to be of custom type")

			// ARRANGE: Simulate previous account has been cleared and we are in a new block
			err = k.PendingTransfers.Clear(ctx, nil)
			assert.NoError(t, err)

			// ACT: Trying to register again the account fails.
			_, err = tC.serverCall(server, ctx, msg)

			// ASSERT: The function did't change the state.
			assert.Error(t, err, "expected error when account is already registered")
			assert.ErrorContains(t, err, "account has already been registered", "expected a different error")

			nAccount, _ = k.NumOfAccounts.Get(ctx, accountProperties.DestinationDomain)
			assert.Equal(t, uint64(1), nAccount, "expected no change in number of account")

			_, err = k.PendingTransfers.Get(ctx, customAddress.String())
			assert.Error(t, err, "no account should have been added to pending transfers")

			// ARRANGE: Create a new account but without funding it.
			mocks.ResetTest(t, ctx, k, m)

			customAddress = types.GenerateAddress(accountProperties)
			registeredAcc = m.AccountKeeper.NewAccountWithAddress(ctx, customAddress)
			m.AccountKeeper.Accounts[registeredAcc.GetAddress().String()] = registeredAcc

			// ACT
			_, err = tC.serverCall(server, ctx, msg)

			// ASSERT: One account has been added but no prending transfers because the balance was empty.
			assert.NoError(t, err, "expected no error during account registration")

			nAccount, _ = k.NumOfAccounts.Get(ctx, accountProperties.DestinationDomain)
			assert.Equal(t, uint64(1), nAccount, "expected only one account registered")

			_, err = k.PendingTransfers.Get(ctx, customAddress.String())
			assert.Error(t, err, "expected no pending transfers")

			// ARRANGE: Create a new account with unsupported type.
			mocks.ResetTest(t, ctx, k, m)

			customAddress = types.GenerateAddress(accountProperties)
			unsupportedAcc := authtypes.ModuleAccount{
				BaseAccount: authtypes.NewBaseAccountWithAddress(customAddress),
			}
			m.AccountKeeper.Accounts[unsupportedAcc.GetAddress().String()] = &unsupportedAcc

			// ACT
			_, err = tC.serverCall(server, ctx, msg)

			// ASSERT
			assert.Error(t, err, "expected error during account registration")
			assert.ErrorContains(t, err, "unsupported account type", "expected a different error")
		})
	}
}
