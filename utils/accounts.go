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
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"autocctp.dev/types"
)

func ValidPropertiesTest(withCaller bool) types.AccountProperties {
	mintRecipient := make([]byte, 32)
	copy(mintRecipient[12:], sdk.AccAddress(AddressBytesTest()))

	fallbackRecipient := AddressTest()

	destinationCaller := []byte{}
	if withCaller {
		destinationCaller = make([]byte, 32)
		copy(destinationCaller[12:], sdk.AccAddress(AddressBytesTest()))
	}

	return types.AccountProperties{
		DestinationDomain: 0,
		MintRecipient:     mintRecipient,
		FallbackRecipient: fallbackRecipient,
		DestinationCaller: destinationCaller,
	}
}

// DummyAccountTest returns a dummy AutoCCTP account for testing.
func DummyAccountTest(withCaller bool) types.Account {
	accAddr := sdk.AccAddress(AddressBytesTest())
	baseAcc := authtypes.NewBaseAccountWithAddress(accAddr)
	acc := types.Account{
		BaseAccount:       baseAcc,
		DestinationDomain: randomDestinationDomainTest(),
		FallbackRecipient: AddressTest(),
		MintRecipient:     randomBytesTest(32),
	}

	if withCaller {
		acc.DestinationCaller = randomBytesTest(32)
	}

	return acc
}

// AddressTest is a test util to generate a bech32 address with "noble" prefix.
func AddressTest() string {
	return generateNobleAddressTest(AddressBytesTest())
}

// AddressBytesTest is a test util which returns the bytes of an address from a private key
// generated using the secp256k1 algorithm.
func AddressBytesTest() []byte {
	key := secp256k1.GenPrivKey()
	return key.PubKey().Address().Bytes()
}

func generateNobleAddressTest(bytes []byte) string {
	address, _ := sdk.Bech32ifyAddressBytes("noble", bytes)
	return address
}

// SDKConfigTest is a test util which set the Cosmos SDK Config to use
// "noble" prefix for accounts without sealing it.
func SDKConfigTest() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("noble", "noblepub")
}
