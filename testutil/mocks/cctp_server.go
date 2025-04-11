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

	"cosmossdk.io/math"

	"autocctp.dev/types"
)

var _ types.CCTPService = CCTPServer{}

type MockCounter struct {
	// NumDepositForBurn keep track of the number of times the associated method is called.
	NumDepositForBurn int
	// NumDepositForBurnWithCaller keep track of the number of times the associated method is called.
	NumDepositForBurnWithCaller int
}

type CCTPServer struct {
	// Failing defines if calls to the CCTPServer return an error response.
	Failing           bool
	MaxTransferAmount int64
	// MockCounter is used to check if the proper method has been called.
	MockCounter *MockCounter
}

func (c CCTPServer) DepositForBurn(_ context.Context, msg *cctptypes.MsgDepositForBurn) (*cctptypes.MsgDepositForBurnResponse, error) {
	if c.Failing {
		return nil, errors.New("error calling deposit for burn api")
	}

	c.MockCounter.NumDepositForBurn += 1

	return nil, nil
}

func (c CCTPServer) DepositForBurnWithCaller(_ context.Context, msg *cctptypes.MsgDepositForBurnWithCaller) (*cctptypes.MsgDepositForBurnWithCallerResponse, error) {
	if c.Failing {
		return nil, errors.New("error calling deposit for burn with caller api")
	}

	c.MockCounter.NumDepositForBurnWithCaller += 1

	return nil, nil
}

func (c CCTPServer) PerMessageBurnLimit(context.Context, *cctptypes.QueryGetPerMessageBurnLimitRequest) (*cctptypes.QueryGetPerMessageBurnLimitResponse, error) {
	return &cctptypes.QueryGetPerMessageBurnLimitResponse{
		BurnLimit: cctptypes.PerMessageBurnLimit{
			Denom:  "uusdc",
			Amount: math.NewInt(c.MaxTransferAmount),
		},
	}, nil
}
