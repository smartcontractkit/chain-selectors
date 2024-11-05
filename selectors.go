package chain_selectors

import (
	_ "embed"
	"fmt"
	"strconv"

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
var chainSelectorToDetails = loadAllChainSelector()

func loadAllChainSelector() map[uint64]ChainDetails {
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
	for k, v := range chainSelectorToDetails {
		copyMap[k] = v
	}

	return copyMap
}

func GetSelectorFamily(selector uint64) (string, error) {
	// previously selector_families.yml includes both real and test chains, therefore we check both maps
	details, exist := chainSelectorToDetails[selector]
	if exist {
		return details.Family, nil
	}

	return "", fmt.Errorf("chain detail not found for selector %d", selector)
}

// ChainIdFromSelector is for backward compatibility support, it used to return uint64 for chainID so we preserve the behavior
// Deprecated: Call GetChainIdFromSelector directly
func ChainIdFromSelector(chainSelectorId uint64) (uint64, error) {
	chainId, err := GetChainIdFromSelector(chainSelectorId)
	if err != nil {
		return 0, err
	}

	parseInt, err := strconv.ParseUint(chainId, 10, 64)
	if err != nil {
		return 0, err
	}
	return parseInt, fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

func GetChainIdFromSelector(chainSelectorId uint64) (string, error) {
	chainDetail, ok := chainSelectorToDetails[chainSelectorId]
	if ok {
		return chainDetail.ChainID, nil
	}

	return "0", fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

// SelectorFromChainId is for backward compatibility support, it used to take uint64 as chainID so we preserve the behavior
// Deprecated: Call SelectorFromChainIdAndFamily directly
func SelectorFromChainId(chainId uint64) (uint64, error) {
	return SelectorFromChainIdAndFamily(strconv.FormatUint(chainId, 10), FamilyEVM)
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

// ChainIdFromName is for backward compatibility support, it used to return uint64 as chain ID so we preserve the behavior
// Deprecated: Call ChainIdFromNameAndFamily directly
func ChainIdFromName(name string) (uint64, error) {
	chainID, err := ChainIdFromNameAndFamily(name, FamilyEVM)
	if err != nil {
		return 0, err
	}

	// assume chain family is evm and chain id can be converted to uint64
	parseInt, err := strconv.ParseUint(chainID, 10, 64)
	if err != nil {
		return 0, err
	}

	return parseInt, nil
}

// NameFromChainId is for backward compatibility support
// Deprecated: Call SelectorFromChainId directly
func NameFromChainId(chainId uint64) (string, error) {
	selector, err := SelectorFromChainIdAndFamily(strconv.FormatUint(chainId, 10), FamilyEVM)
	if err != nil {
		return "", fmt.Errorf("chain name not found for chain %d", chainId)
	}

	details, exist := chainSelectorToDetails[selector]
	if !exist {
		return "", fmt.Errorf("chain selector not found for chain %d", chainId)
	}

	if details.Name == "" {
		return strconv.FormatUint(chainId, 10), nil
	}

	return details.Name, nil
}

func ChainIdFromNameAndFamily(name string, family string) (string, error) {
	// if family is missing use EVM as default
	if family == "" {
		family = FamilyEVM
	}

	for _, v := range chainSelectorToDetails {
		if v.Name == name && family == v.Family {
			return v.ChainID, nil
		}
	}

	return "0", fmt.Errorf("chain not found for name %s and family %s", name, family)
}

// TestChainIds is for backward compatibility support, it used to return uint64 as chain ID so we preserve the behavior
func TestChainIds() []uint64 {
	chainIds := make([]uint64, 0, len(testSelectorsMap))
	for _, details := range testSelectorsMap {
		parseInt, err := strconv.ParseUint(details.ChainID, 10, 64)
		if err != nil {
			continue
		}

		chainIds = append(chainIds, parseInt)
	}
	return chainIds
}

var chainsBySelector = make(map[uint64]Chain)
var chainsByChainID = make(map[string]Chain)
var chainsByEvmChainID = make(map[string]Chain)

func init() {
	for _, ch := range ALL {
		chainsBySelector[ch.Selector] = ch
		chainsByChainID[ch.ChainID] = ch
		if ch.Family == FamilyEVM {
			chainsByEvmChainID[ch.ChainID] = ch
		}
	}
}

func chainByChainID(chainID string) (Chain, bool) {
	ch, exists := chainsByChainID[chainID]
	return ch, exists
}

// Deprecated: Call chainByChainID directly
func ChainByEvmChainID(evmChainID uint64) (Chain, bool) {
	chainID := strconv.FormatUint(evmChainID, 10)
	ch, exists := chainsByEvmChainID[chainID]
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
