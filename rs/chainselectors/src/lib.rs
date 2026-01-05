pub mod generated_chains;

#[cfg(test)]
mod tests {
    use super::*;

    use generated_chains::{ChainId, ChainName, ChainSelector};
    use std::str::FromStr;

    #[test]
    fn chain_from_id() {
        assert_eq!(
            ChainName::try_from(ChainId(1)).unwrap(),
            ChainName::EthereumMainnet,
        );
    }

    #[test]
    fn chain_from_unknown_id() {
        assert_eq!(
            ChainName::try_from(ChainId(0)).unwrap_err().to_string(),
            "unknown chain id: 0",
        );
    }

    #[test]
    fn chain_from_str() {
        assert_eq!(
            ChainName::from_str("ethereum-mainnet").unwrap(),
            ChainName::EthereumMainnet,
        );
    }

    #[test]
    fn chain_from_str_unknown() {
        let e = ChainName::from_str("ethereum-dummy-x").unwrap_err();
        assert_eq!(e.to_string(), "unknown chain name: ethereum-dummy-x");
    }

    #[test]
    fn chain_selector_from_chain() {
        assert_eq!(
            ChainSelector::from(ChainName::EthereumMainnet),
            ChainSelector(5009297550715157269),
        );
    }

    #[test]
    fn chain_from_selector() {
        assert_eq!(
            ChainName::try_from(ChainSelector(5009297550715157269)).unwrap(),
            ChainName::EthereumMainnet,
        );
    }

    #[test]
    fn chain_from_selector_unknown() {
        assert_eq!(
            ChainName::try_from(ChainSelector(1))
                .unwrap_err()
                .to_string(),
            "unknown chain selector: 1"
        );
    }
}
