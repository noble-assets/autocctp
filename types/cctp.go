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
	"context"
	"errors"
	"fmt"

	"github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/cosmos/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
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

func NewCCTPServer(msgServer CCTPMsgServer, queryServer CCTPQueryServer) CCTPService {
	if msgServer == nil {
		panic("CCTP msg server cannot be nil")
	}
	if queryServer == nil {
		panic("CCTP query server cannot be nil")
	}

	return &CCTPServer{
		MsgSever:    msgServer,
		QueryServer: queryServer,
	}
}

var _ CCTPService = CCTPServer{}

type CCTPServer struct {
	MsgSever    CCTPMsgServer
	QueryServer CCTPQueryServer
}

// DepositForBurn implements CCTPService.
func (c CCTPServer) DepositForBurn(ctx context.Context, msg *types.MsgDepositForBurn) (*types.MsgDepositForBurnResponse, error) {
	return c.MsgSever.DepositForBurn(ctx, msg)
}

// DepositForBurnWithCaller implements CCTPService.
func (c CCTPServer) DepositForBurnWithCaller(ctx context.Context, msg *types.MsgDepositForBurnWithCaller) (*types.MsgDepositForBurnWithCallerResponse, error) {
	return c.MsgSever.DepositForBurnWithCaller(ctx, msg)
}

// PerMessageBurnLimit implements CCTPService.
func (c CCTPServer) PerMessageBurnLimit(ctx context.Context, req *types.QueryGetPerMessageBurnLimitRequest) (*types.QueryGetPerMessageBurnLimitResponse, error) {
	return c.QueryServer.PerMessageBurnLimit(ctx, req)
}
