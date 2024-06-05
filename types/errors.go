package types

import "cosmossdk.io/errors"

var ErrMalformedMemo = errors.Register(ModuleName, 1, "malformed memo")
