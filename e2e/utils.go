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

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/stretchr/testify/require"

	abci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/jsonpb"

	"autocctp.dev/types"
)

func (s AutoCCTPSuite) RegisterAutoCCTPAccount(t *testing.T, ctx context.Context, validator *cosmos.ChainNode, destinationDomain, mintRecipient, fallbackRecipient, destinationCaller string) string {
	var err error
	var hash string
	if destinationCaller == "" {
		hash, err = validator.ExecTx(ctx, s.sender.KeyName(), "autocctp", "register-account", destinationDomain, mintRecipient, fallbackRecipient)
	} else {
		hash, err = validator.ExecTx(ctx, s.sender.KeyName(), "autocctp", "register-account", destinationDomain, mintRecipient, fallbackRecipient, destinationCaller)
	}
	require.NoError(t, err, "expected no error registering the autocctp account")
	return hash
}

func GetAutoCCTPAccount(t *testing.T, ctx context.Context, validator *cosmos.ChainNode, destinationDomain, mintRecipient, fallbackRecipient, destinationCaller string) (string, bool) {
	var raw []byte
	var err error

	if destinationCaller == "" {
		raw, _, err = validator.ExecQuery(ctx, "autocctp", "address", destinationDomain, mintRecipient, fallbackRecipient)
	} else {
		raw, _, err = validator.ExecQuery(ctx, "autocctp", "address", destinationDomain, mintRecipient, fallbackRecipient, destinationCaller)
	}
	require.NoError(t, err, "expected no error querying the autocctp account")

	var res types.QueryAddressResponse
	require.NoError(t, json.Unmarshal(raw, &res), "expected no error parsing address response")

	return res.Address, res.Exists
}

func GetAutoCCTPStatsByDestinationDomain(t *testing.T, ctx context.Context, validator *cosmos.ChainNode, destinationDomain string) types.QueryStatsByDestinationDomainResponse {
	var raw []byte
	var err error

	raw, _, err = validator.ExecQuery(ctx, "autocctp", "stats", destinationDomain)
	require.NoError(t, err, "expected no error querying the autocctp stats")

	var res types.QueryStatsByDestinationDomainResponse
	require.NoError(t, jsonpb.Unmarshal(bytes.NewReader(raw), &res), "expected no error parsing stats response")

	return res
}

type Fee struct {
	Amount sdk.Coins `json:"amount"`
}

type AuthInfo struct {
	Fee Fee `json:"fee"`
}

type Tx struct {
	AuthInfo AuthInfo `json:"auth_info"`
}

type TxResponse struct {
	Tx Tx `json:"tx"`
}

func TxFee(ctx context.Context, validator *cosmos.ChainNode, hash string) (sdk.Coins, error) {
	raw, _, err := validator.ExecQuery(ctx, "tx", hash)
	if err != nil {
		return nil, err
	}

	var res TxResponse
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return res.Tx.AuthInfo.Fee.Amount, nil
}

func QueryEvents(ctx context.Context, validator *cosmos.ChainNode, height string) ([]abci.Event, error) {
	raw, _, err := validator.ExecQuery(ctx, "block-results", height)
	if err != nil {
		return nil, err
	}

	var res coretypes.ResultBlockResults
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, err
	}

	return res.FinalizeBlockEvents, nil
}

func QueryTransaction(ctx context.Context, validator *cosmos.ChainNode, hash string) (*sdk.TxResponse, error) {
	raw, _, err := validator.ExecQuery(ctx, "tx", hash)
	if err != nil {
		return nil, err
	}

	var res sdk.TxResponse
	err = jsonpb.Unmarshal(bytes.NewReader(raw), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func QueryStats(ctx context.Context, validator *cosmos.ChainNode, destinationDomain string) (*types.QueryStatsByDestinationDomainResponse, error) {
	raw, _, err := validator.ExecQuery(ctx, "autocctp", "stats", destinationDomain)
	if err != nil {
		return nil, err
	}

	var res types.QueryStatsByDestinationDomainResponse
	err = jsonpb.Unmarshal(bytes.NewReader(raw), &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func QueryCCTPBurningAndMintingPaused(t *testing.T, ctx context.Context, validator *cosmos.ChainNode) *cctptypes.QueryGetBurningAndMintingPausedResponse {
	var raw []byte
	var err error
	raw, _, err = validator.ExecQuery(ctx, "cctp", "show-burning-and-minting-paused")

	require.NoError(t, err, "expected no error querying cctp burning and minting paused")

	var res cctptypes.QueryGetBurningAndMintingPausedResponse
	require.NoError(t, json.Unmarshal(raw, &res), "expected no error parsing burning and minting paused response")

	return &res
}
