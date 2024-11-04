package chain_selectors

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains.go

//go:embed test_selectors_restructured.yml
var testSelectorsYml []byte

//go:embed selector_restructured.yml
var selectorYml []byte

type ChainDetails struct {
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
var testChainIDToSelectorMapForFamily = make(map[string]map[string]uint64)
var selectorToChainDetails = loadChainDetailsBySelector()
var testSelectorToChainDetailsMap = loadTestChainDetailsBySelector()

func loadTestChainDetailsBySelector() map[uint64]ChainDetails {
	type yamlData struct {
		Selectors map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data yamlData
	err := yaml.Unmarshal(testSelectorsYml, &data)
	if err != nil {
		panic(err)
	}

	for k, v := range data.Selectors {
		if v.Family == "" {
			continue
		}

		// update testChainIDToSelectorMapForFamily
		_, exist := testChainIDToSelectorMapForFamily[v.Family]
		if exist {
			testChainIDToSelectorMapForFamily[v.Family][v.ChainID] = k
		} else {
			testChainIDToSelectorMapForFamily[v.Family] = make(map[string]uint64)
		}
	}

	return data.Selectors
}

func loadChainDetailsBySelector() map[uint64]ChainDetails {
	type yamlData struct {
		Selectors map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data yamlData
	err := yaml.Unmarshal(selectorYml, &data)
	if err != nil {
		panic(err)
	}

	for k, v := range data.Selectors {
		if v.Family == "" {
			continue
		}

		// update chainIDToSelectorMapForFamily
		_, exist := chainIDToSelectorMapForFamily[v.Family]
		if !exist {
			chainIDToSelectorMapForFamily[v.Family] = make(map[string]uint64)
		}

		chainIDToSelectorMapForFamily[v.Family][v.ChainID] = k
	}

	return data.Selectors
}

func ChainSelectorToChainDetails() map[uint64]ChainDetails {
	copyMap := make(map[uint64]ChainDetails, len(selectorToChainDetails))
	for k, v := range selectorToChainDetails {
		copyMap[k] = v
	}

	return copyMap
}

func TestChainSelectorToChainDetails() map[uint64]ChainDetails {
	copyMap := make(map[uint64]ChainDetails, len(selectorToChainDetails))
	for k, v := range testSelectorToChainDetailsMap {
		copyMap[k] = v
	}

	return copyMap
}

func GetSelectorFamily(selector uint64) (string, error) {
	// previously selector_families.yml includes both real and test chains, therefore we check both maps
	details, exist := selectorToChainDetails[selector]
	if exist {
		return details.Family, nil

	}

	details, exist = testSelectorToChainDetailsMap[selector]
	if exist {
		return details.Family, nil
	}

	return "", fmt.Errorf("chain detail not found for selector %d", selector)
}

func ChainIdFromSelector(chainSelectorId uint64) (string, error) {
	chainDetail, ok := selectorToChainDetails[chainSelectorId]
	if ok {
		return chainDetail.ChainID, nil
	}

	chainDetail, ok = testSelectorToChainDetailsMap[chainSelectorId]
	if ok {
		return chainDetail.ChainID, nil
	}

	return "0", fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
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
		return 0, fmt.Errorf("chain selector not found for chainID %v, family %v", chainId, family)
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
		return "", fmt.Errorf("chain family not found for chain %v, family %v", chainId, family)
	}

	selector, exist := selectorMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain selector not found for chain %v, family %v", chainId, family)
	}

	details, exist := selectorToChainDetails[selector]
	if !exist {
		return "", fmt.Errorf("chain details not found for chain %v, family %v", chainId, family)
	}

	// when name is missing use chainID
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

	return "0", fmt.Errorf("chain not found for name %s and family %s", name, family)
}

func TestChainIds() []uint64 {
	chainIds := make([]uint64, 0, len(testSelectorToChainDetailsMap))
	for k := range testSelectorToChainDetailsMap {
		chainIds = append(chainIds, k)
	}
	return chainIds
}

var chainsBySelector = make(map[uint64]Chain)
var chainsByEvmChainID = make(map[string]Chain)

func init() {
	for _, ch := range ALL {
		chainsBySelector[ch.Selector] = ch
		if ch.Family == FamilyEVM {
			chainsByEvmChainID[ch.ChainID] = ch
		}
	}
}

func ChainByEvmChainID(evmChainID string) (Chain, bool) {
	ch, exists := chainsByEvmChainID[evmChainID]
	return ch, exists
}

func ChainBySelector(sel uint64) (Chain, bool) {
	ch, exists := chainsBySelector[sel]
	return ch, exists
}

func IsEvm(chainSel uint64) (bool, error) {
	chain, exists := ChainBySelector(chainSel)
	if !exists {
		return false, fmt.Errorf("chain %d not found", chainSel)
	}

	if chain.Family == FamilyEVM {
		return true, nil
	}
	return false, nil
}
