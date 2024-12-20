// Code generated by go generate please DO NOT EDIT
package chain_selectors

type SolanaChain struct {
	ChainID  string
	Selector uint64
	Name     string
	VarName  string
}

var (
	SOLANA_DEVNET  = SolanaChain{ChainID: "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG", Selector: 16423721717087811551, Name: "solana-devnet"}
	SOLANA_MAINNET = SolanaChain{ChainID: "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d", Selector: 124615329519749607, Name: "solana-mainnet"}
	SOLANA_TESTNET = SolanaChain{ChainID: "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY", Selector: 6302590918974934319, Name: "solana-testnet"}
)

var SolanaALL = []SolanaChain{
	SOLANA_DEVNET,
	SOLANA_MAINNET,
	SOLANA_TESTNET,
}
