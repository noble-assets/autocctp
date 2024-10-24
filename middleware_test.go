package autocctp_test

import (
	"encoding/json"
	"fmt"
	"testing"

	autocctp "autocctp.dev"
	mock "autocctp.dev/test/mock"
	"autocctp.dev/types"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	cctpAmount                 = "99"
	feeRecipient               = "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue"
	autocctpDepositForBurnMemo = types.Memo{
		DepositForBurn: &types.DepositForBurn{
			DestinationDomain: 0,
			MintRecipient:     []byte("mintRecipient"),
			Amount:            &cctpAmount,
			FeeRecipient:      &feeRecipient,
		},
	}
	autocctpDepositForBurnWithCallerMemo = types.Memo{
		DepositForBurnWithCaller: &types.DepositForBurnWithCaller{
			DepositForBurn: types.DepositForBurn{
				DestinationDomain: 0,
				MintRecipient:     []byte("mintRecipient"),
				Amount:            &cctpAmount,
				FeeRecipient:      &feeRecipient,
			},
			DestinationCaller: []byte("destinationCaller"),
		},
	}
	transferPacket = transfertypes.FungibleTokenPacketData{
		Denom:    "transfer/channel-0/uusdc",
		Amount:   "100",
		Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
		Receiver: authtypes.NewModuleAddress("autocctp").String(),
	}
	packet = channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
	}
)

