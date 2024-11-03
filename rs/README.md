# Chain Selectors (Rust)

This folder contains the Rust implementation of the Chain Selectors module.

## Installation

```shell
cargo add --git https://github.com/smartcontractkit/chain-selectors --tag <TAG>
```

## Usage

```rust
use std::str::FromStr;
use chainselectors::generated_chains;

fn main() {
    match generated_chains::ChainName::try_from(420) {
        Ok(c) => {
            assert_eq!(c, generated_chains::ChainName::EthereumTestnetGoerliOptimism1);
        }
        Err(_) => {
            panic!("Failed to convert chain id to chain name");
        }
    }
    
    let selector = generated_chains::chain_selector(generated_chains::ChainName::EthereumTestnetGoerliOptimism1);
    assert_eq!(
        selector,
        2664363617261496610,
    );

    match generated_chains::ChainName::from_str("ethereum-testnet-goerli-optimism-1") {
        Ok(chain) => {
            assert_eq!(chain, generated_chains::ChainName::EthereumTestnetGoerliOptimism1);
        }
        Err(_) => {
            panic!("Failed to parse chain name");
        }
    }
}
```

## Dev Guide

### Pre-requisites

As part of the code generation, `rustfmt` is used to format the generated code. If you don't have rust toolchain installed, it will run in docker.

### Build

To build the Chain Selectors module, run the `go generate ./rs` from the root of the project, or `go generate ./...` to generate all modules (not just rust).
