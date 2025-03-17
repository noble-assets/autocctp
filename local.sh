#!/bin/bash
SIMD="./simapp/build/simd"

for arg in "$@"; do
    case $arg in
    -r | --reset)
        rm -rf .autocctp
        shift
        ;;
    esac
done

if ! [ -f .autocctp/data/priv_validator_state.json ]; then
    $SIMD init validator --chain-id "autocctp-1" --home .autocctp &>/dev/null

    $SIMD keys add validator --home .autocctp --keyring-backend test &>/dev/null
    $SIMD genesis add-genesis-account validator 2000000ustake,1000000000uusdc --home .autocctp --keyring-backend test

    TEMP=.autocctp/genesis.json
    touch $TEMP && jq '.app_state.bank.denom_metadata += [{ "description": "Circle USD Coin", "denom_units": [{ "denom": "uusdc", "exponent": 0, "aliases": ["microusdc"] }, { "denom": "usdc", "exponent": 6 }], "base": "uusdc", "display": "usdc", "name": "Circle USD Coin", "symbol": "USDC" }]' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json
    touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json

    touch $TEMP && jq '.app_state."fiat-tokenfactory".mintingDenom = { "denom": "uusdc" }' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json
    touch $TEMP && jq '.app_state."fiat-tokenfactory".paused = { "paused": false }' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json
    touch $TEMP && jq '.app_state."fiat-tokenfactory".mintersList += [{"address": "noble12l2w4ugfz4m6dd73yysz477jszqnfughxvkss5", "allowance": { "denom": "uusdc", "amount": "1000000000000" }}]' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json

    touch $TEMP && jq '.app_state.cctp.per_message_burn_limit_list += [{ "denom": "uusdc", "amount": "1000000000000" }]' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json
    touch $TEMP && jq '.app_state.cctp.token_messenger_list +=[{"domain_id": "0", "address": "AAAAAAAAAAAAAAAAvT+oG1i6kqghNgOLJa3scGavMVU="}]' .autocctp/config/genesis.json >$TEMP && mv $TEMP .autocctp/config/genesis.json

    $SIMD genesis gentx validator 1000000ustake --chain-id "autocctp-1" --home .autocctp --keyring-backend test &>/dev/null
    $SIMD genesis collect-gentxs --home .autocctp &>/dev/null

    sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .autocctp/config/config.toml
fi

$SIMD start --home .autocctp
