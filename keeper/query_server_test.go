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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"autocctp.dev/keeper"
	"autocctp.dev/types"
	"autocctp.dev/utils"
	"autocctp.dev/utils/mocks"
)

func TestAddress(t *testing.T) {
	utils.SDKConfigTest()

	validProperties := utils.ValidPropertiesTest(false)
	address := types.GenerateAddress(validProperties)

	validPropertiesWithCaller := utils.ValidPropertiesTest(true)
	addressWithCaller := types.GenerateAddress(validPropertiesWithCaller)

	testCases := []struct {
		name        string
		req         *types.QueryAddress
		setup       func(*mocks.AccountKeeper, sdk.Context)
		expResponse *types.QueryAddressResponse
		errContains string
	}{
		{
			name:        "fail with nil request",
			req:         nil,
			setup:       func(*mocks.AccountKeeper, sdk.Context) {},
			errContains: sdkerrors.ErrInvalidRequest.Error(),
		},
		{
			name:        "fail when valid properties fails",
			req:         &types.QueryAddress{},
			setup:       func(*mocks.AccountKeeper, sdk.Context) {},
			errContains: types.ErrInvalidAccountProperties.Error(),
		},
		{
			name: "valid request but account does not exists",
			req: &types.QueryAddress{
				DestinationDomain: validProperties.DestinationDomain,
				MintRecipient:     validProperties.MintRecipient,
				FallbackRecipient: validProperties.FallbackRecipient,
			},
			setup: func(ak *mocks.AccountKeeper, ctx sdk.Context) {},
			expResponse: &types.QueryAddressResponse{
				Address: address.String(),
				Exists:  false,
			},
		},
		{
			name: "valid request with nil destination caller",
			req: &types.QueryAddress{
				DestinationDomain: validProperties.DestinationDomain,
				MintRecipient:     validProperties.MintRecipient,
				FallbackRecipient: validProperties.FallbackRecipient,
			},
			setup: func(ak *mocks.AccountKeeper, ctx sdk.Context) {
				base := ak.NewAccountWithAddress(ctx, address)
				account := types.NewAccount(
					authtypes.NewBaseAccount(base.GetAddress(), base.GetPubKey(), base.GetAccountNumber(), base.GetSequence()),
					validProperties,
				)
				ak.Accounts[address.String()] = account
			},
			expResponse: &types.QueryAddressResponse{
				Address: address.String(),
				Exists:  true,
			},
		},
		{
			name: "valid request with complete data but address is not associated with autocctp account",
			req: &types.QueryAddress{
				DestinationDomain: validPropertiesWithCaller.DestinationDomain,
				MintRecipient:     validPropertiesWithCaller.MintRecipient,
				FallbackRecipient: validPropertiesWithCaller.FallbackRecipient,
				DestinationCaller: validPropertiesWithCaller.DestinationCaller,
			},
			setup: func(ak *mocks.AccountKeeper, ctx sdk.Context) {
				base := ak.NewAccountWithAddress(ctx, addressWithCaller)
				account := authtypes.NewBaseAccount(base.GetAddress(), base.GetPubKey(), base.GetAccountNumber(), base.GetSequence())
				ak.Accounts[address.String()] = account
			},
			expResponse: &types.QueryAddressResponse{
				Address: addressWithCaller.String(),
				Exists:  false,
			},
		},
		{
			name: "valid request with complete data",
			req: &types.QueryAddress{
				DestinationDomain: validPropertiesWithCaller.DestinationDomain,
				MintRecipient:     validPropertiesWithCaller.MintRecipient,
				FallbackRecipient: validPropertiesWithCaller.FallbackRecipient,
				DestinationCaller: validPropertiesWithCaller.DestinationCaller,
			},
			setup: func(ak *mocks.AccountKeeper, ctx sdk.Context) {
				base := ak.NewAccountWithAddress(ctx, addressWithCaller)
				account := types.NewAccount(
					authtypes.NewBaseAccount(base.GetAddress(), base.GetPubKey(), base.GetAccountNumber(), base.GetSequence()),
					validPropertiesWithCaller,
				)
				ak.Accounts[addressWithCaller.String()] = account
			},
			expResponse: &types.QueryAddressResponse{
				Address: addressWithCaller.String(),
				Exists:  true,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// ARRANGE
			mocks, k, ctx := mocks.AutoCCTPKeeper(t)
			server := keeper.NewQueryServer(k)
			tC.setup(mocks.AccountKeeper, ctx)

			// ACT
			resp, err := server.Address(ctx, tC.req)

			// ASSERT
			if tC.errContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errContains, "epxected a different error")
				require.Nil(t, resp, "expected nil response when receiving an error")
			} else {
				require.NoError(t, err, "expected no error")
				require.Equal(t, tC.expResponse, resp, "expected a different response")
			}
		})
	}
}

