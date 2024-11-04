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
var selectorsMap = loadYML(selectorYml)
var testSelectorsMap = loadYML(testSelectorsYml)
var chainIdToChainSelector = loadAllChainIDToChainSelector()

func loadAllChainIDToChainSelector() map[uint64]ChainDetails {
	output := make(map[uint64]ChainDetails, len(selectorsMap)+len(testSelectorsMap))
	for k, v := range selectorsMap {
		output[k] = v
	}
	for k, v := range testSelectorsMap {
		output[k] = v
	}
	return output
}

func loadYML(yml []byte) map[uint64]ChainDetails {
	type yamlData struct {
		Selectors map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data yamlData
	err := yaml.Unmarshal(yml, &data)
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
	copyMap := make(map[uint64]ChainDetails, len(selectorsMap))
	for k, v := range chainIdToChainSelector {
		copyMap[k] = v
	}

	return copyMap
}

func TestChainSelectorToChainDetails() map[uint64]ChainDetails {
	copyMap := make(map[uint64]ChainDetails, len(selectorsMap))
	for k, v := range testSelectorsMap {
		copyMap[k] = v
	}

	return copyMap
}

func GetSelectorFamily(selector uint64) (string, error) {
	// previously selector_families.yml includes both real and test chains, therefore we check both maps
	details, exist := chainIdToChainSelector[selector]
	if exist {
		return details.Family, nil
	}

	return "", fmt.Errorf("chain detail not found for selector %d", selector)
}

func ChainIdFromSelector(chainSelectorId uint64) (string, error) {
	chainDetail, ok := chainIdToChainSelector[chainSelectorId]
	if ok {
		return chainDetail.ChainID, nil
	}

	return "0", fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

// SelectorFromChainId is for backward compatibility support
func SelectorFromChainId(chainId string) (uint64, error) {
	return SelectorFromChainIdAndFamily(chainId, FamilyEVM)
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

// ChainIdFromName is for backward compatibility support
func ChainIdFromName(name string) (string, error) {
	return ChainIdFromNameAndFamily(name, FamilyEVM)
}

func ChainIdFromNameAndFamily(name string, family string) (string, error) {
	// if family is missing use EVM as default
	if family == "" {
		family = FamilyEVM
	}

	for _, v := range chainIdToChainSelector {
		if v.Name == name && family == v.Family {
			return v.ChainID, nil
		}
	}

	return "0", fmt.Errorf("chain not found for name %s and family %s", name, family)
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
