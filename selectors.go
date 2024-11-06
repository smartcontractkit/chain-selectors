package chain_selectors

import "fmt"

const (
	FamilyEVM      = "evm"
	FamilySolana   = "solana"
	FamilyStarknet = "starknet"
	FamilyCosmos   = "cosmos"
	FamilyAptos    = "aptos"
)

func GetSelectorFamily(selector uint64) (string, error) {
	if _, exist := evmChainsBySelector[selector]; exist {
		return FamilyEVM, nil
	}
	if _, exist := solanaChainIdBySelector[selector]; exist {
		return FamilySolana, nil
	}
	return "", fmt.Errorf("unknown chain selector %d", selector)
}
