package simapp

import (
	"github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory"
	ftfkeeper "github.com/circlefin/noble-fiattokenfactory/x/fiattokenfactory/keeper"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	autocctp "autocctp.dev"
	"autocctp.dev/types"
)

type BankKeeper interface {
	authtypes.BankKeeper
	types.BankKeeper
}

type HandlerOptions struct {
	cdc codec.Codec
	ante.HandlerOptions
	FTFKeeper  *ftfkeeper.Keeper
	BankKeeper BankKeeper
}

func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errors.Wrap(errorstypes.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, errors.Wrap(errorstypes.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.FTFKeeper == nil {
		return nil, errors.Wrap(errorstypes.ErrLogic, "fiat tokenfactory keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, errors.Wrap(errorstypes.ErrLogic, "sign mode handler is required for ante builder")
	}

	sigVerificationDecorator := autocctp.NewSigVerificationDecorator(
		options.FTFKeeper,
		options.BankKeeper,
		options.AccountKeeper,
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
	)

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		fiattokenfactory.NewIsPausedDecorator(options.cdc, options.FTFKeeper),
		fiattokenfactory.NewIsBlacklistedDecorator(options.FTFKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),

		// Custom signature verification for AutoCCTP accounts.
		sigVerificationDecorator,

		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
