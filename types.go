package chain_selectors

const (
	FamilyEVM      = "evm"
	FamilySolana   = "solana"
	FamilyStarknet = "starknet"
	FamilyCosmos   = "cosmos"
	FamilyAptos    = "aptos"
	FamilySui      = "sui"
	FamilyTron     = "tron"
	FamilyTon      = "ton"
	FamilyCanton   = "canton"
)

type ChainDetails struct {
	ChainSelector uint64 `yaml:"selector"`
	ChainName     string `yaml:"name"`
	IsTestnet     bool   `yaml:"is_testnet,omitempty"`
}
