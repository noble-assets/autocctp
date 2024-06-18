package autocctp

import (
	"encoding/json"
	"errors"
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
	app porttypes.IBCModule

	bankKeeper types.BankKeeper
	server     cctptypes.MsgServer
}

func NewMiddleware(app porttypes.IBCModule, bankKeeper types.BankKeeper, keeper *cctpkeeper.Keeper) Middleware {
	return Middleware{
		app:        app,
		bankKeeper: bankKeeper,
		server:     cctpkeeper.NewMsgServerImpl(keeper),
	}
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
		return channeltypes.NewErrorAcknowledgement(errors.New("malformed memo"))
	}
	sender := types.GenerateAddress(packet.GetDestChannel(), data.Sender)

	if memo.DepositForBurn != nil && memo.DepositForBurnWithCaller == nil {
		amount, _ := math.NewIntFromString(data.Amount)

		if memo.DepositForBurn.Amount != nil {
			if memo.DepositForBurn.FeeRecipient != nil {
				feeRecipient, err := sdk.AccAddressFromBech32(*memo.DepositForBurn.FeeRecipient)
				if err != nil {
					return channeltypes.NewErrorAcknowledgement(errors.New("failed to decode fee recipient"))
				}

				packetAmount := amount
				amount, ok := math.NewIntFromString(*memo.DepositForBurn.Amount)
				if !ok {
					return channeltypes.NewErrorAcknowledgement(errors.New("failed to decode specified amount"))
				}

				feeAmount := packetAmount.Sub(amount)
				if !feeAmount.IsPositive() {
					return channeltypes.NewErrorAcknowledgement(errors.New("specified amount must be strictly less than packet amount"))
				}

				err = m.bankKeeper.SendCoins(
					ctx, types.ModuleAddress, feeRecipient,
					sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(feeAmount.BigInt()))),
				)
				if err != nil {
					return channeltypes.NewErrorAcknowledgement(errors.New("failed to execute fee transfer"))
				}
			} else {
				return channeltypes.NewErrorAcknowledgement(errors.New("specified amount without a fee recipient"))
			}
		}

		err = m.bankKeeper.SendCoins(
			ctx, types.ModuleAddress, sender,
			sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(amount.BigInt()))),
		)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		msg := &cctptypes.MsgDepositForBurn{
			From:              sender.String(),
			Amount:            amount,
			DestinationDomain: memo.DepositForBurn.DestinationDomain,
			MintRecipient:     memo.DepositForBurn.MintRecipient,
			BurnToken:         denom,
		}

		goCtx := sdk.WrapSDKContext(ctx)
		res, err := m.server.DepositForBurn(goCtx, msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		return channeltypes.NewResultAcknowledgement([]byte(fmt.Sprintf("{\"nonce\":%d}", res.Nonce)))
	} else if memo.DepositForBurn == nil && memo.DepositForBurnWithCaller != nil {
		amount, _ := math.NewIntFromString(data.Amount)

		if memo.DepositForBurnWithCaller.Amount != nil {
			if memo.DepositForBurnWithCaller.FeeRecipient != nil {
				feeRecipient, err := sdk.AccAddressFromBech32(*memo.DepositForBurnWithCaller.FeeRecipient)
				if err != nil {
					return channeltypes.NewErrorAcknowledgement(errors.New("failed to decode fee recipient"))
				}

				packetAmount := amount
				amount, ok := math.NewIntFromString(*memo.DepositForBurnWithCaller.Amount)
				if !ok {
					return channeltypes.NewErrorAcknowledgement(errors.New("failed to decode specified amount"))
				}

				feeAmount := packetAmount.Sub(amount)
				if !feeAmount.IsPositive() {
					return channeltypes.NewErrorAcknowledgement(errors.New("specified amount must be strictly less than packet amount"))
				}

				err = m.bankKeeper.SendCoins(
					ctx, types.ModuleAddress, feeRecipient,
					sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(feeAmount.BigInt()))),
				)
				if err != nil {
					return channeltypes.NewErrorAcknowledgement(errors.New("failed to execute fee transfer"))
				}
			} else {
				return channeltypes.NewErrorAcknowledgement(errors.New("specified amount without a fee recipient"))
			}
		}

		err = m.bankKeeper.SendCoins(
			ctx, types.ModuleAddress, sender,
			sdk.NewCoins(sdk.NewCoin(denom, sdk.NewIntFromBigInt(amount.BigInt()))),
		)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		msg := &cctptypes.MsgDepositForBurnWithCaller{
			From:              sender.String(),
			Amount:            amount,
			DestinationDomain: memo.DepositForBurnWithCaller.DestinationDomain,
			MintRecipient:     memo.DepositForBurnWithCaller.MintRecipient,
			BurnToken:         denom,
			DestinationCaller: memo.DepositForBurnWithCaller.DestinationCaller,
		}

		goCtx := sdk.WrapSDKContext(ctx)
		res, err := m.server.DepositForBurnWithCaller(goCtx, msg)
		if err != nil {
			return channeltypes.NewErrorAcknowledgement(err)
		}

		return channeltypes.NewResultAcknowledgement([]byte(fmt.Sprintf("{\"nonce\":%d}", res.Nonce)))
	} else {
		return channeltypes.NewErrorAcknowledgement(errors.New("malformed memo"))
	}
}

func (m Middleware) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return m.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (m Middleware) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return m.app.OnTimeoutPacket(ctx, packet, relayer)
}
