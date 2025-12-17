package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_canton.go

//go:embed selectors_canton.yml
var cantonSelectorsYml []byte

var (
	cantonChainsByChainId  = parseCantonYml(cantonSelectorsYml)
	cantonChainsBySelector = loadAllCantonSelectors(cantonChainsByChainId)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Canton {
		if _, exists := cantonChainsByChainId[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for Canton chain %s because it already exists", chainID)
			continue
		}
		cantonChainsByChainId[chainID] = chainDetails
		cantonChainsBySelector[chainDetails.ChainSelector] = CantonChain{
			ChainID:  chainID,
			Selector: chainDetails.ChainSelector,
			Name:     chainDetails.ChainName,
		}
	}

	err := validateCantonChainID(cantonChainsByChainId)
	if err != nil {
		panic(err)
	}
}

func parseCantonYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByName map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	if err := yaml.Unmarshal(ymlFile, &data); err != nil {
		panic(err)
	}

	return data.SelectorsByName
}

func loadAllCantonSelectors(in map[string]ChainDetails) map[uint64]CantonChain {
	output := make(map[uint64]CantonChain, len(cantonChainsByChainId))
	for chainID, v := range in {
		output[v.ChainSelector] = CantonChain{
			ChainID:  chainID,
			Selector: v.ChainSelector,
			Name:     v.ChainName,
		}
	}
	return output
}

func validateCantonChainID(data map[string]ChainDetails) error {
	// Add validation logic if needed
	return nil
}

func CantonChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(cantonChainsByChainId))
	for k, v := range cantonChainsByChainId {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func CantonNameFromChainId(chainID string) (string, error) {
	details, exist := cantonChainsByChainId[chainID]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain: %v", chainID)
	}
	if details.ChainName == "" {
		return "", fmt.Errorf("chain name is empty for chain: %v", chainID)
	}
	return details.ChainName, nil
}

func CantonChainIdFromSelector(selector uint64) (string, error) {
	chain, exist := cantonChainsBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain not found for chain selector %d", selector)
	}

	return chain.ChainID, nil
}

func CantonChainBySelector(selector uint64) (CantonChain, bool) {
	chain, exists := cantonChainsBySelector[selector]
	return chain, exists
}
