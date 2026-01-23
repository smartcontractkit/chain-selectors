package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_ton.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_ton.yml
var tonSelectorsYml []byte

var (
	tonSelectorsMap      = parseTonYml(tonSelectorsYml)
	tonChainIdBySelector = make(map[uint64]int32)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Ton {
		if _, exists := tonSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %d because it already exists", chainID)
			continue
		}
		tonSelectorsMap[chainID] = chainDetails
		tonChainIdBySelector[chainDetails.ChainSelector] = chainID
	}

	for k, v := range tonSelectorsMap {
		tonChainIdBySelector[v.ChainSelector] = k
	}
}

func parseTonYml(ymlFile []byte) map[int32]ChainDetails {
	type ymlData struct {
		SelectorsByTonChainId map[int32]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByTonChainId
}

func TonChainIdToChainSelector() map[int32]uint64 {
	copyMap := make(map[int32]uint64, len(tonSelectorsMap))
	for k, v := range tonSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func TonNameFromChainId(chainId int32) (string, error) {
	details, exist := tonSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func TonChainIdFromSelector(selector uint64) (int32, error) {
	chainId, exist := tonChainIdBySelector[selector]
	if !exist {
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}

func TonIsMainnetChain(chainID int32) (bool, error) {
	details, exist := tonSelectorsMap[chainID]
	if !exist {
		return false, fmt.Errorf("chain not found for chain ID: %v", chainID)
	}
	return details.IsMainnet, nil
}
