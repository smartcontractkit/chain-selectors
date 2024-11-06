package chain_selectors

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

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

func parseSolanaYml(ymlFile []byte) map[string]chainDetails {
	type ymlData struct {
		SelectorsBySolanaChainId map[string]chainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}
	return data.SelectorsBySolanaChainId
}
