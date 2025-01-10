package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_tron.go

//go:embed selectors_tron.yml
var tronSelectorsYml []byte

var (
	tronSelectorsMap      = parseTronYml(tronSelectorsYml)
	tronChainIdBySelector = make(map[uint64]uint64)
)

func init() {
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
