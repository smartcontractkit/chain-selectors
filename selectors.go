package chain_selectors

import (
	"fmt"
	"strconv"
)

const (
	FamilyEVM      = "evm"
	FamilySolana   = "solana"
	FamilyStarknet = "starknet"
	FamilyCosmos   = "cosmos"
	FamilyAptos    = "aptos"
)

func GetSelectorFamily(selector uint64) (string, error) {
	// check EVM
	_, exist := evmChainsBySelector[selector]
	if exist {
		return FamilyEVM, nil
	}

	// check solana
	_, exist = solanaChainIdBySelector[selector]
	if exist {
		return FamilySolana, nil
	}

	// check aptos
	_, exist = aptosChainIdBySelector[selector]
	if exist {
		return FamilyAptos, nil
	}

	return "", fmt.Errorf("unknown chain selector %d", selector)
}

func GetChainIDFromSelector(selector uint64) (string, error) {
	destChainFamily, err := GetSelectorFamily(selector)
	if err != nil {
		return "", err
	}

	var id uint64
	var destChainID string
	switch destChainFamily {
	case FamilyEVM:
		id, err = ChainIdFromSelector(selector)
		if err != nil {
			return "", fmt.Errorf("failed to get %v chain ID from selector %d: %w", destChainFamily, selector, err)
		}
		destChainID = fmt.Sprintf("%d", id)
	case FamilySolana:
		destChainID, err = SolanaChainIdFromSelector(selector)
		if err != nil {
			return "", fmt.Errorf("failed to get %v chain ID from selector %d: %w", destChainFamily, selector, err)
		}
	case FamilyAptos:
		id, err = AptosChainIdFromSelector(selector)
		if err != nil {
			return "", fmt.Errorf("failed to get %v chain ID from selector %d: %w", destChainFamily, selector, err)
		}
		destChainID = fmt.Sprintf("%d", id)
	}

	return destChainID, nil
}

func GetChainDetailsByChainIDAndFamily(chainID string, family string) (ChainDetails, error) {
	switch family {
	case FamilyEVM:
		evmChainId, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := evmChainIdToChainSelector[evmChainId]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilySolana:
		details, exist := solanaSelectorsMap[chainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilyAptos:
		aptosChainId, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := aptosSelectorsMap[aptosChainId]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	default:
		return ChainDetails{}, fmt.Errorf("family %s is not yet support", family)
	}
}
