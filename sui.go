package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_sui.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_sui.yml
var suiSelectorsYml []byte

var (
	suiSelectorsMap     = parseSuiYml(suiSelectorsYml)
	suiChainsBySelector = make(map[uint64]SuiChain)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Sui {
		if _, exists := suiSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %d because it already exists", chainID)
			continue
		}
		suiSelectorsMap[chainID] = chainDetails
		suiChainsBySelector[chainDetails.ChainSelector] = SuiChain{
			ChainID:  chainID,
			Selector: chainDetails.ChainSelector,
			Name:     chainDetails.ChainName,
		}
	}

	for _, v := range SuiALL {
		suiChainsBySelector[v.Selector] = v
	}
}

func parseSuiYml(ymlFile []byte) map[uint64]ChainDetails {
	type ymlData struct {
		SelectorsBySuiChainId map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	err = validateSuiChainID(data.SelectorsBySuiChainId)
	if err != nil {
		panic(err)
	}

	return data.SelectorsBySuiChainId
}

func validateSuiChainID(data map[uint64]ChainDetails) error {
	// TODO: https://smartcontract-it.atlassian.net/browse/NONEVM-890
    return nil
}

func SuiChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(suiSelectorsMap))
	for k, v := range suiSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func SuiNameFromChainId(chainId uint64) (string, error) {
	details, exist := suiSelectorsMap[chainId]
	if !exist {
		// Try remote datasource if enabled
		if remoteDetails, ok := getRemoteChainByID(FamilySui, fmt.Sprint(chainId)); ok {
			if remoteDetails.ChainName == "" {
				return fmt.Sprint(chainId), nil
			}
			return remoteDetails.ChainName, nil
		}
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func SuiChainIdFromSelector(selector uint64) (uint64, error) {
	chain, exist := suiChainsBySelector[selector]
	if !exist {
		// Try remote datasource if enabled (selectors are globally unique)
		if _, chainID, _, ok := getRemoteChainBySelector(selector); ok {
			id, err := strconv.ParseUint(chainID, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid chain id from remote datasource for selector %d: %w", selector, err)
			}
			return id, nil
		}
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chain.ChainID, nil
}

func SuiChainBySelector(selector uint64) (SuiChain, bool) {
	chain, exist := suiChainsBySelector[selector]
	if exist {
		return chain, true
	}
	// Try remote datasource if enabled (selectors are globally unique)
	if _, chainID, details, ok := getRemoteChainBySelector(selector); ok {
		id, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return SuiChain{}, false
		}
		return SuiChain{
			ChainID:  id,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}, true
	}
	return SuiChain{}, false
}
