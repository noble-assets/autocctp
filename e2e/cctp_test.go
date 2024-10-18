package e2e

import (
	"context"
	"testing"

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

	err = testutil.WaitForBlocks(ctx, 10, noble, gaia)
	require.NoError(t, err, "failed to wait for blocks")

	// Step 1: Mint some USDC
	// Step 2: Send USDC to gaia
	// Step 3: Check USDC balance on gaia
	// Step 4: Send 1/2 USDC back to noble with autocctp msg
	// Step 5: Check USDC balance on noble - should not exist here
	// Setp 6: Send the remaining 1/2 USDC back to noble without autocctp msg
}
