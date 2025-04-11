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
	"strconv"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"autocctp.dev/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return queryServer{Keeper: keeper}
}

// Address implements types.QueryServer.
func (q queryServer) Address(ctx context.Context, req *types.QueryAddress) (*types.QueryAddressResponse, error) {
	if req == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("cannot be nil")
	}

	accountProperties, err := req.GetAccountProperties()
	if err != nil {
		return nil, types.ErrInvalidInputs.Wrap(err.Error())
	}
	if err := q.ValidateAccountProperties(accountProperties); err != nil {
		return nil, types.ErrInvalidAccountProperties.Wrap(err.Error())
	}

	address := types.GenerateAddress(accountProperties)

	exists := false
	if q.accountKeeper.HasAccount(ctx, address) {
		account := q.accountKeeper.GetAccount(ctx, address)
		_, exists = account.(*types.Account)
	}

	return &types.QueryAddressResponse{
		Address: address.String(),
		Exists:  exists,
	}, nil
}

// Stats implements types.QueryServer.
func (q queryServer) Stats(ctx context.Context, req *types.QueryStats) (*types.QueryStatsResponse, error) {
	if req == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("cannot be nil")
	}

	stats := make(map[uint32]types.DomainStats)

	numOfAccountsPerDestination, err := q.GetNumOfAccountPerDestination(ctx)
	if err != nil {
		return &types.QueryStatsResponse{}, err
	}

	for destinationDomain, numOfAccount := range numOfAccountsPerDestination {
		// We intentionally not return an error here since we favor partial information
		// compared to no information.
		numOfTransfers, err := q.NumOfTransfers.Get(ctx, destinationDomain)
		if err != nil {
			q.logger.Error("unable to get number of transfers", "destination domain", strconv.Itoa(int(destinationDomain)), "err", err)
		}

		totalTransferred, err := q.TotalTransferred.Get(ctx, destinationDomain)
		if err != nil {
			q.logger.Error("unable to get total transferred", "destination domain", strconv.Itoa(int(destinationDomain)), "err", err)
		}

		stats[destinationDomain] = types.DomainStats{
			Accounts:         numOfAccount,
			Transfers:        numOfTransfers,
			TotalTransferred: totalTransferred,
		}
	}

	return &types.QueryStatsResponse{Stats: stats}, nil
}

// StatsByDestinationDomain implements types.QueryServer.
func (q queryServer) StatsByDestinationDomain(ctx context.Context, req *types.QueryStatsByDestinationDomain) (*types.QueryStatsByDestinationDomainResponse, error) {
	if req == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("cannot be nil")
	}

	numOfAccount, err := q.NumOfAccounts.Get(ctx, req.DestinationDomain)
	if err != nil {
		q.logger.Error("unable to get num of accounts", "destination domain", strconv.Itoa(int(req.DestinationDomain)), "err", err)
	}
	numOfTransfers, err := q.NumOfTransfers.Get(ctx, req.DestinationDomain)
	if err != nil {
		q.logger.Error("unable to get num of transfers", "destination domain", strconv.Itoa(int(req.DestinationDomain)), "err", err)
	}
	totalTransferred, err := q.TotalTransferred.Get(ctx, req.DestinationDomain)
	if err != nil {
		q.logger.Error("unable to get total transferred", "destination domain", strconv.Itoa(int(req.DestinationDomain)), "err", err)
	}

	return &types.QueryStatsByDestinationDomainResponse{
		Accounts:         numOfAccount,
		Transfers:        numOfTransfers,
		TotalTransferred: totalTransferred,
	}, nil
}
