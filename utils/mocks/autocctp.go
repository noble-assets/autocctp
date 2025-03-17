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
	"fmt"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/stretchr/testify/assert"

	"autocctp.dev/keeper"
	"autocctp.dev/types"
)

type Mocks struct {
	AccountKeeper *AccountKeeper
	BankKeeper    *BankKeeper
	FTFKeeper     *FTFKeeper
	CCTPServer    *CCTPServer
}

// AutoCCTPKeeper returns the AutoCCTP keeper with all dependencies mocked and a context.
func AutoCCTPKeeper(t testing.TB) (*Mocks, *keeper.Keeper, sdk.Context) {
	ak := AccountKeeper{
		Accounts: make(map[string]sdk.AccountI),
	}
	bk := BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: NoOpSendRestrictionFn,
	}
	cctps := CCTPServer{
		MockCounter: &MockCounter{},
		Failing:     false,
	}

	mocks := Mocks{
		AccountKeeper: &ak,
		BankKeeper:    &bk,
		FTFKeeper:     &FTFKeeper{},
		CCTPServer:    &cctps,
	}

	k, ctx := autoCCTPKeeperWithMocks(t, &mocks)

	return &mocks, k, ctx
}

// autoCCTPKeeperWithMocks returns an instance of the AutoCCTP keeper and creates all the store dependencies required.
func autoCCTPKeeperWithMocks(t testing.TB, m *Mocks) (*keeper.Keeper, sdk.Context) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	tkey := storetypes.NewTransientStoreKey(fmt.Sprintf("transient_%s", types.ModuleName))
	wrapper := testutil.DefaultContextWithDB(t, key, tkey)

	cfg := MakeTestEncodingConfig("noble")
	types.RegisterInterfaces(cfg.InterfaceRegistry)

	k := keeper.NewKeeper(
		cfg.Codec,
		log.NewNopLogger(),
		runtime.NewKVStoreService(key),
		runtime.NewTransientStoreService(tkey),
		runtime.ProvideEventService(),
		m.AccountKeeper,
		m.BankKeeper,
		m.FTFKeeper,
		m.CCTPServer,
	)

	k.InitGenesis(wrapper.Ctx, *types.DefaultGenesisState())
	return k, wrapper.Ctx
}

func ResetTest(t *testing.T, ctx context.Context, k *keeper.Keeper, m *Mocks) {
	m.CCTPServer.MockCounter.NumDepositForBurn = 0
	m.CCTPServer.MockCounter.NumDepositForBurnWithCaller = 0

	m.BankKeeper.Balances = make(map[string]sdk.Coins)

	m.AccountKeeper.Accounts = make(map[string]sdk.AccountI)

	err := k.PendingTransfers.Clear(ctx, nil)
	assert.NoError(t, err)

	err = k.NumOfAccounts.Clear(ctx, nil)
	assert.NoError(t, err)

	err = k.NumOfTransfers.Clear(ctx, nil)
	assert.NoError(t, err)

	err = k.TotalTransferred.Clear(ctx, nil)
	assert.NoError(t, err)
}

// MakeTestEncodingConfig is a modified testutil.MakeTestEncodingConfig that
// sets a custom Bech32 prefix in the interface registry.
func MakeTestEncodingConfig(prefix string, modules ...module.AppModuleBasic) moduletestutil.TestEncodingConfig {
	aminoCodec := codec.NewLegacyAmino()
	interfaceRegistry := codectestutil.CodecOptions{
		AccAddressPrefix: prefix,
	}.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)

	encCfg := moduletestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          tx.NewTxConfig(codec, tx.DefaultSignModes),
		Amino:             aminoCodec,
	}

	mb := module.NewBasicManager(modules...)

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	mb.RegisterLegacyAminoCodec(encCfg.Amino)
	mb.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}
