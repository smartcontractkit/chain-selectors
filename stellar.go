package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_stellar.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_stellar.yml
var stellarSelectorsYml []byte

var (
	stellarChainsByChainId  = parseStellarYml(stellarSelectorsYml)
	stellarChainsBySelector = loadAllStellarSelectors(stellarChainsByChainId)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Stellar {
		if _, exists := stellarChainsByChainId[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for Stellar chain %s because it already exists", chainID)
			continue
		}
		stellarChainsByChainId[chainID] = chainDetails
		stellarChainsBySelector[chainDetails.ChainSelector] = StellarChain{
			ChainID:     chainID,
			Selector:    chainDetails.ChainSelector,
			Name:        chainDetails.ChainName,
			NetworkType: chainDetails.NetworkType,
		}
	}

	err := validateStellarChainID(stellarChainsByChainId)
	if err != nil {
		panic(err)
	}
}

func parseStellarYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByNetworkId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	if err := yaml.Unmarshal(ymlFile, &data); err != nil {
		panic(err)
	}

	return data.SelectorsByNetworkId
}

func loadAllStellarSelectors(in map[string]ChainDetails) map[uint64]StellarChain {
	output := make(map[uint64]StellarChain, len(stellarChainsByChainId))
	for chainID, v := range in {
		output[v.ChainSelector] = StellarChain{
			ChainID:     chainID,
			Selector:    v.ChainSelector,
			Name:        v.ChainName,
			NetworkType: v.NetworkType,
		}
	}
	return output
}

func validateStellarChainID(data map[string]ChainDetails) error {
	// Chain IDs are SHA-256 hashes of network passphrases
	// Add validation logic if needed
	return nil
}

func StellarChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(stellarChainsByChainId))
	for k, v := range stellarChainsByChainId {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func StellarNameFromChainId(chainID string) (string, error) {
	details, exist := stellarChainsByChainId[chainID]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain: %v", chainID)
	}
	if details.ChainName == "" {
		return "", fmt.Errorf("chain name is empty for chain: %v", chainID)
	}
	return details.ChainName, nil
}

func StellarChainIdFromSelector(selector uint64) (string, error) {
	chain, exist := stellarChainsBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain not found for chain selector %d", selector)
	}

	return chain.ChainID, nil
}

func StellarChainBySelector(selector uint64) (StellarChain, bool) {
	chain, exists := stellarChainsBySelector[selector]
	return chain, exists
}

func StellarNetworkTypeFromChainId(chainId string) (NetworkType, error) {
	if chainDetails, exist := stellarChainsByChainId[chainId]; exist {
		return chainDetails.NetworkType, nil
	}
	return "", fmt.Errorf("chain network type not found for chain %v", chainId)
}
