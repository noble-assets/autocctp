package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

func GenerateAddress(channel string, sender string) sdk.AccAddress {
	bz := []byte(channel + sender)
	return address.Derive([]byte(ModuleName), bz)[12:]
}
