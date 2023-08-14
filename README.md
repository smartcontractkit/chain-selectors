# CCIP chain selectors

CCIP uses its own set of chain selectors represented by uint64 to identify blockchains. This repository contains a mapping between the CCIP chain identifiers (`chainSelectorId`) and the chain identifiers used by the blockchains themselves (`chainId`).

Please refer to the [official documentation](https://docs.chain.link/ccip/supported-networks) to learn more about supported networks and their selectors.

### Installation

`go get github.com/smartcontract/ccip-chain-selectors`

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
    if chainSelector, exists := selectors.EvmChainIdToChainSelector[lookupChainId]; exists {
        fmt.Println("Found chain selector for chain", lookupChainId, ":", chainSelector)
    }
}

```
