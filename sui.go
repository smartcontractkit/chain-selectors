package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_sui.go

//go:embed selectors_sui.yml
var suiSelectorsYml []byte

var (
	suiSelectorsMap     = parseSuiYml(suiSelectorsYml)
	suiChainsBySelector = make(map[uint64]SuiChain)
)

func init() {
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

	validateSuiChainID(data.SelectorsBySuiChainId)
	return data.SelectorsBySuiChainId
}

func validateSuiChainID(data map[uint64]ChainDetails) {
	// TODO: https://smartcontract-it.atlassian.net/browse/NONEVM-890
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
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chain.ChainID, nil
}

func SuiChainBySelector(selector uint64) (SuiChain, bool) {
	chain, exist := suiChainsBySelector[selector]
	return chain, exist
}
