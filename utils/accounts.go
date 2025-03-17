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

package utils

import (
	crand "crypto/rand"
	"math/rand"

	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"autocctp.dev/types"
)

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := crand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func generateRandomDestinationDomain() uint32 {
	return rand.Uint32() % 10
}

// DummyAccount returns a dummy AutoCCTP account for testing.
func DummyAccount(withCaller bool) types.Account {
	accAddr := sdk.AccAddress(TestAddressBytes())
	baseAcc := authtypes.NewBaseAccountWithAddress(accAddr)
	acc := types.Account{
		BaseAccount:       baseAcc,
		DestinationDomain: generateRandomDestinationDomain(),
		FallbackRecipient: TestAddress(),
		MintRecipient:     generateRandomBytes(32),
	}

	if withCaller {
		acc.DestinationCaller = generateRandomBytes(32)
	}

	return acc
}

// TestAddress is a test util to generate a bech32 address with "noble" prefix.
func TestAddress() string {
	return generateNobleAddress(TestAddressBytes())
}

// TestAddressBytes is a test util which returns the bytes of an address from a private key
// generated using the secp256k1 algorithm.
func TestAddressBytes() []byte {
	key := secp256k1.GenPrivKey()
	return key.PubKey().Address().Bytes()
}

func generateNobleAddress(bytes []byte) string {
	address, _ := sdk.Bech32ifyAddressBytes("noble", bytes)
	return address
}

// TestSDKConfig is a test util which set the Cosmos SDK Config to use
// "noble" prefix for accounts without sealing it.
// TODO: move to init()
func TestSDKConfig() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("noble", "noblepub")
}
