# Chain Selectors (Rust)

This folder contains the Rust implementation of the Chain Selectors module.

## User Guide

### Installation

```shell
cargo add --git https://github.com/smartcontractkit/chain-selectors --tag <TAG>
```

### Usage

```rust
use std::str::FromStr;
use chainselectors::generated_chains::{ChainName, ChainSelector, ChainId};

fn main() {
    let chain = ChainName::try_from(ChainId(420)).unwrap();
    assert_eq!(chain, ChainName::EthereumTestnetGoerliOptimism1);
    
    let selector = ChainSelector::from(ChainName::EthereumTestnetGoerliOptimism1);
    assert_eq!(
        selector,
        ChainSelector(2664363617261496610),
    );

    let chain_from_str = ChainName::from_str("ethereum-testnet-goerli-optimism-1").unwrap();
    assert_eq!(chain_from_str, ChainName::EthereumTestnetGoerliOptimism1);
}
```

## Dev Guide

### Pre-requisites

As part of the code generation, `rustfmt` is used to format the generated code. If you don't have rust toolchain installed, it will run in docker.

### Build

To build the Chain Selectors module, run the `go generate ./rs` from the root of the project, or `go generate ./...` to generate all modules (not just rust).
