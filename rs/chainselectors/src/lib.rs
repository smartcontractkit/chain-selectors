pub mod generated_chains;

#[cfg(test)]
mod tests {
    use super::*;
    use generated_chains::ChainName;

    #[test]
    fn convert_from_chain_id() {
        assert_eq!(
            ChainName::from_chain_id(1),
            Some(ChainName::EthereumMainnet)
        );
    }

    #[test]
    fn convert_to_chain_selector() {
        assert_eq!(
            generated_chains::chain_selector(ChainName::EthereumMainnet),
            5009297550715157269
        );
    }
}
