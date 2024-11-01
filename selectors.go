package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains.go

//go:embed test_selectors.yml
var testSelectorsYml []byte

//go:embed selector_restructured.yml
var selectorYml []byte

type chainDetails struct {
	Family  string `yaml:"family"`
	Name    string `yaml:"name"`
	ChainID string `yaml:"chain_id"`
}

const (
	FamilyEVM      = "evm"
	FamilySolana   = "solana"
	FamilyStarknet = "starknet"
	FamilyCosmos   = "cosmos"
	FamilyAptos    = "aptos"
)

var chainIDToSelectorMapForFamily = make(map[string]map[string]uint64)
var selectorToChainDetails = loadChainDetailsBySelector()
var testSelectorsMap = loadTestChains()

func loadTestChains() map[uint64]chainDetails {
	type yamlData struct {
		SelectorFamilies map[uint64]chainDetails `yaml:"selectors"`
	}

	var data yamlData
	err := yaml.Unmarshal(testSelectorsYml, &data)
	if err != nil {
		panic(err)
	}
	return data.SelectorFamilies
}

func loadChainDetailsBySelector() map[uint64]chainDetails {
	type yamlData struct {
		SelectorFamilies map[uint64]chainDetails `yaml:"selectors"`
	}

	var data yamlData
	err := yaml.Unmarshal(selectorYml, &data)
	if err != nil {
		panic(err)
	}

	for k, v := range data.SelectorFamilies {
		if v.Family == "" {
			continue
		}

		// update chainIDToSelectorMapForFamily
		_, exist := chainIDToSelectorMapForFamily[v.Family]
		if exist {
			chainIDToSelectorMapForFamily[v.Family][v.ChainID] = k
		} else {
			chainIDToSelectorMapForFamily[v.Family] = make(map[string]uint64)
		}
	}

	return data.SelectorFamilies
}

func GetSelectorFamily(selector uint64) (string, error) {
	details, exist := selectorToChainDetails[selector]
	if !exist {
		return "", fmt.Errorf("chain detail not found for selector %d", selector)
	}

	return details.Family, nil
}

func ChainSelectorToChainDetails() map[uint64]chainDetails {
	copyMap := make(map[uint64]chainDetails, len(selectorToChainDetails))
	for k, v := range selectorToChainDetails {
		copyMap[k] = v
	}

	return copyMap
}

func ChainIdFromSelector(chainSelectorId uint64) (string, error) {
	chainDetail, ok := selectorToChainDetails[chainSelectorId]
	if !ok {
		return "0", fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
	}
	return chainDetail.ChainID, nil
}

func SelectorFromChainIdAndFamily(chainId string, family string) (uint64, error) {
	// if family is missing use EVM as default
	if family == "" {
		family = FamilyEVM
	}

	selectorMap, exist := chainIDToSelectorMapForFamily[family]
	if !exist {
		return 0, fmt.Errorf("chain selector map not found for family %v", family)
	}

	selector, exist := selectorMap[chainId]
	if !exist {
		return 0, fmt.Errorf("chain selector not found for chain %v", chainId)
	}
	return selector, nil
}

func NameFromChainIdAndFamily(chainId string, family string) (string, error) {
	// if family is missing use EVM as default
	if family == "" {
		family = FamilyEVM
	}

	selectorMap, exist := chainIDToSelectorMapForFamily[family]
	if !exist {
		return "", fmt.Errorf("chain family not found for chain %v, %v", chainId, family)
	}

	selector, exist := selectorMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain selector not found for chain %v, %v", chainId, family)
	}

	details, exist := selectorToChainDetails[selector]
	if !exist {
		return "", fmt.Errorf("chain details not found for chain %v, %v", chainId, family)
	}

	if details.Name == "" {
		return chainId, nil
	}
	return details.Name, nil
}

func ChainIdFromNameAndFamily(name string, family string) (string, error) {
	// if family is missing use EVM as default
	if family == "" {
		family = FamilyEVM
	}

	for _, v := range selectorToChainDetails {
		if v.Name == name && family == v.Family {
			return v.ChainID, nil
		}
	}

	return "0", fmt.Errorf("chain not found for name %s", name)
}

func TestChainIds() []uint64 {
	chainIds := make([]uint64, 0, len(testSelectorsMap))
	for k := range testSelectorsMap {
		chainIds = append(chainIds, k)
	}
	return chainIds
}

var chainsBySelector = make(map[uint64]Chain)
var chainsByEvmChainID = make(map[string]Chain)

func init() {
	for _, ch := range ALL {
		chainsBySelector[ch.Selector] = ch
		chainsByEvmChainID[ch.ChainID] = ch
	}
}

func ChainBySelector(sel uint64) (Chain, bool) {
	ch, exists := chainsBySelector[sel]
	return ch, exists
}

func ChainByEvmChainID(chainID string) (Chain, bool) {
	ch, exists := chainsByEvmChainID[chainID]
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
