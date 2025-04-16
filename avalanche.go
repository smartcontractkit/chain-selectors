package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_avalanche.go

//go:embed selectors_avalanche.yml
var avalancheSelectorsYml []byte

var (
	avalancheSelectorsMap      = parseAvalancheYml(avalancheSelectorsYml)
	avalancheChainIdBySelector = make(map[uint64]string)
)

func init() {
	for k, v := range avalancheSelectorsMap {
		avalancheChainIdBySelector[v.ChainSelector] = k
	}
}

func parseAvalancheYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByAvalancheChainId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByAvalancheChainId
}

func AvalancheChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(avalancheSelectorsMap))
	for k, v := range avalancheSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func AvalancheNameFromChainId(chainId string) (string, error) {
	details, exist := avalancheSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func AvalancheChainIdFromSelector(selector uint64) (string, error) {
	chainId, exist := avalancheChainIdBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}
