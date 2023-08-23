# Chain Selectors

CCIP uses its own set of chain selectors represented by uint64 to identify blockchains. This repository contains a
mapping between the CCIP chain identifiers (`chainSelectorId`) and the chain identifiers used by the blockchains
themselves (`chainId`).

Please refer to the [official documentation](https://docs.chain.link/ccip/supported-networks) to learn more about
supported networks and their selectors.

### Installation

`go get github.com/smartcontractkit/chain-selectors`

### Usage

```go
import (
    "github.com/smartcontractkit/ccip-chain-selectors"
)

func main() {
    // Getting selector based on ChainId
    selector, err := selectors.SelectorFromChainId(420)

    // Getting ChainId based on ChainSelector
    chainId, err := selectors.ChainIdFromSelector(2664363617261496610)

    // Accessing mapping directly
    lookupChainId := uint64(1337)
    if chainSelector, exists := selectors.EvmChainIdToChainSelector()[lookupChainId]; exists {
        fmt.Println("Found chain selector for chain", lookupChainId, ":", chainSelector)
    }
}
```

### Contributing

Any new chains and selectors should be always added to [selectors.yml](selectors.yml) and client libraries should load
details from this file. This ensures that all client libraries are in sync and use the same mapping.

If you need a support for a new language, please open a PR with the following changes:
- Library codebase is in a separate directory
- Library uses selectors.yml as a source of truth
- Proper Github workflow is present to make sure code compiles and tests pass
