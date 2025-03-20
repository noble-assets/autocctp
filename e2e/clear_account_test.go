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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"autocctp.dev/client/cli"
	"autocctp.dev/types"
)

// TestClearAccount_ToFallbackRecipient tests that an AutoCCTP account can be correctly cleared by
// sending funds to the fallback account.
func TestClearAccount_ToFallbackRecipient(t *testing.T) {
	t.Parallel()
	// ARRANGE
	ctx, s := NewAutoCCTPSuite(t, true, false)
	val := s.Chain.Validators[0]
	destinationDomain := fmt.Sprintf("%d", s.destinationDomain)

	// Register the AutoCCTP account.
	_ = s.RegisterAutoCCTPAccount(t, ctx, val, destinationDomain, s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")

	address, exists := GetAutoCCTPAccount(t, ctx, val, destinationDomain, s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")
	require.True(t, exists, "expected the new AutoCCTP account registered")

	// Pause the CCTP module to cause a failure in the clearing of the pending transfers.
	// This way, we can manually clear the AutoCCTP account with the tested tx.
	_ = s.PauseBurningAndMinting(t, ctx, val, s.CircleRoles.Pauser.KeyName())
	resp := GetCCTPBurningAndMintingPaused(t, ctx, val)
	require.True(t, resp.Paused.Paused, "expected the CCTP module to be paused")

	// Transfer funds to the AutCCTP account to be able to clear to fallback address.
	transferAmt := math.NewInt(1_000_000)
	err := val.BankSend(ctx, s.sender.KeyName(), ibc.WalletAmount{
		Address: address,
		Denom:   "uusdc",
		Amount:  transferAmt,
	})
	require.NoError(t, err, "expected no error funding the AutoCCTP account")

	initAmt, err := s.Chain.BankQueryBalance(ctx, s.fallbackRecipient.FormattedAddress(), "uusdc")
	require.NoError(t, err, "expected no error getting fallback recipient initial balance")

	// Restore CCTP unpaused condition.
	_ = s.UpnauseBurningAndMinting(t, ctx, val, s.CircleRoles.Pauser.KeyName())
	resp = GetCCTPBurningAndMintingPaused(t, ctx, val)
	require.False(t, resp.Paused.Paused, "expected the CCTP module to be unpaused")

	// ACT
	_, err = s.ClearAutoCCTPAccount(t, ctx, val, s.sender.KeyName(), address, true)

	// ASSERT
	require.Error(t, err, "expected error when signer is not fallback")

	// ACT
	hash, err := s.ClearAutoCCTPAccount(t, ctx, val, s.fallbackRecipient.KeyName(), address, true)
	require.NoError(t, err, "expected no error clearing the account to fallback recipient")

	// ASSERT
	finalAmt, err := s.Chain.BankQueryBalance(ctx, s.fallbackRecipient.FormattedAddress(), "uusdc")
	require.NoError(t, err, "expected no error getting fallback recipient final balance")
	require.Equal(t, initAmt.Add(transferAmt), finalAmt, "expected the fallback address to have received the funds")

	amt, err := s.Chain.BankQueryBalance(ctx, address, "uusdc")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), amt, "expected empty AutoCCTP account after clearing")

	tx := GetTx(t, ctx, val, hash)
	eventFound := false
	for _, rawEvent := range tx.Events {
		switch rawEvent.Type {
		case "noble.autocctp.v1.AccountCleared":
			event, err := sdk.ParseTypedEvent(rawEvent)
			require.NoError(t, err, "expected no error parsing the event")

			accountCleared, ok := event.(*types.AccountCleared)
			require.True(t, ok)

			require.Equal(t, address, accountCleared.Address, "expected a different address in the event")
			require.Equal(t, s.fallbackRecipient.FormattedAddress(), accountCleared.Receiver, "expected a different receiver in the event")
			eventFound = true
		}
	}
	require.True(t, eventFound, "expected account cleared event to be emitted")

	blockEvents := GetBlockResultsEvents(t, ctx, val, strconv.Itoa(int(tx.Height)))
	eventFound = false
	for _, event := range blockEvents {
		switch event.Type {
		case "circle.cctp.v1.DepositForBurn":
			eventFound = true
		}
	}
	require.False(t, eventFound, "expected no cctp event")
}

