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

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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
	// Meesage inputs validation
	if msg == nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf("msg to register account cannot be nil")
	}

	accountProperties := msg.GetAccountProperties()
	if err := ms.ValidateAccountProperties(accountProperties); err != nil {
		return nil, types.ErrInvalidAccountProperties.Wrap(err.Error())
	}

	// State transition logic.
	address, err := ms.registerAccount(ctx, accountProperties)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to register the account")
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
