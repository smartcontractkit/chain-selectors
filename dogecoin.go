package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_dogecoin.go

//go:embed selectors_dogecoin.yml
var dogecoinSelectorsYml []byte

var (
	dogecoinSelectorsMap      = parseDogecoinYml(dogecoinSelectorsYml)
	dogecoinChainIdBySelector = make(map[uint64]string)
)

func init() {
	for k, v := range dogecoinSelectorsMap {
		dogecoinChainIdBySelector[v.ChainSelector] = k
	}
}

func parseDogecoinYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByDogecoinChainId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByDogecoinChainId
}

func DogecoinChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(dogecoinSelectorsMap))
	for k, v := range dogecoinSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func DogecoinNameFromChainId(chainId string) (string, error) {
	details, exist := dogecoinSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func DogecoinChainIdFromSelector(selector uint64) (string, error) {
	chainId, exist := dogecoinChainIdBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}
