package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_solana.go

//go:embed selectors_solana.yml
var solanaSelectorsYml []byte

var (
	solanaSelectorsMap      = parseSolanaYml(solanaSelectorsYml)
	solanaChainIdBySelector = make(map[uint64]string)
)

func init() {
	for k, v := range solanaSelectorsMap {
		solanaChainIdBySelector[v.ChainSelector] = k
	}
}

func parseSolanaYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsBySolanaChainId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}
	return data.SelectorsBySolanaChainId
}

func SolanaChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(solanaSelectorsMap))
	for k, v := range solanaSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func SolanaNameFromChainId(chainId string) (string, error) {
	details, exist := solanaSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return chainId, nil
	}
	return details.ChainName, nil
}
