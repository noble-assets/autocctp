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

package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"autocctp.dev/client/cli"
	"autocctp.dev/types"
	"autocctp.dev/utils"
)

func TestValidateAndParseDomainFields(t *testing.T) {
	// ARRANGE
	utils.SDKConfigTest()

	testCases := []struct {
		name              string
		destinationDomain string
		mintRecipient     string
		fallbackRecipient string
		destinationCaller string
		errContains       string
		postChecks        func(*types.AccountProperties)
	}{
		{
			name:              "fail when the destination domain is not a number",
			destinationDomain: "",
			errContains:       "invalid destination domain",
		},
		{
			name:              "fail when the destination domain is not supported",
			destinationDomain: "11",
			errContains:       "not supported",
		},
		{
			name:              "fail when the destination domain is noble",
			destinationDomain: "4",
			errContains:       "cannot be source domain",
		},
		{
			name:              "fail when the mint recipient is empty",
			destinationDomain: "0",
			mintRecipient:     "",
			errContains:       "cannot be empty",
		},
		{
			name:              "fail when destination chain is ethereum and mint recipient is a solana address",
			destinationDomain: "1",
			mintRecipient:     "2WjnnBcYf4ff9xyDoH8yevnKF3yhH98DCcdy6PSmjNDa",
			errContains:       "address not in hex format",
		},
		{
			name:              "fail when fallback recipient is empty",
			destinationDomain: "0",
			mintRecipient:     "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			fallbackRecipient: "",
			errContains:       "invalid fallback recipient",
		},
		{
			name:              "fail when fallback recipient is not chain address",
			destinationDomain: "0",
			mintRecipient:     "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			fallbackRecipient: "cosmos1y5azhw4a99s4tm4kwzfwus52tjlvsaywuq3q3m",
			errContains:       "invalid Bech32 prefix",
		},
		{
			name:              "fail when destination caller is not empty and not valid",
			destinationDomain: "0",
			mintRecipient:     "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			fallbackRecipient: "noble1h8tqx833l3t2s45mwxjz29r85dcevy93wk63za",
			destinationCaller: "invalid",
			errContains:       "invalid destination caller",
		},
		{
			name:              "fail when destination caller is not a destination address",
			destinationDomain: "0",
			mintRecipient:     "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			fallbackRecipient: "noble1h8tqx833l3t2s45mwxjz29r85dcevy93wk63za",
			destinationCaller: "2WjnnBcYf4ff9xyDoH8yevnKF3yhH98DCcdy6PSmjNDa",
			errContains:       "invalid destination caller",
		},
		{
			name:              "success when mint recipient is an ethereum address",
			destinationDomain: "0",
			mintRecipient:     "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			fallbackRecipient: "noble1h8tqx833l3t2s45mwxjz29r85dcevy93wk63za",
			errContains:       "",
			postChecks: func(aP *types.AccountProperties) {
				require.Equal(t, 32, len(aP.MintRecipient), "expected mint recipient 32 bytes")
				require.Equal(t, 0, len(aP.DestinationCaller), "expected empty destinationc caller")
			},
		},
		{
			name:              "success when addresses are ethereum address",
			destinationDomain: "0",
			mintRecipient:     "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			fallbackRecipient: "noble1h8tqx833l3t2s45mwxjz29r85dcevy93wk63za",
			destinationCaller: "0xaB537dC791355d986A4f7a9a53f3D8810fd870D1",
			errContains:       "",
			postChecks: func(aP *types.AccountProperties) {
				require.Equal(t, 32, len(aP.MintRecipient), "expected mint recipient 32 bytes")
				require.Equal(t, 32, len(aP.DestinationCaller), "expected destination caller 32 bytes")
			},
		},
		{
			name:              "success when mint recipient is an aptos address",
			destinationDomain: "9",
			mintRecipient:     "0xeeff357ea5c1a4e7bc11b2b17ff2dc2dcca69750bfef1e1ebcaccf8c8018175b",
			fallbackRecipient: "noble1h8tqx833l3t2s45mwxjz29r85dcevy93wk63za",
			errContains:       "",
			postChecks: func(aP *types.AccountProperties) {
				require.Equal(t, 32, len(aP.MintRecipient), "expected mint recipient 32 bytes")
				require.Equal(t, 0, len(aP.DestinationCaller), "expected empty destinationc caller")
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// ACT
			accountProperties, err := cli.ValidateAndParseAccountFields(
				tC.destinationDomain,
				tC.mintRecipient,
				tC.fallbackRecipient,
				tC.destinationCaller,
			)

			// ASSERT
			if tC.errContains != "" {
				require.Error(t, err, "expected an error")
				require.ErrorContains(t, err, tC.errContains, "epxected a different error")
				require.Nil(t, accountProperties, "expected nil response when receiving an error")
			} else {
				require.NoError(t, err, "expected no error")
				tC.postChecks(accountProperties)
			}
		})
	}
}
