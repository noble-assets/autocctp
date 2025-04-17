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

package autocctp

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"autocctp.dev/keeper"
	"autocctp.dev/types"
)

// SigVerificationGasConsumer is a wrapper around the default function provided by the
// Cosmos SDK that supports AutoCCTP account public keys. If the public key is an AutoCCTP
// key, skip gas consumption.
func SigVerificationGasConsumer(meter storetypes.GasMeter, sig signing.SignatureV2, params authtypes.Params) error {
	switch sig.PubKey.(type) {
	case *types.PubKey:
		return nil
	default:
		return ante.DefaultSigVerificationGasConsumer(meter, sig, params)
	}
}

//

var _ sdk.AnteDecorator = SigVerificationDecorator{}

// SigVerificationDecorator is a custom ante handler used to verify signerless registration messages for AutoCCTP accounts.
type SigVerificationDecorator struct {
	autocctp   keeper.Keeper
	underlying sdk.AnteDecorator
}

// NewSigVerificationDecorator returns the AutoCCTP signature verification ante handler which
// wraps another ante handler. The custom signature verification allows to broadcast a
// signerless tx to create an AutoCCTP account.
func NewSigVerificationDecorator(
	autocctp keeper.Keeper,
	underlying sdk.AnteDecorator,
) SigVerificationDecorator {
	if underlying == nil {
		panic("underlying ante decorator cannot be nil")
	}

	return SigVerificationDecorator{
		autocctp:   autocctp,
		underlying: underlying,
	}
}

// AnteHandle check if the transaction contains a single message of type `MsgRegisterAccountSignerlessly` and
// if the signature verification has to be skipped.
func (d SigVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if msgs := tx.GetMsgs(); len(msgs) == 1 {
		msg, ok := msgs[0].(*types.MsgRegisterAccountSignerlessly)
		if !ok {
			return d.underlying.AnteHandle(ctx, tx, simulate, next)
		}

		address := types.GenerateAddress(msg.GetAccountProperties())

		// Check if the balance of the wannabe AutoCCTP account is not zero to prevent spam
		// attacks of AutoCCTP accounts.
		mintToken := d.ftf.GetMintingDenom(ctx)
		balance := d.bank.GetBalance(ctx, address, mintToken.Denom)
		if balance.Amount.LT(types.GetMinimumTransferAmount()) || msg.Signer != address.String() {
			return d.underlying.AnteHandle(ctx, tx, simulate, next)
		}

		return next(ctx, tx, simulate)
	}

	return d.underlying.AnteHandle(ctx, tx, simulate, next)
}
