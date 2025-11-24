# Chain Selectors

CCIP uses its own set of chain selectors represented by uint64 to identify blockchains. It is a random
integer generated as follows:

```python3
python3
>>> import random
>>> random.randint(1, 2**64-1)
```

The scheme is used for several reasons:

- Global uniqueness across blockchain families
- Very unlikely to collide with existing chain ID schemes, reducing confusion
- Efficient on/off-chain representation
- No preference towards any family or chain
- Decoupled from chain name which may change over time with rebrands/forks

This repository contains a
mapping between the custom chain identifiers (`chainSelectorId`) chain names and the chain identifiers
used by the blockchains themselves (`chainId`). For solana we use the base58 encoded genesis hash as the chain id.

Please refer to the [official documentation](https://docs.chain.link/ccip/supported-networks) to learn more about
supported networks and their selectors.

### Installation

`go get github.com/smartcontractkit/chain-selectors`

### Usage

```go
import (
    chainselectors "github.com/smartcontractkit/chain-selectors"
)

func main() {
    // -------------------Chains agnostic --------------------:

    // Getting chain family based on selector
    family, err := GetSelectorFamily(2664363617261496610)

    // -------------------For EVM chains--------------------

    // Getting selector based on ChainId
    selector, err := chainselectors.SelectorFromChainId(420)

    // Getting ChainId based on ChainSelector
    chainId, err := chainselectors.ChainIdFromSelector(2664363617261496610)

    // Getting ChainName based on ChainId
    chainName, err := chainselectors.NameFromChainId(420)

    // Getting ChainId based on the ChainName
    chainId, err := chainselectors.ChainIdFromName("binance_smart_chain-testnet")

    // Accessing mapping directly
    lookupChainId := uint64(1337)
    if chainSelector, exists := chainselectors.EvmChainIdToChainSelector()[lookupChainId]; exists {
        fmt.Println("Found evm chain selector for chain", lookupChainId, ":", chainSelector)
    }

    // -------------------Solana Chain --------------------:

    // Getting chain family based on selector
    family, err := SolanaNameFromChainId("5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d")

    // Getting chain id from chain selector
	chainId, err := chainselectors.SolanaChainIdFromSelector(124615329519749607)

    // Accessing mapping directly
    lookupChainId := "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"
    if chainSelector, exists:= chainselectors.SolanaChainIdToChainSelector()[lookupChainId]; exists {
        fmt.Println("Found solana chain selector for chain", lookupChainId, ":", chainSelector)
    }
}
```

### Adding additional chains at runtime

You can add additional chains at runtime by setting the `EXTRA_SELECTORS_FILE` environment variable to point to either:

1. A **local YAML file** path
2. An **HTTP/HTTPS URL** to a remote YAML file

This is useful for adding custom chains, test networks, or getting the latest chains without updating the library version.

#### Using the Official All Chains File

To automatically get the latest chains without updating your dependency version:

```bash
export EXTRA_SELECTORS_FILE=https://raw.githubusercontent.com/smartcontractkit/chain-selectors/main/all_selectors.yml
```

This file contains all production chains from all blockchain families and is updated whenever new chains are added to the repository.

**Note:** You can also pin to a specific version or commit:

```bash
# Pin to specific release tag
export EXTRA_SELECTORS_FILE=https://raw.githubusercontent.com/smartcontractkit/chain-selectors/v1.2.3/all_selectors.yml

# Pin to specific commit
export EXTRA_SELECTORS_FILE=https://raw.githubusercontent.com/smartcontractkit/chain-selectors/abc123def/all_selectors.yml
```

#### Using a Custom File

```bash
# Local file
export EXTRA_SELECTORS_FILE=/path/to/custom-chains.yml

# Remote custom file
export EXTRA_SELECTORS_FILE=https://raw.githubusercontent.com/your-org/configs/main/custom-chains.yml
```

#### File Format

The YAML file must use the multi-family format with blockchain family keys:

```yaml
evm:
  999999:
    selector: 1234567890123456789
    name: "my-custom-chain"
solana:
  "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d":
    selector: 1111111111111111111
    name: "solana-custom"
aptos:
  4:
    selector: 2222222222222222222
    name: "aptos-custom"
# Also supports: tron, ton, sui, starknet sections as needed
```

#### Behavior Notes

- **Remote URLs** are fetched once at process startup with a 30-second timeout
- **Chains that already exist** in your version are skipped (logged as warnings)
- **Only new chains** are added to the runtime maps
- For **air-gapped environments**, use local file paths
- **File parsing errors** will cause the process to panic (fail-fast)

#### Use Cases

This feature solves the version update problem: downstream services no longer need to update their dependency version every time a new chain is added. Instead, they can point to the remote `all_selectors.yml` file and automatically receive new chains at startup.

### Contributing

#### Naming new chains

Chain names must respect the following format:
`<blockchain>-<type>-<network_instance>`

When a component requires more than 1 word, use snake-case to connect them, e.g `polygon-zkevm`.

| Parameter        | Description                                          | Example                                  |
| ---------------- | ---------------------------------------------------- | ---------------------------------------- |
| blockchain       | Name of the chain                                    | `ethereum`, `avalanche`, `polygon-zkevm` |
| type             | Type of network                                      | `testnet`, `mainnet`, `devnet`           |
| network_instance | [Only if not mainnet] Identifier of specific network | `alfajores`, `holesky`, `sepolia`, `1`   |

More on `network_instance`: only include it if `type` is not mainnet. This is because legacy testnet instances are often dropped after a new one is spun up, e.g Ethereum Rinkeby.

Rules for `network_instance`:

1. If chain has an officially-named testnet, use it, e.g
   `celo-testnet-alfajores`, `ethereum-testnet-holesky`
2. If not above, and chain is a rollup, use the name of its settlement network, e.g `base-testnet-sepolia`
3. If not above, use a number, e.g `bsc-testnet-1`

Example chain names that comply with the format:

```
astar-mainnet
astar-testnet-shibuya
celo-mainnet
celo-testnet-sepolia
polygon-zkevm-mainnet
polygon-zkevm-testnet-cardona
ethereum-mainnet
ethereum-testnet-sepolia
ethereum-testnet-holesky
optimism-mainnet
optimism-testnet-sepolia
bsc-mainnet
bsc-testnet-1
```

You may find some existing names follow a legacy naming pattern: `<blockchain>-<type>-<network_name>-<parachain>-<rollup>-<rollup_instance>`. Those names are kept as is due to complexity of migration. The transition form legacy pattern to the new pattern is motivated by chain migrations, e.g Celo migrating from an L1 into an L2, rendering the legacy name stale.

#### Adding new chains

Any new chains and selectors should be always added to [selectors.yml](selectors.yml) and client libraries should load
details from this file. This ensures that all client libraries are in sync and use the same mapping.
To add a new chain, please add new entry to the `selectors.yml` file and use the following format:

Make sure to run `go generate` after making any changes.

```yaml
$chain_id:
  selector: $chain_selector as uint64
  name: $chain_name as string # Although name is optional parameter, please provide it and respect the format described below
```

[selectors.yml](selectors.yml) file is divided into sections based on the blockchain type.
Please make sure to add new entries to the both sections and keep them sorted by chain id within these sections.

If you need to add a new chain for testing purposes (e.g. running tests with simulated environment) don't mix it with
the main file and use [test_selectors.yml](test_selectors.yml) instead. This file is used only for testing purposes.

#### Adding new client libraries

If you need a support for a new language, please open a PR with the following changes:

- Library codebase is in a separate directory
- Library uses selectors.yml as a source of truth
- Proper Github workflow is present to make sure code compiles and tests pass
