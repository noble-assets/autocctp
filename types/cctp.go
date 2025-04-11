package types

import (
	"context"

	"github.com/circlefin/noble-cctp/x/cctp/types"
)

func NewCCTPServer(msgServer CCTPMsgServer, queryServer CCTPQueryServer) CCTPService {
	if msgServer == nil {
		panic("CCTP msg server cannot be nil")
	}
	if queryServer == nil {
		panic("CCTP query server cannot be nil")
	}

	return &CCTPServer{
		MsgSever:    msgServer,
		QueryServer: queryServer,
	}
}

var _ CCTPService = CCTPServer{}

type CCTPServer struct {
	MsgSever    CCTPMsgServer
	QueryServer CCTPQueryServer
}

// DepositForBurn implements CCTPService.
func (c CCTPServer) DepositForBurn(ctx context.Context, msg *types.MsgDepositForBurn) (*types.MsgDepositForBurnResponse, error) {
	return c.MsgSever.DepositForBurn(ctx, msg)
}

// DepositForBurnWithCaller implements CCTPService.
func (c CCTPServer) DepositForBurnWithCaller(ctx context.Context, msg *types.MsgDepositForBurnWithCaller) (*types.MsgDepositForBurnWithCallerResponse, error) {
	return c.MsgSever.DepositForBurnWithCaller(ctx, msg)
}

// PerMessageBurnLimit implements CCTPService.
func (c CCTPServer) PerMessageBurnLimit(ctx context.Context, req *types.QueryGetPerMessageBurnLimitRequest) (*types.QueryGetPerMessageBurnLimitResponse, error) {
	return c.QueryServer.PerMessageBurnLimit(ctx, req)
}
