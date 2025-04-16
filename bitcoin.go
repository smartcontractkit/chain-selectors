package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_bitcoin.go

//go:embed selectors_bitcoin.yml
var bitcoinSelectorsYml []byte

var (
	bitcoinSelectorsMap      = parseBitcoinYml(bitcoinSelectorsYml)
	bitcoinChainIdBySelector = make(map[uint64]string)
)

func init() {
	for k, v := range bitcoinSelectorsMap {
		bitcoinChainIdBySelector[v.ChainSelector] = k
	}
}

func parseBitcoinYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByBitcoinChainId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByBitcoinChainId
}

func BitcoinChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(bitcoinSelectorsMap))
	for k, v := range bitcoinSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func BitcoinNameFromChainId(chainId string) (string, error) {
	details, exist := bitcoinSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func BitcoinChainIdFromSelector(selector uint64) (string, error) {
	chainId, exist := bitcoinChainIdBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}