func TestMiddleware(t *testing.T) {
	testcases := []struct {
		name          string
		getPacket     func() channeltypes.Packet
		expectSuccess bool
	}{
		{
			name: "Success: Packet not transfer type, so move on",
			getPacket: func() channeltypes.Packet {
				return channeltypes.Packet{}
			},
			expectSuccess: true,
		},
		{
			name: "Success: Receiver address is not the autocctp module address, so do not perform autocctp transfer",
			getPacket: func() channeltypes.Packet {
				transferPacket.Receiver = "cosmos1vzxkv3lxccnttr9rs0002s93sgw72h7ghukuhs"
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: true,
		},
		{
			name: "Success: Sender chain is source chain, so move on",
			getPacket: func() channeltypes.Packet {
				transferPacket.Denom = "testDenom"
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: true,
		},
		{
			name: "Fail: types.Memo exists but is empty",
			getPacket: func() channeltypes.Packet {
				memo := types.Memo{}
				memobz, _ := json.Marshal(memo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurn - forward amount is greater than transferred amount",
			getPacket: func() channeltypes.Packet {
				overAmount := "200"
				autocctpDepositForBurnMemo.DepositForBurn.Amount = &overAmount
				memobz, _ := json.Marshal(autocctpDepositForBurnMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurnWithCaller - foward amount is greater than transferred amount",
			getPacket: func() channeltypes.Packet {
				overAmount := "200"
				autocctpDepositForBurnWithCallerMemo.DepositForBurnWithCaller.Amount = &overAmount
				memobz, _ := json.Marshal(autocctpDepositForBurnWithCallerMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurn - fee recipient is nil",
			getPacket: func() channeltypes.Packet {
				autocctpDepositForBurnMemo.DepositForBurn.FeeRecipient = nil
				memobz, _ := json.Marshal(autocctpDepositForBurnMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurnWithCaller - fee recipient is nil",
			getPacket: func() channeltypes.Packet {
				autocctpDepositForBurnWithCallerMemo.DepositForBurnWithCaller.FeeRecipient = nil
				memobz, _ := json.Marshal(autocctpDepositForBurnWithCallerMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurn - fee recipient is invalid bech32",
			getPacket: func() channeltypes.Packet {
				feeRecipient := "invalidbech32"
				autocctpDepositForBurnMemo.DepositForBurn.FeeRecipient = &feeRecipient
				memobz, _ := json.Marshal(autocctpDepositForBurnMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				return channeltypes.Packet{
					SourcePort:         "transfer",
					SourceChannel:      "channel-0",
					DestinationPort:    "transfer",
					DestinationChannel: "channel-1",
					Data:               transferData,
				}
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurnWithCaller - fee recipient is invalid bech32",
			getPacket: func() channeltypes.Packet {
				feeRecipient := "invalidbech32"
				autocctpDepositForBurnWithCallerMemo.DepositForBurnWithCaller.FeeRecipient = &feeRecipient
				memobz, _ := json.Marshal(autocctpDepositForBurnWithCallerMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				return channeltypes.Packet{
					SourcePort:         "transfer",
					SourceChannel:      "channel-0",
					DestinationPort:    "transfer",
					DestinationChannel: "channel-1",
					Data:               transferData,
				}
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurn - forward amount is invalid",
			getPacket: func() channeltypes.Packet {
				cctpAmount := "👻"
				autocctpDepositForBurnMemo.DepositForBurn.Amount = &cctpAmount
				memobz, _ := json.Marshal(autocctpDepositForBurnMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
		{
			name: "Fail: DepositForBurnWithCaller - forward amount is invalid",
			getPacket: func() channeltypes.Packet {
				cctpAmount := "👻"
				autocctpDepositForBurnWithCallerMemo.DepositForBurnWithCaller.Amount = &cctpAmount
				memobz, _ := json.Marshal(autocctpDepositForBurnWithCallerMemo)
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "transfer/channel-0/uusdc",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
					Memo:     string(memobz),
				}
				transferData, _ := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				packet.Data = transferData
				return packet
			},
			expectSuccess: false,
		},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("Case: %s", tc.name), func(t *testing.T) {
			logger := log.NewTestLogger(t)
			ctx := sdk.NewContext(store.NewCommitMultiStore(dbm.NewMemDB(), logger, metrics.NewNoOpMetrics()), tmproto.Header{}, false, logger)
			// Setting up mocks
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			packet := tc.getPacket()
			bankKeeper := mock.NewMockBankKeeper(ctl)
			ibcModule := mock.NewMockIBCModule(ctl)
			cctpServer := mock.NewMockMsgServer(ctl)
			ibcModule.EXPECT().OnRecvPacket(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(channeltypes.NewResultAcknowledgement([]byte("")))
			bankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil).AnyTimes()
			middleware := autocctp.NewMiddleware(ibcModule, bankKeeper, cctpServer)

			ack := middleware.OnRecvPacket(ctx, packet, nil)

			if tc.expectSuccess {
				require.True(t, ack.Success())
			} else {
				require.False(t, ack.Success())
			}
		})
	}
}

func TestMiddleware_DepositForBurn_Success(t *testing.T) {
	logger := log.NewTestLogger(t)
	ctx := sdk.NewContext(store.NewCommitMultiStore(dbm.NewMemDB(), logger, metrics.NewNoOpMetrics()), tmproto.Header{}, false, logger)
	// Setting up mocks
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	bankKeeper := mock.NewMockBankKeeper(ctl)
	ibcModule := mock.NewMockIBCModule(ctl)
	cctpServer := mock.NewMockMsgServer(ctl)
	ibcModule.EXPECT().OnRecvPacket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(channeltypes.NewResultAcknowledgement([]byte(""))).Times(1)
	bankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).Times(2)
	cctpServer.EXPECT().DepositForBurn(gomock.Any(), gomock.Any()).Return(&cctptypes.MsgDepositForBurnResponse{
		Nonce: 10,
	}, nil).Times(1)
	middleware := autocctp.NewMiddleware(ibcModule, bankKeeper, cctpServer)

	memo := types.Memo{
		DepositForBurn: &types.DepositForBurn{
			DestinationDomain: 0,
			MintRecipient:     []byte("mintRecipient"),
			Amount:            &cctpAmount,
			FeeRecipient:      &feeRecipient,
		},
	}
	memobz, _ := json.Marshal(memo)
	testTransferPacket := transfertypes.FungibleTokenPacketData{
		Denom:    "transfer/channel-0/uusdc",
		Amount:   "100",
		Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
		Receiver: authtypes.NewModuleAddress("autocctp").String(),
		Memo:     string(memobz),
	}
	testTransferData, _ := transfertypes.ModuleCdc.MarshalJSON(&testTransferPacket)
	testPacket := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               testTransferData,
	}

	res := middleware.OnRecvPacket(ctx, testPacket, nil)

	require.True(t, res.Success())
	var ack channeltypes.Acknowledgement
	err := channeltypes.SubModuleCdc.UnmarshalJSON(res.Acknowledgement(), &ack)
	require.NoError(t, err)
	require.Equal(t, "{\"nonce\":10}", string(ack.GetResult()))
}

func TestMiddleware_DepositForBurnWithCaller_Success(t *testing.T) {
	logger := log.NewTestLogger(t)
	ctx := sdk.NewContext(store.NewCommitMultiStore(dbm.NewMemDB(), logger, metrics.NewNoOpMetrics()), tmproto.Header{}, false, logger)
	// Setting up mocks
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	bankKeeper := mock.NewMockBankKeeper(ctl)
	ibcModule := mock.NewMockIBCModule(ctl)
	cctpServer := mock.NewMockMsgServer(ctl)
	ibcModule.EXPECT().OnRecvPacket(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(channeltypes.NewResultAcknowledgement([]byte(""))).Times(1)
	bankKeeper.EXPECT().SendCoins(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).Times(2)
	cctpServer.EXPECT().DepositForBurnWithCaller(gomock.Any(), gomock.Any()).Return(&cctptypes.MsgDepositForBurnWithCallerResponse{
		Nonce: 10,
	}, nil).Times(1)
	middleware := autocctp.NewMiddleware(ibcModule, bankKeeper, cctpServer)

	memo := types.Memo{
		DepositForBurnWithCaller: &types.DepositForBurnWithCaller{
			DepositForBurn: types.DepositForBurn{
				DestinationDomain: 0,
				MintRecipient:     []byte("mintRecipient"),
				Amount:            &cctpAmount,
				FeeRecipient:      &feeRecipient,
			},
			DestinationCaller: []byte("destinationCaller"),
		},
	}
	memobz, _ := json.Marshal(memo)
	testTransferPacket := transfertypes.FungibleTokenPacketData{
		Denom:    "transfer/channel-0/uusdc",
		Amount:   "100",
		Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
		Receiver: authtypes.NewModuleAddress("autocctp").String(),
		Memo:     string(memobz),
	}
	testTransferData, _ := transfertypes.ModuleCdc.MarshalJSON(&testTransferPacket)
	testPacket := channeltypes.Packet{
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               testTransferData,
	}

	res := middleware.OnRecvPacket(ctx, testPacket, nil)

	require.True(t, res.Success())
	var ack channeltypes.Acknowledgement
	err := channeltypes.SubModuleCdc.UnmarshalJSON(res.Acknowledgement(), &ack)
	require.NoError(t, err)
	require.Equal(t, "{\"nonce\":10}", string(ack.GetResult()))
}
