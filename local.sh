alias florind=./simapp/build/simd

for arg in "$@"
do
    case $arg in
        -r|--reset)
        rm -rf .florin
        shift
        ;;
    esac
done

if ! [ -f .florin/data/priv_validator_state.json ]; then
  florind init validator --chain-id "florin-1" --home .florin &> /dev/null

  florind keys add validator --home .florin --keyring-backend test &> /dev/null
  florind add-genesis-account validator 1000000ustake --home .florin --keyring-backend test
  BLACKLIST_OWNER=$(florind keys add blacklist-owner --home .florin --keyring-backend test --output json | jq .address)
  florind add-genesis-account blacklist-owner 10000000uusdc --home .florin --keyring-backend test
  BLACKLIST_PENDING_OWNER=$(florind keys add blacklist-pending-owner --home .florin --keyring-backend test --output json | jq .address)
  florind add-genesis-account blacklist-pending-owner 10000000uusdc --home .florin --keyring-backend test

  TEMP=.florin/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.blacklist_state.owner = '$BLACKLIST_OWNER'' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json

  florind gentx validator 1000000ustake --chain-id "florin-1" --home .florin --keyring-backend test &> /dev/null
  florind collect-gentxs --home .florin &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .florin/config/config.toml
fi

florind start --home .florin
