package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed selectors.yml
var selectorsYml []byte
var evmChainIdToChainSelector = parseYml()

func parseYml() map[uint64]uint64 {
	type ymlData struct {
		Selectors map[uint64]uint64 `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(selectorsYml, &data)
	if err != nil {
		panic(err)
	}

	return data.Selectors
}

func EvmChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(evmChainIdToChainSelector))
	for k, v := range evmChainIdToChainSelector {
		copyMap[k] = v
	}
	return copyMap
}

func ChainIdFromSelector(chainSelectorId uint64) (uint64, error) {
	for k, v := range evmChainIdToChainSelector {
		if v == chainSelectorId {
			return k, nil
		}
	}
	return 0, fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

func SelectorFromChainId(chainId uint64) (uint64, error) {
	if chainSelectorId, exist := evmChainIdToChainSelector[chainId]; exist {
		return chainSelectorId, nil
	}
	return 0, fmt.Errorf("chain selector not found for chain %d", chainId)
}
