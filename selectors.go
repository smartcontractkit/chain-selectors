package chain_selectors

import (
	"fmt"
	"regexp"
	"strconv"
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

	// check sui
	_, exist = suiChainsBySelector[selector]
	if exist {
		family := FamilySui

		chainID, err := SuiChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := suiSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return chainInfo{
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

	// check ton
	_, exist = tonChainIdBySelector[selector]
	if exist {
		family := FamilyTon

		chainID, err := TonChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := tonSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %d for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      fmt.Sprintf("%d", chainID),
			ChainDetails: details,
		}, nil
	}

	// check starknet
	_, exist = starknetChainsBySelector[selector]
	if exist {
		family := FamilyStarknet

		chainID, err := StarknetChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := starknetSelectorsMap[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      chainID,
			ChainDetails: details,
		}, nil
	}

	// check canton
	_, exist = cantonChainsBySelector[selector]
	if exist {
		family := FamilyCanton

		chainID, err := CantonChainIdFromSelector(selector)
		if err != nil {
			return chainInfo{}, fmt.Errorf("failed to get %v chain ID from selector %d: %w", chainID, selector, err)
		}

		details, exist := cantonChainsByChainId[chainID]
		if !exist {
			return chainInfo{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return chainInfo{
			Family:       family,
			ChainID:      chainID,
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

func GetChainNameFromSelector(selector uint64) (string, error) {
	chainInfo, err := getChainInfo(selector)
	if err != nil {
		return "", fmt.Errorf("unknown chain selector %d", selector)
	}

	return chainInfo.ChainDetails.ChainName, nil
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
	case FamilySui:
		suiChainId, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := suiSelectorsMap[suiChainId]
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

	case FamilyTon:
		tonChainId, err := strconv.ParseInt(chainID, 10, 32)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}
		details, exist := tonSelectorsMap[int32(tonChainId)]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilyStarknet:
		details, exist := starknetSelectorsMap[chainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	case FamilyCanton:
		details, exist := cantonChainsByChainId[chainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil
	default:
		return ChainDetails{}, fmt.Errorf("family %s is not yet support", family)
	}
}

// ExtractNetworkEnvName returns chain env identifier from the full network name, for e.g. blockchain-mainnet returns mainnet.
func ExtractNetworkEnvName(networkName string) (string, error) {
	// Create a regexp pattern that matches any of the three.
	re := regexp.MustCompile(`(mainnet|testnet|devnet)`)
	envName := re.FindString(networkName)
	if envName == "" {
		return "", fmt.Errorf("failed to extract network env name from : %s", networkName)
	}
	return envName, nil
}
