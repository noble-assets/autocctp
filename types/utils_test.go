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
	"bytes"
	"testing"

	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"autocctp.dev/types"
	"autocctp.dev/utils"
)

func TestValidateMintRecipient(t *testing.T) {
	testCases := []struct {
		name        string
		address     []byte
		errContains string
	}{
		{
			name:        "fail when the mint recipient is nil",
			address:     nil,
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "fail when the mint recipient is empty initialized bytes",
			address:     make([]byte, 0),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "fail when the mint recipient is zero address",
			address:     make([]byte, cctptypes.MintRecipientLen),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "fail when the mint recipient is less than 32 bytes",
			address:     bytes.Repeat([]byte{0x01}, cctptypes.MintRecipientLen-1),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "fail when the mint recipient is more than 32 bytes",
			address:     bytes.Repeat([]byte{0x01}, cctptypes.MintRecipientLen+1),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "pass when the mint recipient is valid",
			address:     bytes.Repeat([]byte{0x01}, cctptypes.MintRecipientLen),
			errContains: "",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			err := types.ValidateMintRecipient(tC.address)
			if tC.errContains == "" {
				require.NoError(t, err, "expected no error validating mint recipient")
			} else {
				require.Error(t, err, "expected an error validating mint recipient")
				require.ErrorContains(t, err, tC.errContains, "expected a different error")
			}
		})
	}
}

func TestValidateDestinationCaller(t *testing.T) {
	testCases := []struct {
		name        string
		address     []byte
		errContains string
	}{
		{
			name:        "fail with zero address of 32 bytes",
			address:     make([]byte, cctptypes.DestinationCallerLen),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "fail with address shorter than 32 bytes",
			address:     bytes.Repeat([]byte{0x01}, cctptypes.DestinationCallerLen-1),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "fail with address longer than 32 bytes",
			address:     bytes.Repeat([]byte{0x01}, cctptypes.DestinationCallerLen+1),
			errContains: errorstypes.ErrInvalidAddress.Error(),
		},
		{
			name:        "pass with empty address",
			address:     nil,
			errContains: "",
		},
		{
			name:        "pass with zero-length address",
			address:     []byte{},
			errContains: "",
		},
		{
			name:        "pass with valid 32-byte non-zero address",
			address:     bytes.Repeat([]byte{0x01}, cctptypes.DestinationCallerLen),
			errContains: "",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			err := types.ValidateDestinationCaller(tC.address)
			if tC.errContains == "" {
				require.NoError(t, err, "expected no error validating destination caller")
			} else {
				require.Error(t, err, "expected an error validating destination caller")
				require.ErrorContains(t, err, tC.errContains, "expected a different error")
			}
		})
	}
}

func TestValidateExistingAccount(t *testing.T) {
	// ARRANGE: Create a new account
	addr := sdk.AccAddress(utils.AddressBytesTest())
	baseAcc := &authtypes.BaseAccount{Address: addr.String()}

	// ACT: Validate account
	err := types.ValidateExistingAccount(baseAcc, addr)

	// ASSERT
	require.NoError(t, err)

	// ARRANGE: Create a new account with non zero sequence
	err = baseAcc.SetSequence(1)
	require.NoError(t, err, "expected no error setting the sequence")

	// ACT: Validate account
	err = types.ValidateExistingAccount(baseAcc, addr)

	// ASSERT
	require.Error(t, err, "expecting an error when the account is new but with a non zero sequence")
	require.ErrorContains(t, err, "attempting to register an existing user")

	// ARRANGE: Account created signerlessly
	baseAcc = &authtypes.BaseAccount{Address: addr.String()}
	err = baseAcc.SetPubKey(&types.PubKey{Key: addr})
	require.NoError(t, err, "expected no error setting the pub key")
	err = baseAcc.SetSequence(1)
	require.NoError(t, err, "expected no error setting the sequence")

	// ACT: Validate account
	err = types.ValidateExistingAccount(baseAcc, addr)

	// ASSERT
	require.NoError(t, err, "expecting no error when account was created singerlessly")

	// ARRANGE: Change the expected address from previous test
	expAddress := sdk.AccAddress(utils.AddressBytesTest())

	// ACT: Validate account
	err = types.ValidateExistingAccount(baseAcc, expAddress)

	// ASSERT
	require.Error(t, err, "expecting an error when the address is different than the pub key")
	require.ErrorContains(t, err, "attempting to register an existing user")

	// ARRANGE: The account as an invalid pub key type
	err = baseAcc.SetPubKey(secp256k1.GenPrivKey().PubKey())
	require.NoError(t, err, "expected no error setting the pub key")

	// ACT: Validate account
	err = types.ValidateExistingAccount(baseAcc, addr)

	// ASSERT
	require.Error(t, err, "expecting an error when the pub key type is wrong")
	require.ErrorContains(t, err, "attempting to register an existing user")
}
