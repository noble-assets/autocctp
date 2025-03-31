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
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/event"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"autocctp.dev/types"
)

type Keeper struct {
	logger       log.Logger
	eventService event.Service

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	ftfKeeper     types.FiatTokenfactoryKeeper

	cctpServer types.CCTPServer

	// NumOfAccounts keeps track of the number of accounts registered per destination domain.
	NumOfAccounts collections.Map[uint32, uint64]
	// NumOfTransfers keeps track of the number of transfers executed per destination domain.
	NumOfTransfers collections.Map[uint32, uint64]
	// TotalTransferred keeps track of the total value transferred per destination domain.
	TotalTransferred collections.Map[uint32, uint64]

	// PendingTransfers is a transient map that keeps track of the pending transfers for the current block.
	PendingTransfers collections.Map[string, types.Account]
}

func NewKeeper(
	cdc codec.Codec,
	logger log.Logger,
	storeService store.KVStoreService,
	transientService store.TransientStoreService,
	eventService event.Service,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	cctpServer types.CCTPServer,
	ftfKeeper types.FiatTokenfactoryKeeper,
) *Keeper {
	builder := collections.NewSchemaBuilder(storeService)
	transientBuilder := collections.NewSchemaBuilderFromAccessor(transientService.OpenTransientStore)

	keeper := &Keeper{
		logger:       logger.With("module", types.ModuleName),
		eventService: eventService,

		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		ftfKeeper:     fiattokenfactory,

		cctpServer: cctpServer,

		NumOfAccounts:    collections.NewMap(builder, types.NumOfAccountsPrefix, "num_of_accounts", collections.Uint32Key, collections.Uint64Value),
		NumOfTransfers:   collections.NewMap(builder, types.NumOfTransfersPrefix, "num_of_transfers", collections.Uint32Key, collections.Uint64Value),
		TotalTransferred: collections.NewMap(builder, types.TotalTransferredPrefix, "total_transferred", collections.Uint32Key, collections.Uint64Value),

		PendingTransfers: collections.NewMap(transientBuilder, types.PendingTransfersPrefix, "pending_transfers", collections.StringKey, codec.CollValue[types.Account](cdc)),
	}

	if _, err := builder.Build(); err != nil {
		panic(err)
	}

	return keeper
}

// CCTPServer is a method used to set the CCTP server into the AutoCCTP keeper. This method is
// required because is the only way to inject a server dependency into the keeper using deepinject.
func (k *Keeper) SetCCTPServer(cctpServer types.CCTPServer) {
	k.cctpServer = cctpServer
}

// ValidateAccountProperties returns an error if any account properties is not valid.
func (k *Keeper) ValidateAccountProperties(accountProperties types.AccountProperties) error {
	if err := types.ValidateMintRecipient(accountProperties.MintRecipient); err != nil {
		return types.ErrInvalidMintRecipient.Wrap(err.Error())
	}

	_, err := k.accountKeeper.AddressCodec().StringToBytes(accountProperties.FallbackRecipient)
	if err != nil {
		return types.ErrInvalidFallbackRecipient.Wrap(err.Error())
	}

	if err := types.ValidateDestinationCaller(accountProperties.DestinationCaller); err != nil {
		return types.ErrInvalidDestinationCaller.Wrap(err.Error())
	}

	return nil
}

// SendRestrictionFn checks every transfer executed on the Noble chain to see if
// the recipient is an AutoCCTP account, allowing us to mark them for clearing.
func (k *Keeper) SendRestrictionFn(ctx context.Context, _, toAddr sdk.AccAddress, coins sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	rawAccount := k.accountKeeper.GetAccount(ctx, toAddr)
	if rawAccount == nil {
		return toAddr, nil
	}

	account, ok := rawAccount.(*types.Account)
	if !ok {
		return toAddr, nil
	}

	mintingDenom := k.ftfKeeper.GetMintingDenom(ctx).Denom
	if coins.AmountOf(mintingDenom).IsPositive() {
		err = k.PendingTransfers.Set(ctx, account.Address, *account)
	}

	return toAddr, err
}

// registerAccount handles the registration of a custom account given the AutoCCTP properties.
//
// CONTRACT: the function assumes properties have already been validated.
func (k Keeper) registerAccount(ctx context.Context, accountProperties types.AccountProperties) (string, error) {
	address := types.GenerateAddress(accountProperties)

	if k.accountKeeper.HasAccount(ctx, address) {
		rawAccount := k.accountKeeper.GetAccount(ctx, address)

		if err := types.ValidateExistingAccount(rawAccount, address); err != nil {
			return "", fmt.Errorf("error validating existing account: %w", err)
		}

		switch account := rawAccount.(type) {
		case *authtypes.BaseAccount:
			rawAccount = types.NewAccount(account, accountProperties)
			k.accountKeeper.SetAccount(ctx, rawAccount)

			k.IncrementNumOfAccounts(ctx, accountProperties.DestinationDomain)
		case *types.Account:
			return "", errors.New("account has already been registered")
		default:
			return "", fmt.Errorf("unsupported account type: %T", rawAccount)
		}

		mintingToken := k.ftfKeeper.GetMintingDenom(ctx)
		if !k.bankKeeper.GetBalance(ctx, address, mintingToken.Denom).IsZero() {
			account, _ := rawAccount.(*types.Account)
			if err := k.PendingTransfers.Set(ctx, address.String(), *account); err != nil {
				k.logger.Error("error registering pending transfer for address %s", address.String())
			}
		}

		return address.String(), nil
	}

	base := k.accountKeeper.NewAccountWithAddress(ctx, address)
	baseAccount := authtypes.NewBaseAccount(base.GetAddress(), base.GetPubKey(), base.GetAccountNumber(), base.GetSequence())

	account := types.NewAccount(baseAccount, accountProperties)

	k.accountKeeper.SetAccount(ctx, account)
	k.IncrementNumOfAccounts(ctx, accountProperties.DestinationDomain)

	return address.String(), nil
}
