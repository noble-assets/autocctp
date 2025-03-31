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

package cli

import (
	"errors"
	"fmt"
	"strconv"

	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/btcutil/base58"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"autocctp.dev/types"
)

// Domain represents a destination domain supported by CCTP.
//
// https://developers.circle.com/stablecoins/supported-domains
type Domain uint32

const (
	ETHEREUM Domain = iota
	AVALANCHE
	OPTIMISM
	ARBITRUM
	NOBLE
	SOLANA
	BASE
	POLYGON
	SUI
	APTOS
	UNICHAIN
)

var supportedDomains = map[uint32]Domain{
	0:  ETHEREUM,
	1:  AVALANCHE,
	2:  OPTIMISM,
	3:  ARBITRUM,
	4:  NOBLE,
	5:  SOLANA,
	6:  BASE,
	7:  POLYGON,
	8:  SUI,
	9:  APTOS,
	10: UNICHAIN,
}

// ValidateDestinationDomain returns a Domain type or an error if the domain is Noble or invalid.
func ValidateDestinationDomain(destinationDomain string) (Domain, error) {
	dD, err := strconv.ParseUint(destinationDomain, cctptypes.BaseTen, cctptypes.DomainBitLen)
	if err != nil {
		return 0, fmt.Errorf("invalid destination domain: %w", err)
	}

	domain, supported := supportedDomains[uint32(dD)]
	if !supported {
		return 0, fmt.Errorf("destination domain %s is not supported", destinationDomain)
	}
	if domain == NOBLE {
		return 0, errors.New("destination domain cannot be source domain")
	}

	return domain, nil
}

// parseAddress parses an encoded address into a 32 length byte array. If the encoded bytes are
// less than 32, the function left pads to obtain the length used in cross-chain transfers.
func (d Domain) parseAddress(address string) ([]byte, error) {
	var bz []byte
	switch d {
	case ETHEREUM, AVALANCHE, OPTIMISM, ARBITRUM, BASE, POLYGON, SUI, APTOS, UNICHAIN:
		if !isHex(address) {
			return nil, errors.New("address not in hex format")
		}
		bz = common.FromHex(address)
	case SOLANA:
		bz = base58.Decode(address)
		if len(bz) == 0 {
			return nil, errors.New("address not valid base58")
		}
	case NOBLE:
		return nil, errors.New("destination domain cannot be source domain")
	default:
		return nil, fmt.Errorf("destination domain %d is not supported", uint32(d))
	}

	return LeftPadBytes(bz)
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
	destinationDomain, mintRecipient, fallbackRecipient, destinationCaller string,
) (*types.AccountProperties, error) {
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

	return &types.AccountProperties{
		DestinationDomain: uint32(domain),
		MintRecipient:     recipient,
		FallbackRecipient: fallbackRecipient,
		DestinationCaller: caller,
	}, nil
}