func TestStatsByDestinationDomain(t *testing.T) {
	destinationDomain := uint32(0)

	testCases := []struct {
		name        string
		req         *types.QueryStatsByDestinationDomain
		setup       func(*keeper.Keeper, sdk.Context)
		expResponse *types.QueryStatsByDestinationDomainResponse
		errContains string
	}{
		{
			name:        "fail with nil request",
			req:         nil,
			setup:       func(*keeper.Keeper, sdk.Context) {},
			errContains: sdkerrors.ErrInvalidRequest.Error(),
		},
		{
			name: "valid request but empty state",
			req: &types.QueryStatsByDestinationDomain{
				DestinationDomain: destinationDomain,
			},
			setup: func(*keeper.Keeper, sdk.Context) {},
			expResponse: &types.QueryStatsByDestinationDomainResponse{
				Accounts:         0,
				Transfers:        0,
				TotalTransferred: 0,
			},
		},
		{
			name: "valid request but no data for selected destination",
			req: &types.QueryStatsByDestinationDomain{
				DestinationDomain: destinationDomain,
			},
			setup: func(k *keeper.Keeper, ctx sdk.Context) {
				err := k.NumOfAccounts.Set(ctx, destinationDomain+1, 3)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfTransfers.Set(ctx, destinationDomain+1, 10)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.TotalTransferred.Set(ctx, destinationDomain+1, 1_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
			},
			expResponse: &types.QueryStatsByDestinationDomainResponse{
				Accounts:         0,
				Transfers:        0,
				TotalTransferred: 0,
			},
		},
		{
			name: "valid request with data only for selected destination",
			req: &types.QueryStatsByDestinationDomain{
				DestinationDomain: destinationDomain,
			},
			setup: func(k *keeper.Keeper, ctx sdk.Context) {
				err := k.NumOfAccounts.Set(ctx, destinationDomain, 3)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfTransfers.Set(ctx, destinationDomain, 10)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.TotalTransferred.Set(ctx, destinationDomain, 1_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
			},
			expResponse: &types.QueryStatsByDestinationDomainResponse{
				Accounts:         3,
				Transfers:        10,
				TotalTransferred: 1_000_000,
			},
		},
		{
			name: "valid request with data for multiple destinations",
			req: &types.QueryStatsByDestinationDomain{
				DestinationDomain: destinationDomain,
			},
			setup: func(k *keeper.Keeper, ctx sdk.Context) {
				err := k.NumOfAccounts.Set(ctx, destinationDomain, 3)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfAccounts.Set(ctx, destinationDomain+1, 3)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfTransfers.Set(ctx, destinationDomain, 10)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.NumOfTransfers.Set(ctx, destinationDomain+1, 10)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.TotalTransferred.Set(ctx, destinationDomain, 1_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
				err = k.TotalTransferred.Set(ctx, destinationDomain+1, 1_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
			},
			expResponse: &types.QueryStatsByDestinationDomainResponse{
				Accounts:         3,
				Transfers:        10,
				TotalTransferred: 1_000_000,
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// ARRANGE
			_, k, ctx := mocks.AutoCCTPKeeper(t)
			server := keeper.NewQueryServer(k)
			tC.setup(k, ctx)

			// ACT
			resp, err := server.StatsByDestinationDomain(ctx, tC.req)

			// ASSERT
			if tC.errContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errContains, "expected a different error")
				require.Nil(t, resp, "expected nil response when receiving an error")
			} else {
				require.NoError(t, err, "expected no error")
				require.Equal(t, tC.expResponse, resp, "expected a different response")
			}
		})
	}
}

func TestStats(t *testing.T) {
	testCases := []struct {
		name        string
		req         *types.QueryStats
		setup       func(*keeper.Keeper, sdk.Context)
		expResponse *types.QueryStatsResponse
		errContains string
	}{
		{
			name:        "invalid nil request",
			req:         nil,
			setup:       func(*keeper.Keeper, sdk.Context) {},
			errContains: sdkerrors.ErrInvalidRequest.Error(),
		},
		{
			name:  "valid request but empty state",
			req:   &types.QueryStats{},
			setup: func(*keeper.Keeper, sdk.Context) {},
			expResponse: &types.QueryStatsResponse{
				Stats: map[uint32]types.DomainStats{},
			},
		},
		{
			name: "valid request and data for only one destination",
			req:  &types.QueryStats{},
			setup: func(k *keeper.Keeper, ctx sdk.Context) {
				err := k.NumOfAccounts.Set(ctx, 0, 3)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfTransfers.Set(ctx, 0, 10)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.TotalTransferred.Set(ctx, 0, 1_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
			},
			expResponse: &types.QueryStatsResponse{
				Stats: map[uint32]types.DomainStats{
					0: {
						Accounts:         3,
						Transfers:        10,
						TotalTransferred: 1_000_000,
					},
				},
			},
		},
		{
			name: "valid request and data for multiple destinations",
			req:  &types.QueryStats{},
			setup: func(k *keeper.Keeper, ctx sdk.Context) {
				// Destination 0
				err := k.NumOfAccounts.Set(ctx, 0, 3)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfTransfers.Set(ctx, 0, 10)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.TotalTransferred.Set(ctx, 0, 1_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
				// Destination 1
				err = k.NumOfAccounts.Set(ctx, 1, 1)
				require.NoError(t, err, "expected no error setting the number of accounts")
				// Destination 3
				err = k.NumOfAccounts.Set(ctx, 3, 3_000)
				require.NoError(t, err, "expected no error setting the number of accounts")
				err = k.NumOfTransfers.Set(ctx, 3, 1_000_000)
				require.NoError(t, err, "expected no error setting the number of transfers")
				err = k.TotalTransferred.Set(ctx, 3, 1_000_000_000)
				require.NoError(t, err, "expected no error setting the total transferred")
			},
			expResponse: &types.QueryStatsResponse{
				Stats: map[uint32]types.DomainStats{
					0: {
						Accounts:         3,
						Transfers:        10,
						TotalTransferred: 1_000_000,
					},
					1: {
						Accounts:         1,
						Transfers:        0,
						TotalTransferred: 0,
					},
					3: {
						Accounts:         3_000,
						Transfers:        1_000_000,
						TotalTransferred: 1_000_000_000,
					},
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// ARRANGE
			_, k, ctx := mocks.AutoCCTPKeeper(t)
			server := keeper.NewQueryServer(k)
			tC.setup(k, ctx)

			// ACT
			resp, err := server.Stats(ctx, tC.req)

			// ASSERT
			if tC.errContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errContains, "expected a different error")
				require.Nil(t, resp, "expected nil response when receiving an error")
			} else {
				require.NoError(t, err, "expected no error")
				require.Equal(t, tC.expResponse, resp, "expected a different response")
			}
		})
	}
}
