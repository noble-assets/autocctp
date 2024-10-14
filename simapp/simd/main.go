package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	simapp "autocctp.dev/simapp"

	"github.com/noble-assets/noble/v7/cmd"
	//"autocctp.dev/simapp/simd/cmd"
	"github.com/noble-assets/noble/v7/app"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.ChainID,
		app.ModuleBasics,
		simapp.New,
	)

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
