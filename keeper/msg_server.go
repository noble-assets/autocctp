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

package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"

	"autocctp.dev/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return msgServer{Keeper: keeper}
}

// RegisterAccount is the server entrypoint to register a new AutoCCTP account.
func (ms msgServer) RegisterAccount(ctx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	// Message inputs validation.
	accountProperties := msg.GetAccountProperties()
	if err := ms.ValidateAccountProperties(accountProperties); err != nil {
		return nil, types.ErrInvalidAccountProperties.Wrap(err.Error())
	}

	// State transition logic.
	address, err := ms.registerAccount(ctx, accountProperties)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register the account")
	}

	return &types.MsgRegisterAccountResponse{Address: address}, ms.eventService.EventManager(ctx).Emit(ctx, &types.AccountRegistered{
		Address:           address,
		DestinationDomain: msg.DestinationDomain,
		MintRecipient:     msg.MintRecipient,
		FallbackRecipient: msg.FallbackRecipient,
		DestinationCaller: msg.DestinationCaller,
		Signerlessly:      false,
	})
}

// RegisterAccountSignerlessly is the server entrypoint to register a new AutoCCTP account
// signerlessly.
func (ms msgServer) RegisterAccountSignerlessly(ctx context.Context, msg *types.MsgRegisterAccountSignerlessly) (*types.MsgRegisterAccountSignerlesslyResponse, error) {
	// Message inputs validation
	if msg == nil {
		return nil, errorstypes.ErrInvalidRequest.Wrapf("msg to register account signerlessly cannot be nil")
	}

	accountProperties := msg.GetAccountProperties()
	if err := ms.ValidateAccountProperties(accountProperties); err != nil {
		return nil, types.ErrInvalidAccountProperties.Wrap(err.Error())
	}

	// State transition logic.
	address, err := ms.registerAccount(ctx, accountProperties)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register then account signerlessly")
	}

	return &types.MsgRegisterAccountSignerlesslyResponse{Address: address}, ms.eventService.EventManager(ctx).Emit(ctx, &types.AccountRegistered{
		Address:           address,
		DestinationDomain: msg.DestinationDomain,
		MintRecipient:     msg.MintRecipient,
		FallbackRecipient: msg.FallbackRecipient,
		DestinationCaller: msg.DestinationCaller,
		Signerlessly:      true,
	})
}

// ClearAccount is the server entrypoint to retry the CCTP transfer associated with an AutoCCTP
// account or to clear the account sending funds to the fallback address.
func (ms msgServer) ClearAccount(ctx context.Context, msg *types.MsgClearAccount) (*types.MsgClearAccountResponse, error) {
	// Message inputs validation
	if msg == nil {
		return nil, errorstypes.ErrInvalidRequest.Wrapf("msg to clear an account cannot be nil")
	}

	address, err := ms.accountKeeper.AddressCodec().StringToBytes(msg.Address)
	if err != nil {
		return nil, errorstypes.ErrInvalidAddress.Wrapf("failed to decode autocctp address: %s", err.Error())
	}

	rawAccount := ms.accountKeeper.GetAccount(ctx, address)
	if rawAccount == nil {
		return nil, types.ErrInvalidClearingAccount.Wrapf("account does not exist")
	}
	account, ok := rawAccount.(*types.Account)
	if !ok {
		return nil, types.ErrInvalidClearingAccount.Wrapf("account is not an autocctp account")
	}

	if msg.Fallback && msg.Signer != account.FallbackRecipient {
		return nil, errorstypes.ErrUnauthorized.Wrapf("msg sender must be fallback account: %s != %s", msg.Signer, account.FallbackRecipient)
	}

	mintingToken := ms.ftfKeeper.GetMintingDenom(ctx)
	balance := ms.bankKeeper.GetBalance(ctx, address, mintingToken.Denom)
	if balance.IsZero() {
		return nil, types.ErrInvalidClearingAccount.Wrapf("account does not require clearing")
	}

	// State transition logic.
	err = ms.clearAccount(ctx, account, sdk.NewCoins(balance), msg.Fallback)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to clear the account")
	}

	return &types.MsgClearAccountResponse{}, nil
}
