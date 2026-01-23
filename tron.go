package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_tron.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_tron.yml
var tronSelectorsYml []byte

var (
	tronSelectorsMap      = parseTronYml(tronSelectorsYml)
	tronChainIdBySelector = make(map[uint64]uint64)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Tron {
		if _, exists := tronSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %d because it already exists", chainID)
			continue
		}
		tronSelectorsMap[chainID] = chainDetails
		tronChainIdBySelector[chainDetails.ChainSelector] = chainID
	}

	for k, v := range tronSelectorsMap {
		tronChainIdBySelector[v.ChainSelector] = k
	}
}

func parseTronYml(ymlFile []byte) map[uint64]ChainDetails {
	type ymlData struct {
		SelectorsByTronChainId map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByTronChainId
}

func TronChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(tronSelectorsMap))
	for k, v := range tronSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func TronNameFromChainId(chainId uint64) (string, error) {
	details, exist := tronSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func TronChainIdFromSelector(selector uint64) (uint64, error) {
	chainId, exist := tronChainIdBySelector[selector]
	if !exist {
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}

func TronIsMainnetChain(chainID uint64) (bool, error) {
	details, exist := tronSelectorsMap[chainID]
	if !exist {
		return false, fmt.Errorf("chain not found for chain ID: %v", chainID)
	}
	return details.IsMainnet, nil
}
