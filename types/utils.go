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
	"bytes"
	"errors"
	"fmt"
	"strconv"

	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateExistingAccount is a utility for checking if an existing account is eligible to
// be registered as an AutoCCTP account.
//
// A valid account must satisfy one of the following conditions.
//
// 1. Is a new account:
//   - A nil PubKey, i.e. the account never sent a transaction AND
//   - A 0 sequence.
//
// 2. Is an account registered signerlessy:
//   - A non nil PubKey with the custom type.
//   - Can have any sequence value.
func ValidateExistingAccount(account sdk.AccountI, address sdk.AccAddress) error {
	pubKey := account.GetPubKey()

	isNewAccount := pubKey == nil && account.GetSequence() == 0
	isValidPubKey := pubKey != nil && pubKey.Equals(&PubKey{Key: address})

	if !isNewAccount && !isValidPubKey {
		return fmt.Errorf("attempting to register an existing user account with address: %s", address.String())
	}
	return nil
}

func ValidateMintRecipient(address []byte) error {
	emptyByteArr := make([]byte, cctptypes.MintRecipientLen)
	if len(address) != cctptypes.MintRecipientLen || bytes.Equal(address, emptyByteArr) {
		return errorstypes.ErrInvalidAddress.Wrap("must be 32 bytes different than the zero address")
	}
	return nil
}

func ValidateDestinationCaller(address []byte) error {
	emptyByteArr := make([]byte, cctptypes.DestinationCallerLen)
	if len(address) != 0 {
		if len(address) != cctptypes.DestinationCallerLen || bytes.Equal(address, emptyByteArr) {
			return errorstypes.ErrInvalidAddress.Wrap("must be 32 bytes different than the zero address")
		}
	}
	return nil
}

func ParseDestinationDomain(destinationDomain string) (uint32, error) {
	dD, err := strconv.ParseUint(destinationDomain, cctptypes.BaseTen, cctptypes.DomainBitLen)
	if err != nil {
		return 0, fmt.Errorf("invalid destination domain: %w", err)
	}

	return uint32(dD), nil
}

// ValidateDestinationDomain returns a Domain type or an error if the domain is Noble or invalid.
func ValidateDestinationDomain(destinationDomain uint32) (Domain, error) {
	domain, supported := supportedDomains[destinationDomain]
	if !supported {
		return 0, fmt.Errorf("destination domain %d is not supported", destinationDomain)
	}
	if domain == NOBLE {
		return 0, errors.New("destination domain cannot be source domain")
	}

	return domain, nil
}

// isHex returns true if str begins with '0x' or '0X', and false otherwise.
func isHex(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// LeftPadBytes left pads a byte array to be length 32.
func LeftPadBytes(bz []byte) ([]byte, error) {
	if len(bz) > 32 {
		return nil, fmt.Errorf("padding error, expected less than 32 bytes, got %d", len(bz))
	}
	if len(bz) == 32 {
		return bz, nil
	}

	res := make([]byte, 32)
	copy(res[32-len(bz):], bz)
	return res, nil
}

func ValidateAndParseAccountFields(
	destinationDomain uint32,
	mintRecipient, fallbackRecipient, destinationCaller string,
) (*AccountProperties, error) {
	domain, err := ValidateDestinationDomain(destinationDomain)
	if err != nil {
		return nil, err
	}

	if len(mintRecipient) == 0 {
		return nil, errors.New("invalid mint recipient: cannot be empty")
	}
	recipient, err := domain.parseAddress(mintRecipient)
	if err != nil {
		return nil, fmt.Errorf("invalid mint recipient %s: %w", mintRecipient, err)
	}

	if _, err := sdk.AccAddressFromBech32(fallbackRecipient); err != nil {
		return nil, fmt.Errorf("invalid fallback recipient %s: %w", fallbackRecipient, err)
	}

	caller := []byte{}
	if len(destinationCaller) != 0 {
		caller, err = domain.parseAddress(destinationCaller)
		if err != nil {
			return nil, fmt.Errorf("invalid destination caller %s: %w", destinationCaller, err)
		}
	}

	return &AccountProperties{
		DestinationDomain: uint32(domain),
		MintRecipient:     recipient,
		FallbackRecipient: fallbackRecipient,
		DestinationCaller: caller,
	}, nil
}
