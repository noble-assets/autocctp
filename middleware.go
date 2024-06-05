package autocctp

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"

	cctpkeeper "github.com/circlefin/noble-cctp/x/cctp/keeper"
	cctptypes "github.com/circlefin/noble-cctp/x/cctp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v4/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v4/modules/core/exported"
	"github.com/noble-assets/autocctp/types"
)

var _ porttypes.IBCModule = &Middleware{}

type Middleware struct {
	app    porttypes.IBCModule
	keeper *cctpkeeper.Keeper
}

func NewMiddleware(app porttypes.IBCModule, keeper *cctpkeeper.Keeper) Middleware {
	return Middleware{app: app, keeper: keeper}
}

func (m Middleware) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, version string) (string, error) {
	return m.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (m Middleware) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty channeltypes.Counterparty, counterpartyVersion string) (version string, err error) {
	return m.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

func (m Middleware) OnChanOpenAck(ctx sdk.Context, portID string, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return m.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (m Middleware) OnChanOpenConfirm(ctx sdk.Context, portID string, channelID string) error {
	return m.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (m Middleware) OnChanCloseInit(ctx sdk.Context, portID string, channelID string) error {
	return m.app.OnChanCloseInit(ctx, portID, channelID)
}

func (m Middleware) OnChanCloseConfirm(ctx sdk.Context, portID string, channelID string) error {
	return m.app.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the porttypes.IBCModule interface.
func (m Middleware) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	ack := m.app.OnRecvPacket(ctx, packet, relayer)
	if !ack.Success() {
		return ack
	}

	var data transfertypes.FungibleTokenPacketData
	_ = transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data)
	if data.Receiver != types.ModuleAddress.String() {
		return ack
	}

	if transfertypes.SenderChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		return ack
	}
	voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
	unprefixedDenom := data.Denom[len(voucherPrefix):]
	denom := unprefixedDenom
	denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
	if denomTrace.Path != "" {
		denom = denomTrace.IBCDenom()
	}

	var memo types.Memo
	err := json.Unmarshal([]byte(data.GetMemo()), &memo)
	if err != nil {
		return channeltypes.NewErrorAcknowledgement(types.ErrMalformedMemo)
	}

	amount, _ := math.NewIntFromString(data.Amount)
	if memo.DepositForBurn != nil && memo.DepositForBurnWithCaller == nil {
		msg := &cctptypes.MsgDepositForBurn{
			From:              types.GenerateAddress(packet.GetDestChannel(), data.Sender).String(),
			Amount:            amount,
			DestinationDomain: memo.DepositForBurn.DestinationDomain,
			MintRecipient:     memo.DepositForBurn.MintRecipient,
			BurnToken:         denom,
		}

		goCtx := sdk.WrapSDKContext(ctx)
		res, err := cctpkeeper.NewMsgServerImpl(m.keeper).DepositForBurn(goCtx, msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		return channeltypes.NewResultAcknowledgement([]byte(fmt.Sprintf("{\"nonce\":%d}", res.Nonce)))
	} else if memo.DepositForBurn == nil && memo.DepositForBurnWithCaller != nil {
		msg := &cctptypes.MsgDepositForBurnWithCaller{
			From:              types.GenerateAddress(packet.GetDestChannel(), data.Sender).String(),
			Amount:            amount,
			DestinationDomain: memo.DepositForBurnWithCaller.DestinationDomain,
			MintRecipient:     memo.DepositForBurnWithCaller.MintRecipient,
			BurnToken:         denom,
			DestinationCaller: memo.DepositForBurnWithCaller.DestinationCaller,
		}

		goCtx := sdk.WrapSDKContext(ctx)
		res, err := cctpkeeper.NewMsgServerImpl(m.keeper).DepositForBurnWithCaller(goCtx, msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		return channeltypes.NewResultAcknowledgement([]byte(fmt.Sprintf("{\"nonce\":%d}", res.Nonce)))
	} else {
		return channeltypes.NewErrorAcknowledgement(types.ErrMalformedMemo)
	}
}

func (m Middleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return m.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (m Middleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return m.app.OnTimeoutPacket(ctx, packet, relayer)
}
