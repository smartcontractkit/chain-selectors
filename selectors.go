package chain_selectors

import (
	_ "embed"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains.go

//go:embed selectors.yml
var selectorsYml []byte

//go:embed test_selectors.yml
var testSelectorsYml []byte

//go:embed selector_families.yml
var selectorFamiliesYml []byte

type chainDetails struct {
	ChainSelector uint64 `yaml:"selector"`
	ChainName     string `yaml:"name"`
}

const (
	FamilyEVM      = "evm"
	FamilySolana   = "solana"
	FamilyStarknet = "starknet"
	FamilyCosmos   = "cosmos"
	FamilyAptos    = "aptos"
)

var selectorsMap = parseYml(selectorsYml)
var testSelectorsMap = parseYml(testSelectorsYml)

var evmChainIdToChainSelector = loadAllSelectors()
var selectorToChainFamily = loadSelectorToFamilyMap()

func loadAllSelectors() map[uint64]chainDetails {
	output := make(map[uint64]chainDetails, len(selectorsMap)+len(testSelectorsMap))
	for k, v := range selectorsMap {
		output[k] = v
	}
	for k, v := range testSelectorsMap {
		output[k] = v
	}
	return output
}

func loadSelectorToFamilyMap() map[uint64]string {
	type familyDetails struct {
		Family string `yaml:"family"`
		Name   string `yaml:"name"`
	}

	type yamlData struct {
		SelectorFamilies map[uint64]familyDetails `yaml:"selector_families"`
	}

	var data yamlData
	err := yaml.Unmarshal(selectorFamiliesYml, &data)
	if err != nil {
		panic(err)
	}

	var selectorFamilies = make(map[uint64]string, len(data.SelectorFamilies))
	for k, v := range data.SelectorFamilies {
		selectorFamilies[k] = v.Family
	}

	return selectorFamilies
}

func parseYml(ymlFile []byte) map[uint64]chainDetails {
	type ymlData struct {
		Selectors map[uint64]chainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.Selectors
}

func GetSelectorFamily(selector uint64) (string, error) {
	family, exist := selectorToChainFamily[selector]
	if !exist {
		return "", fmt.Errorf("family not found for selector %d", selector)
	}

	return family, nil
}

func EvmChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(evmChainIdToChainSelector))
	for k, v := range evmChainIdToChainSelector {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func ChainIdFromSelector(chainSelectorId uint64) (uint64, error) {
	for k, v := range evmChainIdToChainSelector {
		if v.ChainSelector == chainSelectorId {
			return k, nil
		}
	}
	return 0, fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

func SelectorFromChainId(chainId uint64) (uint64, error) {
	if chainSelectorId, exist := evmChainIdToChainSelector[chainId]; exist {
		return chainSelectorId.ChainSelector, nil
	}
	return 0, fmt.Errorf("chain selector not found for chain %d", chainId)
}

func NameFromChainId(chainId uint64) (string, error) {
	details, exist := evmChainIdToChainSelector[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %d", chainId)
	}
	if details.ChainName == "" {
		return strconv.FormatUint(chainId, 10), nil
	}
	return details.ChainName, nil
}

func ChainIdFromName(name string) (uint64, error) {
	for k, v := range evmChainIdToChainSelector {
		if v.ChainName == name {
			return k, nil
		}
	}
	chainId, err := strconv.ParseUint(name, 10, 64)
	if err == nil {
		if details, exist := evmChainIdToChainSelector[chainId]; exist && details.ChainName == "" {
			return chainId, nil
		}
	}
	return 0, fmt.Errorf("chain not found for name %s", name)
}

func TestChainIds() []uint64 {
	chainIds := make([]uint64, 0, len(testSelectorsMap))
	for k := range testSelectorsMap {
		chainIds = append(chainIds, k)
	}
	return chainIds
}

var chainsBySelector = make(map[uint64]Chain)
var chainsByEvmChainID = make(map[uint64]Chain)

func init() {
	for _, ch := range ALL {
		chainsBySelector[ch.Selector] = ch
		chainsByEvmChainID[ch.EvmChainID] = ch
	}
}

func ChainBySelector(sel uint64) (Chain, bool) {
	ch, exists := chainsBySelector[sel]
	return ch, exists
}

func ChainByEvmChainID(evmChainID uint64) (Chain, bool) {
	ch, exists := chainsByEvmChainID[evmChainID]
	return ch, exists
}

func IsEvm(chainSel uint64) (bool, error) {
	_, exists := ChainBySelector(chainSel)
	if !exists {
		return false, fmt.Errorf("chain %d not found", chainSel)
	}
	// We always return true since only evm chains are supported atm.
	return true, nil
}
