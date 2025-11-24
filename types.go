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
)

type ChainDetails struct {
	ChainSelector uint64 `yaml:"selector"`
	ChainName     string `yaml:"name"`
}
