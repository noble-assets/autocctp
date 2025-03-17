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
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"

	"autocctp.dev/client/cli"
	"autocctp.dev/types"
)

func TestRegisterAccount(t *testing.T) {
	t.Parallel()

	// ARRANGE
	ctx, s := NewAutoCCTPSuite(t, true, false)
	val := s.Chain.Validators[0]

	// ACT
	_, exists := GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")

	// ASSERT
	require.False(t, exists, "expected no autocctp account")

	// ACT
	hash, err := val.ExecTx(ctx, s.sender.KeyName(), "autocctp", "register-account", fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress())
	require.NoError(t, err)

	// ASSERT
	address, exists := GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")
	require.True(t, exists, "expected the new autocctp account registered")

	tx, err := QueryTransaction(ctx, val, hash)
	require.NoError(t, err, "expected no error querying the tx from hash")
	for _, rawEvent := range tx.Events {
		switch rawEvent.Type {
		case "noble.autocctp.v1.AccountRegistered":
			event, err := sdk.ParseTypedEvent(rawEvent)
			require.NoError(t, err, "expected no error parsing the event")

			accountRegistered, ok := event.(*types.AccountRegistered)
			require.True(t, ok)

			require.Equal(t, address, accountRegistered.Address)
			require.Equal(t, s.destinationDomain, accountRegistered.DestinationDomain)
			require.Equal(t, []byte(nil), accountRegistered.DestinationCaller)
			require.Equal(t, s.fallbackRecipient.FormattedAddress(), accountRegistered.FallbackRecipient)
			mintRecipient := common.BytesToAddress(accountRegistered.MintRecipient[12:]).String()
			require.Equal(t, s.mintRecipient, mintRecipient)
			require.False(t, accountRegistered.Signerlessly, "expected the account to be not registered signerlessly")
		}
	}
}

func TestRegisterAccountSignerlessly(t *testing.T) {
	t.Parallel()

	// ARRANGE
	ctx, s := NewAutoCCTPSuite(t, false, false)
	val := s.Chain.Validators[0]

	// ACT
	address, exists := GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")

	// ASSERT
	require.False(t, exists, "expected no autocctp account")

	// ACT
	_, err := val.ExecTx(ctx, s.sender.KeyName(), "autocctp", "register-account-signerlessly", fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress())

	// ASSERT
	require.Error(t, err, "expected an error when the autocctp account does not have funds to pay fees")

	// ARRANGE: the wannabe auto CCTP account must be registered and must have funds to pay fees.
	transferAmt := math.NewInt(1_000_000)
	err = val.BankSend(ctx, s.sender.KeyName(), ibc.WalletAmount{
		Address: address,
		Denom:   "uusdc",
		Amount:  transferAmt,
	})
	require.NoError(t, err, "expected no error funding the autocctp account")

	_, exists = GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")
	require.False(t, exists, "expected no auto cctp account but a base account registered")

	// ACT
	hash, err := val.ExecTx(ctx, s.sender.KeyName(), "autocctp", "register-account-signerlessly", fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress())
	require.NoError(t, err)

	// ASSERT
	_, exists = GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), "")
	require.True(t, exists, "expected the autocctp account registered")

	resp, err := QueryStats(ctx, val, fmt.Sprintf("%d", s.destinationDomain))
	require.NoError(t, err)

	require.Equal(t, uint64(1), resp.Accounts, "expected a different number of accounts")
	require.Equal(t, uint64(1), resp.Transfers, "expected a different number of transfers")

	fees, err := TxFee(ctx, val, hash)
	require.NoError(t, err)
	transferAmtNoFees := transferAmt.Sub(fees.AmountOf("uusdc"))
	require.Equal(t, transferAmtNoFees.Uint64(), resp.TotalTransferred, "expected total transfer equal to initial amount minus fees")

	tx, err := QueryTransaction(ctx, val, hash)
	require.NoError(t, err, "expected no error querying the tx from hash")
	for _, rawEvent := range tx.Events {
		switch rawEvent.Type {
		case "noble.autocctp.v1.AccountRegistered":
			event, err := sdk.ParseTypedEvent(rawEvent)
			require.NoError(t, err, "expected no error parsing the event")

			accountRegistered, ok := event.(*types.AccountRegistered)
			require.True(t, ok)

			require.Equal(t, address, accountRegistered.Address)
			require.Equal(t, s.destinationDomain, accountRegistered.DestinationDomain)
			require.Equal(t, []byte(nil), accountRegistered.DestinationCaller)
			require.Equal(t, s.fallbackRecipient.FormattedAddress(), accountRegistered.FallbackRecipient)
			mintRecipient := common.BytesToAddress(accountRegistered.MintRecipient[12:]).String()
			require.Equal(t, s.mintRecipient, mintRecipient)
			require.True(t, accountRegistered.Signerlessly, "expected the account to be registered signerlessly")
		}
	}

	stats := GetAutoCCTPStatsByDestinationDomain(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain))
	require.Equal(t, uint64(1), stats.Accounts, "expected a different number of accounts")
	require.Equal(t, uint64(1), stats.Transfers, "expected a different number of transfers")
	require.Equal(t, transferAmtNoFees.Uint64(), stats.TotalTransferred, "expected a different total transferred")
}

