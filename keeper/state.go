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

	"cosmossdk.io/math"

	"autocctp.dev/types"
)

// Setters

func (k *Keeper) IncrementNumOfAccounts(ctx context.Context, destinationDomain uint32) {
	count, _ := k.NumOfAccounts.Get(ctx, destinationDomain)

	if err := k.NumOfAccounts.Set(ctx, destinationDomain, count+1); err != nil {
		k.logger.Error("increment number of accounts", "destination_domain", destinationDomain)
	}
}

func (k *Keeper) IncrementNumOfTransfers(ctx context.Context, destinationDomain uint32) {
	count, _ := k.NumOfTransfers.Get(ctx, destinationDomain)

	if err := k.NumOfTransfers.Set(ctx, destinationDomain, count+1); err != nil {
		k.logger.Error("increment number of transfers", "destination_domain", destinationDomain)
	}
}

func (k *Keeper) IncrementTotalTransferred(ctx context.Context, destinationDomain uint32, amount math.Int) {
	if !amount.IsUint64() {
		k.logger.Error("increment total transferred because invalid amount",
			"destination_domain", destinationDomain,
			"amount", amount.String(),
		)
		return
	}

	tot, _ := k.TotalTransferred.Get(ctx, destinationDomain)
	if err := k.TotalTransferred.Set(ctx, destinationDomain, tot+amount.Uint64()); err != nil {
		k.logger.Error("increment total transferred", "destination_domain", destinationDomain)
	}
}

// Getters

func (k *Keeper) GetPendingTransfers(ctx context.Context) ([]types.Account, error) {
	accounts := []types.Account{}

	if err := k.PendingTransfers.Walk(ctx, nil, func(_ string, account types.Account) (stop bool, err error) {
		accounts = append(accounts, account)

		return false, nil
	}); err != nil {
		return []types.Account{}, err
	}

	return accounts, nil
}

func (k *Keeper) GetNumOfAccountPerDestination(ctx context.Context) (map[uint32]uint64, error) {
	numOfAccounts := make(map[uint32]uint64)

	if err := k.NumOfAccounts.Walk(ctx, nil, func(key uint32, value uint64) (stop bool, err error) {
		numOfAccounts[key] = value

		return false, nil
	}); err != nil {
		return nil, err
	}

	return numOfAccounts, nil
}

func (k *Keeper) GetNumOfTransfersPerDestination(ctx context.Context) (map[uint32]uint64, error) {
	numOfTransfers := make(map[uint32]uint64)

	if err := k.NumOfTransfers.Walk(ctx, nil, func(key uint32, value uint64) (stop bool, err error) {
		numOfTransfers[key] = value

		return false, nil
	}); err != nil {
		return nil, err
	}

	return numOfTransfers, nil
}

func (k *Keeper) GetTotalTransferredPerDestination(ctx context.Context) (map[uint32]uint64, error) {
	totTransferred := make(map[uint32]uint64)

	if err := k.TotalTransferred.Walk(ctx, nil, func(key uint32, value uint64) (stop bool, err error) {
		totTransferred[key] = value

		return false, nil
	}); err != nil {
		return nil, err
	}

	return totTransferred, nil
}
