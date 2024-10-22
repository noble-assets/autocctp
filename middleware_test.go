package autocctp_test

import (
	"fmt"
	"testing"

	autocctp "autocctp.dev"
	mock "autocctp.dev/test/mock"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// 3. types.Memo is empty
// 4. Packet is deposit for burn - amount is nil
// 5. Packet is deposit for burn - fee recipient is nil
// 6. Packet is deposit for burn - fee recipient is invalid bech32
// 7. Packet is deposit for burn - amount is invalid
// 8. Packet is deposit for burn - specified amount is greater than packet amount
// 9. Packet is deposit for burn - fee transfer failed
// 10. Packet is deposit for burn - deposit success - assert ack
// 11. Packet is deposit for burn with caller - amount is nil
// 12. Packet is deposit for burn with caller - fee recipient is nil
// 13. Packet is deposit for burn with caller - fee recipient is invalid bech32
// 14. Packet is deposit for burn with caller - amount is invalid
// 15. Packet is deposit for burn with caller - specified amount is greater than packet amount
// 16. Packet is deposit for burn with caller - fee transfer failed
// 17. Packet is deposit for burn with caller - deposit success - assert ack

var (
	transferDenom = "transfer"
)

func TestMiddleware(t *testing.T) {
	testcases := []struct {
		name          string
		getPacket     func() channeltypes.Packet
		expectSuccess bool
	}{
		{
			name: "Packet not transfer type",
			getPacket: func() channeltypes.Packet {
				return channeltypes.Packet{}
			},
			expectSuccess: true,
		},
		{
			name: "Receiver is not the autocctp module address",
			getPacket: func() channeltypes.Packet {
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "testDenom",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: "cosmos1vzxkv3lxccnttr9rs0002s93sgw72h7ghukuhs",
				}
				transferData, err := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				require.NoError(t, err)
				return channeltypes.Packet{
					SourcePort:         "transfer",
					SourceChannel:      "channel-0",
					DestinationPort:    "transfer",
					DestinationChannel: "channel-1",
					Data:               transferData,
				}
			},
			expectSuccess: true,
		},
		{
			name: "Sender chain is source chain",
			getPacket: func() channeltypes.Packet {
				transferPacket := transfertypes.FungibleTokenPacketData{
					Denom:    "testDenom",
					Amount:   "100",
					Sender:   "cosmos1wnlew8ss0sqclfalvj6jkcyvnwq79fd74qxxue",
					Receiver: authtypes.NewModuleAddress("autocctp").String(),
				}
				transferData, err := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
				require.NoError(t, err)
				return channeltypes.Packet{
					SourcePort:         "transfer",
					SourceChannel:      "channel-0",
					DestinationPort:    "transfer",
					DestinationChannel: "channel-1",
					Data:               transferData,
				}
			},
			expectSuccess: true,
		},
	}
	for _, tc := range testcases {
		t.Run(fmt.Sprintf("Case: %s", tc.name), func(t *testing.T) {
			db := dbm.NewMemDB()
			logger := log.NewTestLogger(t)
			stateStore := store.NewCommitMultiStore(db, logger, metrics.NewNoOpMetrics())
			ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)
			relayer := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			bankKeeper := mock.NewMockBankKeeper(ctl)
			ibcModule := mock.NewMockIBCModule(ctl)
			middleware := autocctp.NewMiddleware(ibcModule, bankKeeper, nil)
			packet := tc.getPacket()
			gomock.InOrder(
				ibcModule.EXPECT().OnRecvPacket(ctx, packet, relayer).Return(channeltypes.NewResultAcknowledgement([]byte(""))),
			)

			ack := middleware.OnRecvPacket(ctx, packet, relayer)

			if tc.expectSuccess {
				require.True(t, ack.Success())
			} else {
				require.False(t, ack.Success())
			}
		})
	}
}
