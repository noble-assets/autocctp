// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/circlefin/noble-cctp/x/cctp/types (interfaces: MsgServer)
//
// Generated by this command:
//
//	mockgen -package=mock -destination=./test/mock/cctp_msgserver.go github.com/circlefin/noble-cctp/x/cctp/types MsgServer
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "github.com/circlefin/noble-cctp/x/cctp/types"
	gomock "go.uber.org/mock/gomock"
)

// MockMsgServer is a mock of MsgServer interface.
type MockMsgServer struct {
	ctrl     *gomock.Controller
	recorder *MockMsgServerMockRecorder
	isgomock struct{}
}

// MockMsgServerMockRecorder is the mock recorder for MockMsgServer.
type MockMsgServerMockRecorder struct {
	mock *MockMsgServer
}

// NewMockMsgServer creates a new mock instance.
func NewMockMsgServer(ctrl *gomock.Controller) *MockMsgServer {
	mock := &MockMsgServer{ctrl: ctrl}
	mock.recorder = &MockMsgServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMsgServer) EXPECT() *MockMsgServerMockRecorder {
	return m.recorder
}

// AcceptOwner mocks base method.
func (m *MockMsgServer) AcceptOwner(arg0 context.Context, arg1 *types.MsgAcceptOwner) (*types.MsgAcceptOwnerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptOwner", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgAcceptOwnerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AcceptOwner indicates an expected call of AcceptOwner.
func (mr *MockMsgServerMockRecorder) AcceptOwner(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptOwner", reflect.TypeOf((*MockMsgServer)(nil).AcceptOwner), arg0, arg1)
}

// AddRemoteTokenMessenger mocks base method.
func (m *MockMsgServer) AddRemoteTokenMessenger(arg0 context.Context, arg1 *types.MsgAddRemoteTokenMessenger) (*types.MsgAddRemoteTokenMessengerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRemoteTokenMessenger", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgAddRemoteTokenMessengerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRemoteTokenMessenger indicates an expected call of AddRemoteTokenMessenger.
func (mr *MockMsgServerMockRecorder) AddRemoteTokenMessenger(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRemoteTokenMessenger", reflect.TypeOf((*MockMsgServer)(nil).AddRemoteTokenMessenger), arg0, arg1)
}

// DepositForBurn mocks base method.
func (m *MockMsgServer) DepositForBurn(arg0 context.Context, arg1 *types.MsgDepositForBurn) (*types.MsgDepositForBurnResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepositForBurn", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgDepositForBurnResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DepositForBurn indicates an expected call of DepositForBurn.
func (mr *MockMsgServerMockRecorder) DepositForBurn(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositForBurn", reflect.TypeOf((*MockMsgServer)(nil).DepositForBurn), arg0, arg1)
}

// DepositForBurnWithCaller mocks base method.
func (m *MockMsgServer) DepositForBurnWithCaller(arg0 context.Context, arg1 *types.MsgDepositForBurnWithCaller) (*types.MsgDepositForBurnWithCallerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepositForBurnWithCaller", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgDepositForBurnWithCallerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DepositForBurnWithCaller indicates an expected call of DepositForBurnWithCaller.
func (mr *MockMsgServerMockRecorder) DepositForBurnWithCaller(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositForBurnWithCaller", reflect.TypeOf((*MockMsgServer)(nil).DepositForBurnWithCaller), arg0, arg1)
}

// DisableAttester mocks base method.
func (m *MockMsgServer) DisableAttester(arg0 context.Context, arg1 *types.MsgDisableAttester) (*types.MsgDisableAttesterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableAttester", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgDisableAttesterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DisableAttester indicates an expected call of DisableAttester.
func (mr *MockMsgServerMockRecorder) DisableAttester(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableAttester", reflect.TypeOf((*MockMsgServer)(nil).DisableAttester), arg0, arg1)
}

// EnableAttester mocks base method.
func (m *MockMsgServer) EnableAttester(arg0 context.Context, arg1 *types.MsgEnableAttester) (*types.MsgEnableAttesterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableAttester", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgEnableAttesterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableAttester indicates an expected call of EnableAttester.
func (mr *MockMsgServerMockRecorder) EnableAttester(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableAttester", reflect.TypeOf((*MockMsgServer)(nil).EnableAttester), arg0, arg1)
}

// LinkTokenPair mocks base method.
func (m *MockMsgServer) LinkTokenPair(arg0 context.Context, arg1 *types.MsgLinkTokenPair) (*types.MsgLinkTokenPairResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LinkTokenPair", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgLinkTokenPairResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LinkTokenPair indicates an expected call of LinkTokenPair.
func (mr *MockMsgServerMockRecorder) LinkTokenPair(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LinkTokenPair", reflect.TypeOf((*MockMsgServer)(nil).LinkTokenPair), arg0, arg1)
}

// PauseBurningAndMinting mocks base method.
func (m *MockMsgServer) PauseBurningAndMinting(arg0 context.Context, arg1 *types.MsgPauseBurningAndMinting) (*types.MsgPauseBurningAndMintingResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseBurningAndMinting", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgPauseBurningAndMintingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PauseBurningAndMinting indicates an expected call of PauseBurningAndMinting.
func (mr *MockMsgServerMockRecorder) PauseBurningAndMinting(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseBurningAndMinting", reflect.TypeOf((*MockMsgServer)(nil).PauseBurningAndMinting), arg0, arg1)
}

// PauseSendingAndReceivingMessages mocks base method.
func (m *MockMsgServer) PauseSendingAndReceivingMessages(arg0 context.Context, arg1 *types.MsgPauseSendingAndReceivingMessages) (*types.MsgPauseSendingAndReceivingMessagesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PauseSendingAndReceivingMessages", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgPauseSendingAndReceivingMessagesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PauseSendingAndReceivingMessages indicates an expected call of PauseSendingAndReceivingMessages.
func (mr *MockMsgServerMockRecorder) PauseSendingAndReceivingMessages(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PauseSendingAndReceivingMessages", reflect.TypeOf((*MockMsgServer)(nil).PauseSendingAndReceivingMessages), arg0, arg1)
}

// ReceiveMessage mocks base method.
func (m *MockMsgServer) ReceiveMessage(arg0 context.Context, arg1 *types.MsgReceiveMessage) (*types.MsgReceiveMessageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveMessage", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgReceiveMessageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReceiveMessage indicates an expected call of ReceiveMessage.
func (mr *MockMsgServerMockRecorder) ReceiveMessage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveMessage", reflect.TypeOf((*MockMsgServer)(nil).ReceiveMessage), arg0, arg1)
}

// RemoveRemoteTokenMessenger mocks base method.
func (m *MockMsgServer) RemoveRemoteTokenMessenger(arg0 context.Context, arg1 *types.MsgRemoveRemoteTokenMessenger) (*types.MsgRemoveRemoteTokenMessengerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRemoteTokenMessenger", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgRemoveRemoteTokenMessengerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveRemoteTokenMessenger indicates an expected call of RemoveRemoteTokenMessenger.
func (mr *MockMsgServerMockRecorder) RemoveRemoteTokenMessenger(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRemoteTokenMessenger", reflect.TypeOf((*MockMsgServer)(nil).RemoveRemoteTokenMessenger), arg0, arg1)
}

// ReplaceDepositForBurn mocks base method.
func (m *MockMsgServer) ReplaceDepositForBurn(arg0 context.Context, arg1 *types.MsgReplaceDepositForBurn) (*types.MsgReplaceDepositForBurnResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceDepositForBurn", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgReplaceDepositForBurnResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceDepositForBurn indicates an expected call of ReplaceDepositForBurn.
func (mr *MockMsgServerMockRecorder) ReplaceDepositForBurn(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceDepositForBurn", reflect.TypeOf((*MockMsgServer)(nil).ReplaceDepositForBurn), arg0, arg1)
}

// ReplaceMessage mocks base method.
func (m *MockMsgServer) ReplaceMessage(arg0 context.Context, arg1 *types.MsgReplaceMessage) (*types.MsgReplaceMessageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceMessage", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgReplaceMessageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReplaceMessage indicates an expected call of ReplaceMessage.
func (mr *MockMsgServerMockRecorder) ReplaceMessage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceMessage", reflect.TypeOf((*MockMsgServer)(nil).ReplaceMessage), arg0, arg1)
}

// SendMessage mocks base method.
func (m *MockMsgServer) SendMessage(arg0 context.Context, arg1 *types.MsgSendMessage) (*types.MsgSendMessageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgSendMessageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockMsgServerMockRecorder) SendMessage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockMsgServer)(nil).SendMessage), arg0, arg1)
}

// SendMessageWithCaller mocks base method.
func (m *MockMsgServer) SendMessageWithCaller(arg0 context.Context, arg1 *types.MsgSendMessageWithCaller) (*types.MsgSendMessageWithCallerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessageWithCaller", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgSendMessageWithCallerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessageWithCaller indicates an expected call of SendMessageWithCaller.
func (mr *MockMsgServerMockRecorder) SendMessageWithCaller(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageWithCaller", reflect.TypeOf((*MockMsgServer)(nil).SendMessageWithCaller), arg0, arg1)
}

// SetMaxBurnAmountPerMessage mocks base method.
func (m *MockMsgServer) SetMaxBurnAmountPerMessage(arg0 context.Context, arg1 *types.MsgSetMaxBurnAmountPerMessage) (*types.MsgSetMaxBurnAmountPerMessageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMaxBurnAmountPerMessage", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgSetMaxBurnAmountPerMessageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetMaxBurnAmountPerMessage indicates an expected call of SetMaxBurnAmountPerMessage.
func (mr *MockMsgServerMockRecorder) SetMaxBurnAmountPerMessage(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMaxBurnAmountPerMessage", reflect.TypeOf((*MockMsgServer)(nil).SetMaxBurnAmountPerMessage), arg0, arg1)
}

// UnlinkTokenPair mocks base method.
func (m *MockMsgServer) UnlinkTokenPair(arg0 context.Context, arg1 *types.MsgUnlinkTokenPair) (*types.MsgUnlinkTokenPairResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlinkTokenPair", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUnlinkTokenPairResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnlinkTokenPair indicates an expected call of UnlinkTokenPair.
func (mr *MockMsgServerMockRecorder) UnlinkTokenPair(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlinkTokenPair", reflect.TypeOf((*MockMsgServer)(nil).UnlinkTokenPair), arg0, arg1)
}

// UnpauseBurningAndMinting mocks base method.
func (m *MockMsgServer) UnpauseBurningAndMinting(arg0 context.Context, arg1 *types.MsgUnpauseBurningAndMinting) (*types.MsgUnpauseBurningAndMintingResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpauseBurningAndMinting", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUnpauseBurningAndMintingResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnpauseBurningAndMinting indicates an expected call of UnpauseBurningAndMinting.
func (mr *MockMsgServerMockRecorder) UnpauseBurningAndMinting(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpauseBurningAndMinting", reflect.TypeOf((*MockMsgServer)(nil).UnpauseBurningAndMinting), arg0, arg1)
}

// UnpauseSendingAndReceivingMessages mocks base method.
func (m *MockMsgServer) UnpauseSendingAndReceivingMessages(arg0 context.Context, arg1 *types.MsgUnpauseSendingAndReceivingMessages) (*types.MsgUnpauseSendingAndReceivingMessagesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpauseSendingAndReceivingMessages", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUnpauseSendingAndReceivingMessagesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnpauseSendingAndReceivingMessages indicates an expected call of UnpauseSendingAndReceivingMessages.
func (mr *MockMsgServerMockRecorder) UnpauseSendingAndReceivingMessages(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpauseSendingAndReceivingMessages", reflect.TypeOf((*MockMsgServer)(nil).UnpauseSendingAndReceivingMessages), arg0, arg1)
}

// UpdateAttesterManager mocks base method.
func (m *MockMsgServer) UpdateAttesterManager(arg0 context.Context, arg1 *types.MsgUpdateAttesterManager) (*types.MsgUpdateAttesterManagerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAttesterManager", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUpdateAttesterManagerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAttesterManager indicates an expected call of UpdateAttesterManager.
func (mr *MockMsgServerMockRecorder) UpdateAttesterManager(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAttesterManager", reflect.TypeOf((*MockMsgServer)(nil).UpdateAttesterManager), arg0, arg1)
}

// UpdateMaxMessageBodySize mocks base method.
func (m *MockMsgServer) UpdateMaxMessageBodySize(arg0 context.Context, arg1 *types.MsgUpdateMaxMessageBodySize) (*types.MsgUpdateMaxMessageBodySizeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMaxMessageBodySize", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUpdateMaxMessageBodySizeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateMaxMessageBodySize indicates an expected call of UpdateMaxMessageBodySize.
func (mr *MockMsgServerMockRecorder) UpdateMaxMessageBodySize(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMaxMessageBodySize", reflect.TypeOf((*MockMsgServer)(nil).UpdateMaxMessageBodySize), arg0, arg1)
}

// UpdateOwner mocks base method.
func (m *MockMsgServer) UpdateOwner(arg0 context.Context, arg1 *types.MsgUpdateOwner) (*types.MsgUpdateOwnerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOwner", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUpdateOwnerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOwner indicates an expected call of UpdateOwner.
func (mr *MockMsgServerMockRecorder) UpdateOwner(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOwner", reflect.TypeOf((*MockMsgServer)(nil).UpdateOwner), arg0, arg1)
}

// UpdatePauser mocks base method.
func (m *MockMsgServer) UpdatePauser(arg0 context.Context, arg1 *types.MsgUpdatePauser) (*types.MsgUpdatePauserResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePauser", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUpdatePauserResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePauser indicates an expected call of UpdatePauser.
func (mr *MockMsgServerMockRecorder) UpdatePauser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePauser", reflect.TypeOf((*MockMsgServer)(nil).UpdatePauser), arg0, arg1)
}

// UpdateSignatureThreshold mocks base method.
func (m *MockMsgServer) UpdateSignatureThreshold(arg0 context.Context, arg1 *types.MsgUpdateSignatureThreshold) (*types.MsgUpdateSignatureThresholdResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSignatureThreshold", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUpdateSignatureThresholdResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSignatureThreshold indicates an expected call of UpdateSignatureThreshold.
func (mr *MockMsgServerMockRecorder) UpdateSignatureThreshold(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSignatureThreshold", reflect.TypeOf((*MockMsgServer)(nil).UpdateSignatureThreshold), arg0, arg1)
}

// UpdateTokenController mocks base method.
func (m *MockMsgServer) UpdateTokenController(arg0 context.Context, arg1 *types.MsgUpdateTokenController) (*types.MsgUpdateTokenControllerResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTokenController", arg0, arg1)
	ret0, _ := ret[0].(*types.MsgUpdateTokenControllerResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTokenController indicates an expected call of UpdateTokenController.
func (mr *MockMsgServerMockRecorder) UpdateTokenController(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTokenController", reflect.TypeOf((*MockMsgServer)(nil).UpdateTokenController), arg0, arg1)
}