func TestFlowIBC(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		destinationCaller bool
	}{
		{
			name:              "no destination caller",
			destinationCaller: false,
		},
		{
			name:              "with destination caller",
			destinationCaller: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			// ARRANGE
			ctx, s := NewAutoCCTPSuite(t, true, true)
			val := s.Chain.Validators[0]

			destinationCaller := ""
			if tC.destinationCaller {
				destinationCaller = s.destinationCaller
			}

			// Transfer funds to the counterparty chain to have an account with USDC balance that can
			// send funds to the AutoCCTP account.
			autocctpToCounterpartyChannelInfo, err := s.IBC.Relayer.GetChannels(ctx, s.IBC.RelayerReporter, s.Chain.Config().ChainID)
			require.NoError(t, err)
			autocctpToCounterpartyChannelID := autocctpToCounterpartyChannelInfo[0].ChannelID

			counterpartyToAutocctpChannelInfo, err := s.IBC.Relayer.GetChannels(ctx, s.IBC.RelayerReporter, s.IBC.CounterpartyChain.Config().ChainID)
			require.NoError(t, err)
			counterpartyToAutocctpChannelID := counterpartyToAutocctpChannelInfo[0].ChannelID

			amountToSend := math.NewInt(1_000)
			transfer := ibc.WalletAmount{
				Address: s.IBC.Account.FormattedAddress(),
				Denom:   "uusdc",
				Amount:  amountToSend,
			}
			_, err = s.Chain.SendIBCTransfer(ctx, autocctpToCounterpartyChannelID, s.sender.KeyName(), transfer, ibc.TransferOptions{})
			require.NoError(t, err)
			require.NoError(t, s.IBC.Relayer.Flush(ctx, s.IBC.RelayerReporter, s.IBC.PathName, autocctpToCounterpartyChannelID), "expected no error relaying MsgRecvPacket & MsgAcknowledgement")

			srcDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", counterpartyToAutocctpChannelID, "uusdc"))
			dstIbcDenom := srcDenomTrace.IBCDenom()

			counterpartyWalletBal, err := s.IBC.CounterpartyChain.GetBalance(ctx, s.IBC.Account.FormattedAddress(), dstIbcDenom)
			require.NoError(t, err)
			require.Equal(t, transfer.Amount, counterpartyWalletBal)

			address, exists := GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), destinationCaller)
			require.False(t, exists, "expected no autocctp account")

			amt, err := s.Chain.BankQueryBalance(ctx, address, "uusdc")
			require.NoError(t, err)
			require.Equal(t, amt, math.ZeroInt(), "expected no initial balance in the autocctp account")

			// ACT
			ibcAmt1 := math.NewInt(100)
			transfer = ibc.WalletAmount{
				Address: address,
				Denom:   dstIbcDenom,
				Amount:  ibcAmt1,
			}
			_, err = s.IBC.CounterpartyChain.SendIBCTransfer(ctx, counterpartyToAutocctpChannelID, s.IBC.Account.KeyName(), transfer, ibc.TransferOptions{})
			require.NoError(t, err)
			require.NoError(t, s.IBC.Relayer.Flush(ctx, s.IBC.RelayerReporter, s.IBC.PathName, counterpartyToAutocctpChannelID), "expected no error relaying MsgRecvPacket & MsgAcknowledgement")

			// ASSERT
			counterpartyWalletBal, err = s.IBC.CounterpartyChain.GetBalance(ctx, s.IBC.Account.FormattedAddress(), dstIbcDenom)
			require.NoError(t, err)
			require.Equal(t, amountToSend.Sub(ibcAmt1), counterpartyWalletBal)

			amt, err = s.Chain.BankQueryBalance(ctx, address, "uusdc")
			require.NoError(t, err)
			require.Equal(t, amt, ibcAmt1, "expected the account to have received funds via IBC")

			// ACT
			hash := s.RegisterAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), destinationCaller)

			// ASSERT
			_, exists = GetAutoCCTPAccount(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain), s.mintRecipient, s.fallbackRecipient.FormattedAddress(), destinationCaller)
			require.True(t, exists, "expected the new autocctp account registered")

			tx, err := QueryTransaction(ctx, val, hash)
			require.NoError(t, err, "expected no error querying the tx from hash")
			blockEvents, err := QueryEvents(ctx, val, strconv.Itoa(int(tx.Height)))
			require.NoError(t, err, "expected no error querying block events")
			eventFound := false
			for _, event := range blockEvents {
				switch event.Type {
				case "circle.cctp.v1.DepositForBurn":
					for _, attribute := range event.Attributes {
						switch attribute.Key {
						case "amount":
							var actual string
							require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
							require.Equal(t, strconv.Itoa(int(ibcAmt1.Int64())), actual, "expected a different amount in cctp event")
						case "destination_domain":
							require.Equal(t, fmt.Sprintf("%d", s.destinationDomain), attribute.Value, "expected a different destination domain in cctp event")
						case "destination_caller":
							expectedBase64 := ""
							if tC.destinationCaller {
								bz := common.FromHex(destinationCaller)
								bz, err := cli.LeftPadBytes(bz)
								require.NoError(t, err, "expected no error padding destination caller address")
								expectedBase64 = base64.StdEncoding.EncodeToString(bz)
							}
							var actual string
							require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
							require.Equal(t, expectedBase64, actual, "expected a different destination caller in cctp event")
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

			amt, err = s.Chain.BankQueryBalance(ctx, address, "uusdc")
			require.NoError(t, err)
			require.Equal(t, math.ZeroInt(), amt, "expected the account to have zero balance after autocctp transfer")

			height, err := s.Chain.Height(ctx)
			require.NoError(t, err)

			// ACT
			ibcAmt2 := math.NewInt(300)
			transfer = ibc.WalletAmount{
				Address: address,
				Denom:   dstIbcDenom,
				Amount:  ibcAmt2,
			}
			_, err = s.IBC.CounterpartyChain.SendIBCTransfer(ctx, counterpartyToAutocctpChannelID, s.IBC.Account.KeyName(), transfer, ibc.TransferOptions{})
			require.NoError(t, err)
			require.NoError(t, s.IBC.Relayer.Flush(ctx, s.IBC.RelayerReporter, s.IBC.PathName, counterpartyToAutocctpChannelID), "expected no error relaying MsgRecvPacket & MsgAcknowledgement")

			// We retrieve the height at which the IBC tx has been executed on the chain with AutoCCTP.
			reg := s.Chain.Config().EncodingConfig.InterfaceRegistry
			heightFound := false
			for !heightFound {
				_, err := cosmos.PollForMessage[*clienttypes.MsgUpdateClient](ctx, s.Chain, reg, height, height, nil)
				if err == nil {
					heightFound = true
				} else {
					height += 1
				}
			}

			blockEvents, err = QueryEvents(ctx, val, strconv.Itoa(int(height)))
			require.NoError(t, err, "expected no error querying block events")
			eventFound = false
			for _, event := range blockEvents {
				switch event.Type {
				case "circle.cctp.v1.DepositForBurn":
					for _, attribute := range event.Attributes {
						switch attribute.Key {
						case "amount":
							var actual string
							require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
							require.Equal(t, strconv.Itoa(int(ibcAmt2.Int64())), actual, "expected a different amount in cctp event")
						case "destination_domain":
							require.Equal(t, fmt.Sprintf("%d", s.destinationDomain), attribute.Value, "expected a different destination domain in cctp event")
						case "destination_caller":
							expectedBase64 := ""
							if tC.destinationCaller {
								bz := common.FromHex(destinationCaller)
								bz, err := cli.LeftPadBytes(bz)
								require.NoError(t, err, "expected no error padding destination caller address")
								expectedBase64 = base64.StdEncoding.EncodeToString(bz)
							}
							var actual string
							require.NoError(t, json.Unmarshal([]byte(attribute.Value), &actual))
							require.Equal(t, expectedBase64, actual, "expected a different destination caller in cctp event")
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

			stats := GetAutoCCTPStatsByDestinationDomain(t, ctx, val, fmt.Sprintf("%d", s.destinationDomain))
			require.Equal(t, uint64(1), stats.Accounts, "expected a different number of accounts")
			require.Equal(t, uint64(2), stats.Transfers, "expected a different number of transfers")
			require.Equal(t, ibcAmt1.Add(ibcAmt2).Uint64(), stats.TotalTransferred, "expected a different total transferred")
		})
	}
}
