pub mod generated_chains;

#[cfg(test)]
mod tests {
    use super::*;

    use generated_chains::{ChainId, ChainName};
    use std::str::FromStr;

    #[test]
    fn chain_from_id() {
        assert_eq!(
            ChainName::try_from(1 as ChainId).unwrap(),
            ChainName::EthereumMainnet
        );
    }

    #[test]
    fn chain_from_unknown_id() {
        match ChainName::try_from(0 as ChainId) {
            Ok(_) => panic!("should have failed for unknown chain"),
            Err(e) => {
                assert_eq!(e.to_string(), "unknown chain id: 0");
            }
        }
    }

    #[test]
    fn chain_from_str() {
        let chain = ChainName::from_str("ethereum-mainnet").unwrap();
        assert_eq!(chain, ChainName::EthereumMainnet);
    }

    #[test]
    fn chain_from_str_unknown() {
        let e = ChainName::from_str("ethereum-dummy-x").unwrap_err();
        assert_eq!(e.to_string(), "unknown chain name: ethereum-dummy-x");
    }

    #[test]
    fn to_chain_selector() {
        assert_eq!(
            generated_chains::chain_selector(ChainName::EthereumMainnet),
            5009297550715157269
        );
    }
}
