package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_beacon.go

//go:embed selectors_beacon.yml
var beaconSelectorsYml []byte

var (
	beaconSelectorsMap      = parseBeaconYml(beaconSelectorsYml)
	beaconChainIdBySelector = make(map[uint64]uint64)
)

func init() {
	for k, v := range beaconSelectorsMap {
		beaconChainIdBySelector[v.ChainSelector] = k
	}
}

func parseBeaconYml(ymlFile []byte) map[uint64]ChainDetails {
	type ymlData struct {
		SelectorsByBeaconChainId map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByBeaconChainId
}

func BeaconChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(beaconSelectorsMap))
	for k, v := range beaconSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func BeaconNameFromChainId(chainId uint64) (string, error) {
	details, exist := beaconSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func BeaconChainIdFromSelector(selector uint64) (uint64, error) {
	chainId, exist := beaconChainIdBySelector[selector]
	if !exist {
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}
