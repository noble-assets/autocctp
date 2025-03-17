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

package mocks

import (
	"context"
	"errors"

	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"autocctp.dev/types"
)

type SendRestrictionFn func(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error)

var (
	_ types.BankKeeper     = BankKeeper{}
	_ cctptypes.BankKeeper = BankKeeper{}
)

type BankKeeper struct {
	// Failing defines if calls to SendCoins return an error response.
	Failing     bool
	Balances    map[string]sdk.Coins
	Restriction SendRestrictionFn
}

func NoOpSendRestrictionFn(_ context.Context, _, toAddr sdk.AccAddress, _ sdk.Coins) (sdk.AccAddress, error) {
	return toAddr, nil
}

// GetBalance implements types.BankKeeper.
func (k BankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	coins := k.Balances[addr.String()]
	amt := coins.AmountOf(denom)

	return sdk.NewCoin(denom, amt)
}

// SendCoinsFromAccountToModule implements types.BankKeeper.
func (k BankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	return nil
}

// SendCoins implements types.BankKeeper.
func (k BankKeeper) SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	if k.Failing {
		return errors.New("errror sending coins")
	}

	fromCoins, found := k.Balances[fromAddr.String()]
	if !found {
		return errors.New("from account not found")
	}

	fromFinalCoins, negativeAmt := fromCoins.SafeSub(amt...)
	if negativeAmt {
		return errors.New("errror during coins deduction")
	}

	toCoins, found := k.Balances[toAddr.String()]
	if !found {
		toCoins = sdk.Coins{}
	}

	toFinalCoins := toCoins.Add(amt...)

	k.Balances[fromAddr.String()] = fromFinalCoins
	k.Balances[toAddr.String()] = toFinalCoins

	return nil
}
