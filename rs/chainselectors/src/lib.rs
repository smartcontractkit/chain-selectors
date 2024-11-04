pub mod generated_chains;

#[cfg(test)]
mod tests {
    use super::*;

    use generated_chains::ChainName;
    use std::str::FromStr;

    #[test]
    fn chain_from_id() {
        assert_eq!(
            ChainName::from_chain_id(1),
            Some(ChainName::EthereumMainnet),
        );
    }

    #[test]
    fn chain_from_unknown_id() {
        assert_eq!(ChainName::from_chain_id(0), None);
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
    fn to_chain_selector() {
        assert_eq!(
            generated_chains::chain_selector(ChainName::EthereumMainnet),
            5009297550715157269
        );
    }

    #[test]
    fn chain_from_selector() {
        assert_eq!(
            ChainName::from_chain_selector(5009297550715157269),
            Some(ChainName::EthereumMainnet),
        );
    }

    #[test]
    fn chain_from_selector_unknown() {
        assert_eq!(ChainName::from_chain_selector(1), None);
    }
}
