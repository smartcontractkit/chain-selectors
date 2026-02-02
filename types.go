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
	FamilyStellar  = "stellar"
)

// NetworkType represents the type of network (testnet or mainnet)
type NetworkType string

const (
	NetworkTypeTestnet  NetworkType = "testnet"
	NetworkTypeMainnet  NetworkType = "mainnet"
	NetworkTypeLocalnet NetworkType = "localnet"
	NetworkTypeFuturenet NetworkType = "futurenet"
)

type ChainDetails struct {
	ChainSelector uint64      `yaml:"selector"`
	ChainName     string      `yaml:"name"`
	NetworkType   NetworkType `yaml:"network_type"`
}
