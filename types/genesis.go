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

package types

import (
	"errors"
	"fmt"
	"sort"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

func (gs *GenesisState) Validate() error {
	keysNumOfAccounts := make([]uint32, 0, len(gs.NumOfAccounts))
	for k := range gs.NumOfAccounts {
		keysNumOfAccounts = append(keysNumOfAccounts, k)
	}

	keysNumOfTransfers := make([]uint32, 0, len(gs.NumOfTransfers))
	for k := range gs.NumOfTransfers {
		keysNumOfTransfers = append(keysNumOfTransfers, k)
	}

	keysTotalTransferred := make([]uint32, 0, len(gs.TotalTransferred))
	for k := range gs.TotalTransferred {
		keysTotalTransferred = append(keysTotalTransferred, k)
	}

	sort.Slice(keysNumOfAccounts, func(i, j int) bool {
		return keysNumOfAccounts[i] < keysNumOfAccounts[j]
	})
	sort.Slice(keysNumOfTransfers, func(i, j int) bool {
		return keysNumOfTransfers[i] < keysNumOfTransfers[j]
	})
	sort.Slice(keysTotalTransferred, func(i, j int) bool {
		return keysTotalTransferred[i] < keysTotalTransferred[j]
	})

	if len(keysTotalTransferred) != len(keysNumOfTransfers) {
		return fmt.Errorf(
			"num of transfers and total transferred should have the same number of stored destinations: %d != %d",
			len(keysNumOfTransfers),
			len(keysTotalTransferred),
		)
	}

	// We can have destination domains with accounts registered and zero transfers or amount transferred
	// but not the opposite.
	if len(keysNumOfTransfers) > len(keysNumOfAccounts) {
		return errors.New("num of transfers has destination domains without accounts")
	}

	for _, keyTotalTransferred := range keysTotalTransferred {
		tot := gs.TotalTransferred[keyTotalTransferred]
		if tot == 0 {
			return fmt.Errorf("trying to register 0 total transferred for destination domain %d", keysTotalTransferred)
		}

		num, found := gs.NumOfTransfers[keyTotalTransferred]
		if !found {
			return fmt.Errorf(
				"destination domain %d is present in total transferred but not in num of transfers",
				keyTotalTransferred,
			)
		}
		// If we have amount transferred to one destination domain, we also must have at least
		// one transfer registered.
		if num == 0 {
			return fmt.Errorf(
				"trying to register total transferred without transfers for destination domain %d",
				keyTotalTransferred,
			)
		}

		acc, found := gs.NumOfAccounts[keyTotalTransferred]
		if !found || acc == 0 {
			return fmt.Errorf(
				"cannot have transfers for destination domain %d without registered accounts",
				keyTotalTransferred,
			)
		}
	}

	return nil
}