// TestClearAccount_ToMintRecipient tests that an AutoCCTP account can be correctly cleared
// re-trying the CCTP transfer.
func TestClearAccount_ToMintRecipient(t *testing.T) {
	t.Parallel()

	// ARRANGE
	ctx, s := NewAutoCCTPSuite(t, false, false)
	val := s.Chain.Validators[0]
	destinationDomain := fmt.Sprintf("%d", s.destinationDomain)

	_ = s.RegisterAutoCCTPAccount(t, ctx, val, destinationDomain, s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")

	address, exists := GetAutoCCTPAccount(t, ctx, val, destinationDomain, s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")
	require.True(t, exists, "expected the new AutoCCTP account registered")

	// Pause the CCTP module to cause a failure in the clearing of the pending transfers.
	// This way, we can manually clear the AutoCCTP account with the tested tx.
	_ = s.PauseBurningAndMinting(t, ctx, val, s.CircleRoles.Pauser.KeyName())
	resp := GetCCTPBurningAndMintingPaused(t, ctx, val)
	require.True(t, resp.Paused.Paused, "expected the CCTP module to be paused")

	transferAmt := math.NewInt(1_000_000)
	err := val.BankSend(ctx, s.sender.KeyName(), ibc.WalletAmount{
		Address: address,
		Denom:   "uusdc",
		Amount:  transferAmt,
	})
	require.NoError(t, err, "expected no error funding the AutoCCTP account")

	initAmt, err := s.Chain.BankQueryBalance(ctx, s.fallbackRecipient.FormattedAddress(), "uusdc")
	require.NoError(t, err, "expected no error getting fallback recipient initial balance")

	// Restore CCTP unpaused condition.
	_ = s.UpnauseBurningAndMinting(t, ctx, val, s.CircleRoles.Pauser.KeyName())
	resp = GetCCTPBurningAndMintingPaused(t, ctx, val)
	require.False(t, resp.Paused.Paused, "expected the CCTP module to be unpaused")

	// ACT
	hash, err := s.ClearAutoCCTPAccount(t, ctx, val, s.sender.KeyName(), address, false)
	require.NoError(t, err, "expected no error clearing the account")

	// ASSERT
	amt, err := s.Chain.BankQueryBalance(ctx, address, "uusdc")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), amt, "expected empty AutoCCTP account after clearing")

	finalAmt, err := s.Chain.BankQueryBalance(ctx, s.fallbackRecipient.FormattedAddress(), "uusdc")
	require.NoError(t, err, "expected no error getting fallback recipient final balance")
	require.Equal(t, initAmt, finalAmt, "expected the fallback address to have same initial funds")

	tx := GetTx(t, ctx, val, hash)
	eventFound := false
	for _, rawEvent := range tx.Events {
		switch rawEvent.Type {
		case "noble.autocctp.v1.AccountCleared":
			eventFound = true
		}
	}
	require.False(t, eventFound, "expected account cleared event to NOT be emitted")

	blockEvents := GetBlockResultsEvents(t, ctx, val, strconv.Itoa(int(tx.Height)))
	eventFound = false
	for _, event := range blockEvents {
		switch event.Type {
		case "circle.cctp.v1.DepositForBurn":

			// We have to iterate here since parsing does not work for CometBFT events. These events
			// have the additional `Attribute: Value` = `mode: EndBlock`
			for _, attribute := range event.Attributes {
				switch attribute.Key {
				case "amount":
					var actual string
					require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
					require.Equal(t, strconv.Itoa(int(transferAmt.Int64())), actual, "expected a different amount in cctp event")
				case "destination_domain":
					require.Equal(t, destinationDomain, attribute.Value, "expected a different destination domain in cctp event")
				case "destination_caller":
					var actual string
					require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
					require.Equal(t, "", actual, "expected a different destination caller in cctp event")
				case "mint_recipient":
					bz := common.FromHex(s.mintRecipient)
					bz, err := cli.LeftPadBytes(bz)
					require.NoError(t, err, "expected no error padding mint recipient address")
					expectedBase64 := base64.StdEncoding.EncodeToString(bz)
					var actual string
					require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
					require.Equal(t, expectedBase64, actual, "expected a different mint recipient in cctp event")
				}
			}

			eventFound = true
		}
	}
	require.True(t, eventFound, "expected cctp event")

	stats := GetAutoCCTPStatsByDestinationDomain(t, ctx, val, destinationDomain)
	require.Equal(t, uint64(1), stats.Accounts, "expected a different number of accounts")
	require.Equal(t, uint64(1), stats.Transfers, "expected a different number of transfers")
	require.Equal(t, transferAmt.Uint64(), stats.TotalTransferred, "expected a different total transferred")
}
