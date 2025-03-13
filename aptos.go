package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_aptos.go

//go:embed selectors_aptos.yml
var aptosSelectorsYml []byte

var (
	aptosSelectorsMap     = parseAptosYml(aptosSelectorsYml)
	aptosChainsBySelector = make(map[uint64]AptosChain)
)

func init() {
	for _, v := range AptosALL {
		aptosChainsBySelector[v.Selector] = v
	}
}

func parseAptosYml(ymlFile []byte) map[uint64]ChainDetails {
	type ymlData struct {
		SelectorsByAptosChainId map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	validateAptosChainID(data.SelectorsByAptosChainId)
	return data.SelectorsByAptosChainId
}

func validateAptosChainID(data map[uint64]ChainDetails) {
	// TODO: https://smartcontract-it.atlassian.net/browse/NONEVM-890
}

func AptosChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(aptosSelectorsMap))
	for k, v := range aptosSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func AptosNameFromChainId(chainId uint64) (string, error) {
	details, exist := aptosSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func AptosChainIdFromSelector(selector uint64) (uint64, error) {
	chain, exist := aptosChainsBySelector[selector]
	if !exist {
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chain.ChainID, nil
}

func AptosChainBySelector(selector uint64) (AptosChain, bool) {
	chain, exist := aptosChainsBySelector[selector]
	return chain, exist
}
