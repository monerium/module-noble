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
  florind genesis add-genesis-account validator 1000000ustake --home .florin --keyring-backend test
  BLACKLIST_OWNER=$(florind keys add blacklist-owner --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account blacklist-owner 10000000uusdc --home .florin --keyring-backend test
  BLACKLIST_PENDING_OWNER=$(florind keys add blacklist-pending-owner --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account blacklist-pending-owner 10000000uusdc --home .florin --keyring-backend test
  BLACKLIST_ADMIN=$(florind keys add blacklist-admin --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account blacklist-admin 10000000uusdc --home .florin --keyring-backend test
  AUTHORITY=$(florind keys add authority --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account authority 1000000uusdc --home .florin --keyring-backend test
  OWNER=$(florind keys add owner --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account owner 10000000uusdc --home .florin --keyring-backend test
  PENDING_OWNER=$(florind keys add pending-owner --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account pending-owner 10000000uusdc --home .florin --keyring-backend test
  SYSTEM=$(florind keys add system --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account system 10000000uusdc --home .florin --keyring-backend test
  ADMIN=$(florind keys add admin --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account admin 10000000uusdc --home .florin --keyring-backend test
  florind keys add alice --home .florin --keyring-backend test --output json &> /dev/null
  florind genesis add-genesis-account alice 10000000uusdc --home .florin --keyring-backend test
  BOB=$(florind keys add bob --home .florin --keyring-backend test --output json | jq .address)
  florind genesis add-genesis-account bob 10000000uusdc --home .florin --keyring-backend test

  TEMP=.florin/genesis.json
  touch $TEMP && jq '.app_state.staking.params.bond_denom = "ustake"' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.blacklist_state.owner = '$BLACKLIST_OWNER'' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.blacklist_state.admins = ['$BLACKLIST_ADMIN']' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.blacklist_state.adversaries = ['$BOB']' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.authority = '$AUTHORITY'' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.owners = { "ueure": '$OWNER' }' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.systems = [{ "denom": "ueure", "address": '$SYSTEM'}]' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.admins = [{ "denom": "ueure", "address": '$ADMIN'}]' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json
  touch $TEMP && jq '.app_state.florin.mint_allowances = [{ "denom": "ueure", "address": '$SYSTEM', "allowance": "1000000000000"}]' .florin/config/genesis.json > $TEMP && mv $TEMP .florin/config/genesis.json

  florind genesis gentx validator 1000000ustake --chain-id "florin-1" --home .florin --keyring-backend test &> /dev/null
  florind genesis collect-gentxs --home .florin &> /dev/null

  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' .florin/config/config.toml
fi

florind start --home .florin
