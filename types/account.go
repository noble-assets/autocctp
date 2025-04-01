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
	"bytes"
	"encoding/binary"
	"fmt"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	_ sdk.AccountI             = &Account{}
	_ authtypes.GenesisAccount = &Account{}
)

// GenerateAddress creates returns a new `sdk.AccAddress` derived from the inputs and
// the module name.
func GenerateAddress(accountProperties AccountProperties) sdk.AccAddress {
	rawDestinationDomain := make([]byte, 4)
	binary.BigEndian.PutUint32(rawDestinationDomain, accountProperties.DestinationDomain)

	bz := append(rawDestinationDomain, accountProperties.MintRecipient...)
	bz = append(bz, []byte(accountProperties.FallbackRecipient)...)
	if len(accountProperties.DestinationCaller) != 0 {
		bz = append(bz, accountProperties.DestinationCaller...)
	}

	return address.Derive([]byte(ModuleName), bz)[12:]
}

// AccountProperties specifies properties for configuring cross-chain token transfers
// and fallback mechanisms for AutoCCTP accounts.
type AccountProperties struct {
	DestinationDomain uint32 // Target domain identifier for cross-chain transfers.
	MintRecipient     []byte // Address where tokens will be minted.
	FallbackRecipient string // Backup recipient address for recovering failed transfers.
	DestinationCaller []byte // Optional address that can finalize transfers on the destination chain.
}

func NewAccount(baseAccount *authtypes.BaseAccount, accountProperties AccountProperties) *Account {
	return &Account{
		BaseAccount:       baseAccount,
		DestinationDomain: accountProperties.DestinationDomain,
		MintRecipient:     accountProperties.MintRecipient,
		FallbackRecipient: accountProperties.FallbackRecipient,
		DestinationCaller: accountProperties.DestinationCaller,
	}
}

func (a *Account) Validate() error {
	if err := ValidateMintRecipient(a.MintRecipient); err != nil {
		return ErrInvalidMintRecipient.Wrap(err.Error())
	}

	// NOTE: this is a deprecated approach but there is no other way
	// to check for the valid bech32 without accessing the module keeper.
	if _, err := sdk.AccAddressFromBech32(a.FallbackRecipient); err != nil {
		return ErrInvalidFallbackRecipient.Wrap(err.Error())
	}

	if err := ValidateDestinationCaller(a.DestinationCaller); err != nil {
		return ErrInvalidDestinationCaller.Wrap(err.Error())
	}

	return a.BaseAccount.Validate()
}

//

var _ cryptotypes.PubKey = &PubKey{}

func (pk *PubKey) String() string {
	return fmt.Sprintf("PubKeyAutoCCTP{%X}", pk.Key)
}

func (pk *PubKey) Address() cryptotypes.Address { return pk.Key }

func (pk *PubKey) Bytes() []byte { return pk.Key }

func (*PubKey) VerifySignature(_ []byte, _ []byte) bool {
	panic("PubKey.VerifySignature should never be invoked with AutoCCTP custom key")
}

func (pk *PubKey) Equals(other cryptotypes.PubKey) bool {
	if _, ok := other.(*PubKey); !ok {
		return false
	}

	return bytes.Equal(pk.Bytes(), other.Bytes())
}

func (*PubKey) Type() string { return "autocctp" }
