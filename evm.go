package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_evm.go
//go:generate go run generate_all_selectors.go

//go:embed selectors.yml
var selectorsYml []byte

//go:embed test_selectors.yml
var testSelectorsYml []byte

var (
	evmSelectorsMap           = parseYml(selectorsYml)
	evmTestSelectorsMap       = parseYml(testSelectorsYml)
	evmChainIdToChainSelector = loadAllEVMSelectors()
	evmChainsBySelector       = make(map[uint64]Chain)
	evmChainsByEvmChainID     = make(map[uint64]Chain)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Evm {
		if _, exists := evmSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %d because it already exists", chainID)
			continue
		}

		evmSelectorsMap[chainID] = chainDetails
		evmChainIdToChainSelector[chainID] = chainDetails
		chain := Chain{
			EvmChainID: chainID,
			Selector:   chainDetails.ChainSelector,
			Name:       chainDetails.ChainName,
		}
		evmChainsBySelector[chainDetails.ChainSelector] = chain
		evmChainsByEvmChainID[chainID] = chain
	}

	for _, ch := range ALL {
		evmChainsBySelector[ch.Selector] = ch
		evmChainsByEvmChainID[ch.EvmChainID] = ch
	}
}

func loadAllEVMSelectors() map[uint64]ChainDetails {
	output := make(map[uint64]ChainDetails, len(evmSelectorsMap)+len(evmTestSelectorsMap))
	for k, v := range evmSelectorsMap {
		output[k] = v
	}
	for k, v := range evmTestSelectorsMap {
		output[k] = v
	}
	return output
}

func parseYml(ymlFile []byte) map[uint64]ChainDetails {
	type ymlData struct {
		SelectorsByEvmChainId map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByEvmChainId
}

func EvmChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(evmChainIdToChainSelector))
	for k, v := range evmChainIdToChainSelector {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

// Deprecated, this only supports EVM chains, use the chain agnostic `GetChainIDFromSelector` instead
func ChainIdFromSelector(chainSelectorId uint64) (uint64, error) {
	for k, v := range evmChainIdToChainSelector {
		if v.ChainSelector == chainSelectorId {
			return k, nil
		}
	}
	// Try remote datasource if enabled (selectors are globally unique)
	if _, chainID, _, ok := getRemoteChainBySelector(chainSelectorId); ok {
		id, _ := strconv.ParseUint(chainID, 10, 64)
		return id, nil
	}
	return 0, fmt.Errorf("chain not found for chain selector %d", chainSelectorId)
}

// Deprecated, this only supports EVM chains, use the chain agnostic `GetChainDetailsByChainIDAndFamily` instead
func SelectorFromChainId(chainId uint64) (uint64, error) {
	if chainSelectorId, exist := evmChainIdToChainSelector[chainId]; exist {
		return chainSelectorId.ChainSelector, nil
	}
	// Try remote datasource if enabled
	if details, ok := getRemoteChainByID(FamilyEVM, strconv.FormatUint(chainId, 10)); ok {
		return details.ChainSelector, nil
	}
	return 0, fmt.Errorf("chain selector not found for chain %d", chainId)
}

// Deprecated, this only supports EVM chains, use the chain agnostic `NameFromChainId` instead
func NameFromChainId(chainId uint64) (string, error) {
	details, exist := evmChainIdToChainSelector[chainId]
	if !exist {
		// Try remote datasource if enabled
		if remoteDetails, ok := getRemoteChainByID(FamilyEVM, strconv.FormatUint(chainId, 10)); ok {
			if remoteDetails.ChainName == "" {
				return strconv.FormatUint(chainId, 10), nil
			}
			return remoteDetails.ChainName, nil
		}
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
	chainIds := make([]uint64, 0, len(evmTestSelectorsMap))
	for k := range evmTestSelectorsMap {
		chainIds = append(chainIds, k)
	}
	return chainIds
}

func ChainBySelector(sel uint64) (Chain, bool) {
	ch, exists := evmChainsBySelector[sel]
	if exists {
		return ch, true
	}
	// Try remote datasource if enabled (selectors are globally unique)
	if _, chainID, details, ok := getRemoteChainBySelector(sel); ok {
		evmChainID, _ := strconv.ParseUint(chainID, 10, 64)
		return Chain{
			EvmChainID: evmChainID,
			Selector:   details.ChainSelector,
			Name:       details.ChainName,
		}, true
	}
	return Chain{}, false
}

func ChainByEvmChainID(evmChainID uint64) (Chain, bool) {
	ch, exists := evmChainsByEvmChainID[evmChainID]
	if exists {
		return ch, true
	}
	// Try remote datasource if enabled
	if details, ok := getRemoteChainByID(FamilyEVM, strconv.FormatUint(evmChainID, 10)); ok {
		return Chain{
			EvmChainID: evmChainID,
			Selector:   details.ChainSelector,
			Name:       details.ChainName,
		}, true
	}
	return Chain{}, false
}

func IsEvm(chainSel uint64) (bool, error) {
	_, exists := ChainBySelector(chainSel)
	if !exists {
		return false, fmt.Errorf("chain %d not found", chainSel)
	}
	// We always return true since only evm chains are supported atm.
	return true, nil
}
