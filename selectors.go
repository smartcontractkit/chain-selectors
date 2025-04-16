package chain_selectors

import (
	"fmt"
	"strconv"
)

const (
	FamilyEVM       = "evm"
	FamilySolana    = "solana"
	FamilyStarknet  = "starknet"
	FamilyCosmos    = "cosmos"
	FamilyAptos     = "aptos"
	FamilyTron      = "tron"
	FamilyAvalanche = "avalanche"
	FamilyBeacon    = "beacon"
	FamilyBitcoin   = "bitcoin"
)

type chainInfo struct {
	Family       string
	ChainID      string
	ChainDetails ChainDetails
}

func getChainInfo(selector uint64) (chainInfo, error) {
	// check EVM
	_, exist := evmChainsBySelector[selector]
	if exist {
		family := FamilyEVM

		evmChainId, err := ChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", family, selector, err)
		}

		details, exist := evmChainIdToChainSelector[evmChainId]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %d for %s", evmChainId, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", evmChainId),
			ChainDetails: details,
		}, nil
	}

	// check solana
	_, exist = solanaChainsBySelector[selector]
	if exist {
		family := FamilySolana

		chainID, err := SolanaChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %s chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := solanaChainIdToChainSelector[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      chainID,
			ChainDetails: details,
		}, nil

	}

	// check avalanche
	_, exist = avalancheChainIdBySelector[selector]
	if exist {
		family := FamilyAvalanche

		chainID, err := AvalancheChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := avalancheSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      chainID,
			ChainDetails: details,
		}, nil
	}

	// check aptos
	_, exist = aptosChainsBySelector[selector]
	if exist {
		family := FamilyAptos

		chainID, err := AptosChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := aptosSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", chainID),
			ChainDetails: details,
		}, nil
	}

	// check beacon
	_, exist = beaconChainIdBySelector[selector]
	if exist {
		family := FamilyBeacon

		chainID, err := BeaconChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := beaconSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", chainID),
			ChainDetails: details,
		}, nil
	}

	// check bitcoin
	_, exist = bitcoinChainIdBySelector[selector]
	if exist {
		family := FamilyBitcoin

		chainID, err := BitcoinChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := bitcoinSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      chainID,
			ChainDetails: details,
		}, nil
	}

	// check tron
	_, exist = tronChainIdBySelector[selector]
	if exist {
		family := FamilyTron

		chainID, err := TronChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := tronSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", chainID),
			ChainDetails: details,
		}, nil
	}

	return chainInfo{}, fmt.Errorf("unknown chain selector %d", selector)
}

func GetSelectorFamily(selector uint64) (string, error) {
	chainInfo, err := getChainInfo(selector)
	if err != nil {
		return "", fmt.Errorf("unknown chain selector %d", selector)
	}

	return chainInfo.Family, nil
}

func GetChainIDFromSelector(selector uint64) (string, error) {
	chainInfo, err := getChainInfo(selector)
	if err != nil {
		return "", fmt.Errorf("unknown chain selector %d", selector)
	}

	return chainInfo.ChainID, nil
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
		details, exist := solanaChainIdToChainSelector[chainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilyAvalanche:
		details, exist := avalancheSelectorsMap[chainID]
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
	case FamilyBeacon:
		beaconChainId, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := beaconSelectorsMap[beaconChainId]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilyBitcoin:
		details, exist := bitcoinSelectorsMap[chainID]
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
