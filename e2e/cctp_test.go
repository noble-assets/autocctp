package e2e

import (
	"context"
	"encoding/json"
	"testing"

	types "autocctp.dev/types"
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/ibc"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestCCTP(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx := context.Background()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	client, network := interchaintest.DockerSetup(t)

	var gw NobleWrapper

	numValidators, numFullNodes := 1, 0
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		NobleChainSpec(ctx, &gw, "noble-1", 2, 1, true),
		{
			Name:          "gaia",
			Version:       "v14.1.0",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	gaia := chains[1].(*cosmos.CosmosChain)

	gw.Chain = chains[0].(*cosmos.CosmosChain)
	noble := gw.Chain

	rly := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	).Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(noble).
		AddChain(gaia).
		AddRelayer(rly, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  noble,
			Chain2:  gaia,
			Relayer: rly,
			Path:    "transfer",
		})

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	}))

	err = rly.StartRelayer(ctx, eRep, "transfer")
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	nobleValidator := noble.Validators[0]
	mintAmount := math.NewInt(1000000000000)
	halfMintAmount := math.NewInt(500000000000)

	// Step 1: Mint 1000000000000 USDC

	nobleUser := interchaintest.GetAndFundTestUsers(t, ctx, "wallet", mintAmount, noble)[0]
	gaiaUser := interchaintest.GetAndFundTestUsers(t, ctx, "wallet", mintAmount, gaia)[0]

	_, err = nobleValidator.ExecTx(ctx, gw.FiatTfRoles.MasterMinter.KeyName(),
		"fiat-tokenfactory", "configure-minter-controller", gw.FiatTfRoles.MinterController.FormattedAddress(), gw.FiatTfRoles.Minter.FormattedAddress(),
	)
	require.NoError(t, err, "failed to execute configure minter controller tx")

	_, err = nobleValidator.ExecTx(ctx, gw.FiatTfRoles.MinterController.KeyName(),
		"fiat-tokenfactory", "configure-minter", gw.FiatTfRoles.Minter.FormattedAddress(), mintAmount.String()+DenomMetadataUsdc.Base,
	)
	require.NoError(t, err, "failed to execute configure minter tx")

	_, err = nobleValidator.ExecTx(ctx, gw.FiatTfRoles.Minter.KeyName(),
		"fiat-tokenfactory", "mint", nobleUser.FormattedAddress(), mintAmount.String()+DenomMetadataUsdc.Base,
	)
	require.NoError(t, err, "failed to execute mint to user tx")

	balance, err := noble.GetBalance(ctx, nobleUser.FormattedAddress(), DenomMetadataUsdc.Base)
	require.NoError(t, err, "failed to get balance")
	require.Equal(t, mintAmount.String(), balance.String())

	// Step 2: Send 1000000000000 USDC to gaia
	srcTx, err := noble.SendIBCTransfer(ctx, "channel-0", nobleUser.KeyName(), ibc.WalletAmount{
		Address: gaiaUser.FormattedAddress(),
		Denom:   DenomMetadataUsdc.Base,
		Amount:  mintAmount,
	}, ibc.TransferOptions{})
	require.NoError(t, err, "failed to execute transfer tx")
	nobleHeight, err := noble.Height(ctx)
	require.NoError(t, err, "failed to get height")
	srcAck, err := testutil.PollForAck(ctx, noble, nobleHeight, nobleHeight+10, srcTx.Packet)
	require.NoError(t, err, "failed to poll for ack")
	require.NoError(t, srcAck.Validate(), "invalid acknowledgement on source chain")
	balance, err = noble.GetBalance(ctx, nobleUser.FormattedAddress(), DenomMetadataUsdc.Base)
	require.NoError(t, err, "failed to get balance")
	require.Equal(t, math.ZeroInt().String(), balance.String())

	// Step 3: Check USDC balance on gaia to be 1000000000000
	srcDenomTrace := transfertypes.ParseDenomTrace(transfertypes.GetPrefixedDenom("transfer", "channel-0", DenomMetadataUsdc.Base))
	dstIbcDenom := srcDenomTrace.IBCDenom()
	balance, err = gaia.GetBalance(ctx, gaiaUser.FormattedAddress(), dstIbcDenom)
	require.NoError(t, err, "failed to get balance")
	require.Equal(t, mintAmount.String(), balance.String())

	// Step 4: Send 500000000000 (1/2) USDC back to noble with autocctp msg
	// where 499999999999 will be sent to destination
	mintRecipient := make([]byte, 32)
	copy(mintRecipient[12:], common.FromHex("0xfCE4cE85e1F74C01e0ecccd8BbC4606f83D3FC90"))
	depositAmount := halfMintAmount.Sub(math.OneInt()).String() // 500000000000 - 1
	feeRecipient := nobleUser.FormattedAddress()
	autocctpAcc := authtypes.NewModuleAddress("autocctp")
	memo := types.Memo{
		DepositForBurn: &types.DepositForBurn{
			DestinationDomain: 0,
			MintRecipient:     mintRecipient,
			Amount:            &depositAmount,
			FeeRecipient:      &feeRecipient,
		},
	}
	memoJSON, err := json.Marshal(memo)
	require.NoError(t, err, "failed to marshal memo")
	dstTx, err := gaia.SendIBCTransfer(ctx, "channel-0", gaiaUser.KeyName(), ibc.WalletAmount{
		Address: autocctpAcc.String(),
		Denom:   dstIbcDenom,
		Amount:  halfMintAmount,
	}, ibc.TransferOptions{
		Memo: string(memoJSON),
	})
	require.NoError(t, err, "failed to execute transfer tx")
	gaiaHeight, err := gaia.Height(ctx)
	require.NoError(t, err, "failed to get height")
	dstAck, err := testutil.PollForAck(ctx, gaia, gaiaHeight, gaiaHeight+10, dstTx.Packet)
	require.NoError(t, err, "failed to poll for ack")
	require.NoError(t, dstAck.Validate(), "invalid acknowledgement on source chain")

	err = testutil.WaitForBlocks(ctx, 5, gaia, noble)
	require.NoError(t, err, "failed to wait for blocks")

	balance, err = gaia.GetBalance(ctx, gaiaUser.FormattedAddress(), dstIbcDenom)
	require.NoError(t, err, "failed to get balance")
	require.Equal(t, halfMintAmount.String(), balance.String())

	// Step 5: Check USDC balance on noble - only one USDC should be there as the rest is forwarded
	balance, err = noble.GetBalance(ctx, nobleUser.FormattedAddress(), DenomMetadataUsdc.Base)
	require.NoError(t, err, "failed to get balance")
	require.Equal(t, math.OneInt().String(), balance.String())

	// Setp 6: Send the remaining 500000000000 USDC back to noble without autocctp msg
	dstTx, err = gaia.SendIBCTransfer(ctx, "channel-0", gaiaUser.KeyName(), ibc.WalletAmount{
		Address: nobleUser.FormattedAddress(),
		Denom:   dstIbcDenom,
		Amount:  halfMintAmount,
	}, ibc.TransferOptions{})
	require.NoError(t, err, "failed to execute transfer tx")
	gaiaHeight, err = gaia.Height(ctx)
	require.NoError(t, err, "failed to get height")
	dstAck, err = testutil.PollForAck(ctx, gaia, gaiaHeight, gaiaHeight+10, dstTx.Packet)
	require.NoError(t, err, "failed to poll for ack")
	require.NoError(t, dstAck.Validate(), "invalid acknowledgement on source chain")

	balance, err = gaia.GetBalance(ctx, gaiaUser.FormattedAddress(), dstIbcDenom)
	require.NoError(t, err, "failed to get balance")
	require.Equal(t, math.ZeroInt().String(), balance.String())

	// Step 7: Check USDC balance on noble to be 500000000001
	balance, err = noble.GetBalance(ctx, nobleUser.FormattedAddress(), DenomMetadataUsdc.Base)
	require.NoError(t, err, "failed to get balance")
	require.Greater(t, halfMintAmount.String(), balance.String())
}
