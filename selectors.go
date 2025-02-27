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
	FamilyTron     = "tron"
)

type ChainInfo struct {
	Family       string
	ChainID      string
	ChainDetails ChainDetails
}

type ChainSelectorsObj struct {
	ExtraChains_evm []ChainInfo
}

func NewChainSelectorsObj(chainDetails ChainInfo) (*ChainSelectorsObj, error) {
	var extraChainsEVM []ChainInfo
	// currentChains_Aptos := AptosALL
	// currentChains_Solana := SolanaALL
	// currentChains_Tron := TronALL

	if chainDetails.Family == FamilyEVM {
		// search through EVM
		// TODO can we do this search in O(1) after we sort them out
		for _, chain := range ALL {
			evmChainID, err := strconv.ParseUint(chainDetails.ChainID, 10, 64)
			if err != nil {
				return &ChainSelectorsObj{}, fmt.Errorf("error converting string to uint64: %v", err)
			}

			if chain.EvmChainID == evmChainID || chain.Selector == chainDetails.ChainDetails.ChainSelector {
				// TODO: add chainselector related validation to check that the random number is valid

				// Conflict detected, return currentChains without adding the new one, donot error
				// TODO: maybe throw a warning log
				return &ChainSelectorsObj{}, nil
			}
		}

		// No conflict, add the new chain
		extraChainsEVM = append(extraChainsEVM, chainDetails)

		return &ChainSelectorsObj{ExtraChains_evm: extraChainsEVM}, nil
	}

	// To extend to new family we can extend the logic like below
	// if chainDetails.Family == FamilyAptos {
	// 	// search through aptos
	// 	currentChains_Aptos
	// }

	return &ChainSelectorsObj{}, nil
}

func getChainInfo(selector uint64, csObj *ChainSelectorsObj) (ChainInfo, error) {
	// check EVM
	_, exist := evmChainsBySelector[selector]
	if exist {
		family := FamilyEVM

		evmChainId, err := ChainIdFromSelector(selector)
		if err != nil {
			return ChainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", family, selector, err)
		}

		details, exist := evmChainIdToChainSelector[evmChainId]
		if !exist {
			return ChainInfo{}, fmt.Errorf("invalid chain id %d for %s", evmChainId, family)
		}

		return ChainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", evmChainId),
			ChainDetails: details,
		}, nil
	}

	// check if the chain exist in extraChains_evm, return if it does
	for _, chain := range csObj.ExtraChains_evm {
		if chain.ChainDetails.ChainSelector == selector {
			return ChainInfo{
				Family:  FamilyEVM, // currently assumes every override is EVM chain
				ChainID: chain.ChainID,
				ChainDetails: ChainDetails{
					ChainSelector: chain.ChainDetails.ChainSelector,
					ChainName:     chain.ChainDetails.ChainName,
				},
			}, nil
		}
	}

	// check solana
	_, exist = solanaChainsBySelector[selector]
	if exist {
		family := FamilySolana

		chainID, err := SolanaChainIdFromSelector(selector)
		if err != nil {
			return ChainInfo{}, fmt.Errorf("failed to get %s chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := solanaChainIdToChainSelector[chainID]
		if !exist {
			return ChainInfo{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return ChainInfo{
			Family:       family,
			ChainID:      chainID,
			ChainDetails: details,
		}, nil
	}

	// check aptos
	_, exist = aptosChainIdBySelector[selector]
	if exist {
		family := FamilyAptos

		chainID, err := AptosChainIdFromSelector(selector)
		if err != nil {
			return ChainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := aptosSelectorsMap[chainID]
		if !exist {
			return ChainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return ChainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", chainID),
			ChainDetails: details,
		}, nil
	}

	// check tron
	_, exist = tronChainIdBySelector[selector]
	if exist {
		family := FamilyTron

		chainID, err := TronChainIdFromSelector(selector)
		if err != nil {
			return ChainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := tronSelectorsMap[chainID]
		if !exist {
			return ChainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return ChainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", chainID),
			ChainDetails: details,
		}, nil
	}

	return ChainInfo{}, fmt.Errorf("unknown chain selector %d", selector)
}

func (obj *ChainSelectorsObj) GetSelectorFamily(selector uint64) (string, error) {
	ChainInfo, err := getChainInfo(selector, obj)
	if err != nil {
		return "", fmt.Errorf("unknown chain selector %d", selector)
	}

	return ChainInfo.Family, nil
}

func (obj *ChainSelectorsObj) GetChainIDFromSelector(selector uint64) (string, error) {
	ChainInfo, err := getChainInfo(selector, obj)
	if err != nil {
		return "", fmt.Errorf("unknown chain selector %d", selector)
	}

	return ChainInfo.ChainID, nil
}

func (obj *ChainSelectorsObj) GetChainDetailsByChainIDAndFamily(chainID string, family string) (ChainDetails, error) {
	switch family {
	case FamilyEVM:
		evmChainId, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := evmChainIdToChainSelector[evmChainId]
		if !exist {
			// Iterate through ExtraChains_evm to find a valid ChainDetails if it doesnot already exist
			for _, extra := range obj.ExtraChains_evm {
				if extra.ChainDetails.ChainSelector != 0 {
					return extra.ChainDetails, nil // Return first valid ExtraChain found
				}
			}
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilySolana:
		details, exist := solanaChainIdToChainSelector[chainID]
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
	case FamilyTron:
		tronChainId, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := tronSelectorsMap[tronChainId]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	default:
		return ChainDetails{}, fmt.Errorf("family %s is not yet support", family)
	}
}
