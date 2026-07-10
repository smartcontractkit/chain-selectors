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
	NetworkTypeTestnet NetworkType = "testnet"
	NetworkTypeMainnet NetworkType = "mainnet"
)

type ChainDetails struct {
	ChainSelector uint64      `yaml:"selector"`
	ChainName     string      `yaml:"name"`
	NetworkType   NetworkType `yaml:"network_type"`
	// Deprecated marks chains that are discouraged or superseded by a newer version.
	// A deprecated chain may still be live; see SunsetAt for when it goes offline.
	Deprecated bool `yaml:"deprecated,omitempty"`
	// SunsetAt is when the chain is (or was) scheduled to go offline, as an RFC 3339
	// datetime or "2006-01-02" date. Empty means no sunset is set.
	SunsetAt string `yaml:"sunset_at,omitempty"`
}
