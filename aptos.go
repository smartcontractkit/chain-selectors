package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_aptos.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_aptos.yml
var aptosSelectorsYml []byte

var (
	aptosSelectorsMap     = parseAptosYml(aptosSelectorsYml)
	aptosChainsBySelector = make(map[uint64]AptosChain)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Aptos {
		if _, exists := aptosSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %d because it already exists", chainID)
			continue
		}
		aptosSelectorsMap[chainID] = chainDetails
		aptosChainsBySelector[chainDetails.ChainSelector] = AptosChain{
			ChainID:  chainID,
			Selector: chainDetails.ChainSelector,
			Name:     chainDetails.ChainName,
		}
	}

	for _, v := range AptosALL {
		aptosChainsBySelector[v.Selector] = v
	}
}

func parseAptosYml(ymlFile []byte) map[uint64]ChainDetails {
	type ymlData struct {
		SelectorsByAptosChainId map[uint64]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	err = validateAptosChainID(data.SelectorsByAptosChainId)
	if err != nil {
		panic(err)
	}

	return data.SelectorsByAptosChainId
}

func validateAptosChainID(data map[uint64]ChainDetails) error {
	// TODO: https://smartcontract-it.atlassian.net/browse/NONEVM-890
	return nil
}

func AptosChainIdToChainSelector() map[uint64]uint64 {
	copyMap := make(map[uint64]uint64, len(aptosSelectorsMap))
	for k, v := range aptosSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func AptosNameFromChainId(chainId uint64) (string, error) {
	details, exist := aptosSelectorsMap[chainId]
	if !exist {
		// Try remote datasource if enabled
		if remoteDetails, ok := getRemoteChainByID(FamilyAptos, fmt.Sprint(chainId)); ok {
			if remoteDetails.ChainName == "" {
				return fmt.Sprint(chainId), nil
			}
			return remoteDetails.ChainName, nil
		}
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return fmt.Sprint(chainId), nil
	}
	return details.ChainName, nil
}

func AptosChainIdFromSelector(selector uint64) (uint64, error) {
	chain, exist := aptosChainsBySelector[selector]
	if !exist {
		// Try remote datasource if enabled (selectors are globally unique)
		if _, chainID, _, ok := getRemoteChainBySelector(selector); ok {
			id, err := strconv.ParseUint(chainID, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid chain id from remote datasource for selector %d: %w", selector, err)
			}
			return id, nil
		}
		return 0, fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chain.ChainID, nil
}

func AptosChainBySelector(selector uint64) (AptosChain, bool) {
	chain, exist := aptosChainsBySelector[selector]
	if exist {
		return chain, true
	}
	// Try remote datasource if enabled (selectors are globally unique)
	if _, chainID, details, ok := getRemoteChainBySelector(selector); ok {
		id, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return AptosChain{}, false
		}
		return AptosChain{
			ChainID:  id,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}, true
	}
	return AptosChain{}, false
}
