package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_starknet.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_starknet.yml
var starknetSelectorsYml []byte

var (
	starknetSelectorsMap     = parseStarknetYml(starknetSelectorsYml)
	starknetChainsBySelector = make(map[uint64]StarknetChain)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Starknet {
		if _, exists := starknetSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %s because it already exists", chainID)
			continue
		}
		starknetSelectorsMap[chainID] = chainDetails
		starknetChainsBySelector[chainDetails.ChainSelector] = StarknetChain{
			ChainID:  chainID,
			Selector: chainDetails.ChainSelector,
			Name:     chainDetails.ChainName,
		}
	}

	for _, v := range StarknetALL {
		starknetChainsBySelector[v.Selector] = v
	}
}

func parseStarknetYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByStarknetChainId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByStarknetChainId
}

func StarknetChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(starknetSelectorsMap))
	for k, v := range starknetSelectorsMap {
		copyMap[k] = v.ChainSelector
	}

	return copyMap
}

func StarknetNameFromChainId(chainId string) (string, error) {
	details, exist := starknetSelectorsMap[chainId]
	if !exist {
		// Try remote datasource if enabled
		if remoteDetails, ok := getRemoteChainByID(FamilyStarknet, chainId); ok {
			if remoteDetails.ChainName == "" {
				return chainId, nil
			}
			return remoteDetails.ChainName, nil
		}
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return chainId, nil
	}

	return details.ChainName, nil
}

func StarknetChainIdFromSelector(selector uint64) (string, error) {
	chain, exist := starknetChainsBySelector[selector]
	if !exist {
		// Try remote datasource if enabled (selectors are globally unique)
		if _, chainID, _, ok := getRemoteChainBySelector(selector); ok {
			return chainID, nil
		}
		return "", fmt.Errorf("chain not found for selector %d", selector)
	}

	return chain.ChainID, nil
}

func StarknetChainBySelector(selector uint64) (StarknetChain, bool) {
	chain, exists := starknetChainsBySelector[selector]
	if exists {
		return chain, true
	}
	// Try remote datasource if enabled (selectors are globally unique)
	if _, chainID, details, ok := getRemoteChainBySelector(selector); ok {
		return StarknetChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}, true
	}
	return StarknetChain{}, false
}
