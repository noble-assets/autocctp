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

package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"autocctp.dev/types"
)

func TestGenesisState_Validate(t *testing.T) {
	testCases := []struct {
		name            string
		genesisModifier func(g *types.GenesisState)
		errContains     string
	}{
		{
			name:            "pass with default genesis",
			genesisModifier: func(g *types.GenesisState) {},
			errContains:     "",
		},
		{
			name: "fails when num of transfers keys > total transferred keys",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfTransfers = map[uint32]uint64{0: 10}
			},
			errContains: "should have the same number of stored destinations",
		},
		{
			name: "fails when num of transfers keys < total transferred keys",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfTransfers = map[uint32]uint64{0: 10}
				g.TotalTransferred = map[uint32]uint64{0: 10, 1: 10}
			},
			errContains: "should have the same number of stored destinations",
		},
		{
			name: "fails when num of transfers keys > num of accounts",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfTransfers = map[uint32]uint64{0: 10}
				g.TotalTransferred = map[uint32]uint64{0: 10}
			},
			errContains: "domains without accounts",
		},
		{
			name: "fails when total transferred is zero",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfAccounts = map[uint32]uint64{0: 10}
				g.NumOfTransfers = map[uint32]uint64{0: 10}
				g.TotalTransferred = map[uint32]uint64{0: 0}
			},
			errContains: "trying to register 0 total transferred",
		},
		{
			name: "fails when there are total transferred but not num of transfers",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfAccounts = map[uint32]uint64{0: 10}
				g.NumOfTransfers = map[uint32]uint64{1: 10}
				g.TotalTransferred = map[uint32]uint64{0: 10}
			},
			errContains: "but not in num of transfers",
		},
		{
			name: "fails when there are total transferred but num of transfers is zero",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfAccounts = map[uint32]uint64{0: 10}
				g.NumOfTransfers = map[uint32]uint64{0: 0}
				g.TotalTransferred = map[uint32]uint64{0: 10}
			},
			errContains: "without transfers",
		},
		{
			name: "fails when num of accounts is not registered for the domain",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfAccounts = map[uint32]uint64{1: 10}
				g.NumOfTransfers = map[uint32]uint64{0: 10}
				g.TotalTransferred = map[uint32]uint64{0: 10}
			},
			errContains: "without registered accounts",
		},
		{
			name: "fails when there are transfers for domain with zero accounts",
			genesisModifier: func(g *types.GenesisState) {
				g.NumOfAccounts = map[uint32]uint64{0: 0}
				g.NumOfTransfers = map[uint32]uint64{0: 10}
				g.TotalTransferred = map[uint32]uint64{0: 10}
			},
			errContains: "without registered accounts",
		},
	}

	for _, tC := range testCases {

		genesis := types.DefaultGenesisState()
		tC.genesisModifier(genesis)

		err := genesis.Validate()

		t.Run(tC.name, func(t *testing.T) {
			if tC.errContains != "" {
				require.Error(t, err, "expected an error validating the genesis")
				require.ErrorContains(t, err, tC.errContains, "expected a different error validating the genesis")
			} else {
				require.NoError(t, err, "expected no error validating the genesis")
			}
		})
	}
}
