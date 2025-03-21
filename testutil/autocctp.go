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

package testutil

import (
	"context"
	crand "crypto/rand"
	"math/rand"
	"strconv"

	"autocctp.dev/keeper"
)

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func randomDestinationDomain() uint32 {
	return rand.Uint32() % 10
}

// PendingTransfers generates a specified number of dummy pending transfers
// and adds them to the state. The parameters `num` controls the number of dummy accounts to
// create and `withCaller` whether the generated AutoCCTP accounts should have an
// associated destination caller.
//
// It returns a slice containing the addresses of the inserted accounts or an error if
// the insertion fails.
func PendingTransfers(ctx context.Context, k *keeper.Keeper, num int, destinationDomain string, withCaller bool) ([]string, error) {
	var addresses []string
	for range num {
		acc := AutoCCTPAccount(withCaller)
		if destinationDomain != "" {
			d, err := strconv.ParseUint(destinationDomain, 10, 32)
			if err != nil {
				return []string{}, err
			}
			acc.DestinationDomain = uint32(d)
		}
		if err := k.PendingTransfers.Set(ctx, acc.Address, acc); err != nil {
			return []string{}, err
		}
		addresses = append(addresses, acc.Address)
	}
	return addresses, nil
}
