package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_solana.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_solana.yml
var solanaSelectorsYml []byte

//go:embed test_selectors_solana.yml
var testSelectorsSolanaYml []byte

var (
	solanaSelectorsMap           = parseSolanaYml(solanaSelectorsYml)
	solanaTestSelectorsMap       = parseSolanaYml(testSelectorsSolanaYml)
	solanaChainIdToChainSelector = loadAllSolanaSelectors()
	solanaChainsBySelector       = make(map[uint64]SolanaChain)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Solana {
		if _, exists := solanaSelectorsMap[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for chain %s because it already exists", chainID)
			continue
		}
		solanaSelectorsMap[chainID] = chainDetails
		solanaChainIdToChainSelector[chainID] = chainDetails
		solanaChainsBySelector[chainDetails.ChainSelector] = SolanaChain{
			ChainID:   chainID,
			Selector:  chainDetails.ChainSelector,
			Name:      chainDetails.ChainName,
			IsMainnet: chainDetails.IsMainnet,
		}
	}

	for _, v := range SolanaALL {
		solanaChainsBySelector[v.Selector] = v
	}
}

func loadAllSolanaSelectors() map[string]ChainDetails {
	output := make(map[string]ChainDetails, len(solanaSelectorsMap)+len(solanaTestSelectorsMap))
	for k, v := range solanaSelectorsMap {
		output[k] = v
	}
	for k, v := range solanaTestSelectorsMap {
		output[k] = v
	}
	return output
}

func parseSolanaYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsBySolanaChainId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	err := yaml.Unmarshal(ymlFile, &data)
	if err != nil {
		panic(err)
	}

	err = validateSolanaChainID(data.SelectorsBySolanaChainId)
	if err != nil {
		panic(err)
	}

	return data.SelectorsBySolanaChainId
}

func validateSolanaChainID(data map[string]ChainDetails) error {
	for genesisHash := range data {
		b, err := base58.Decode(genesisHash)
		if err != nil {
			return fmt.Errorf("failed to decode base58 genesis hash %s: %w", genesisHash, err)
		}
		if len(b) != 32 {
			return fmt.Errorf("decoded genesis hash %s is not 32 bytes long", genesisHash)
		}
	}
	return nil
}

func SolanaChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(solanaChainIdToChainSelector))
	for k, v := range solanaChainIdToChainSelector {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func SolanaNameFromChainId(chainId string) (string, error) {
	details, exist := solanaChainIdToChainSelector[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return chainId, nil
	}
	return details.ChainName, nil
}

func SolanaChainIdFromSelector(selector uint64) (string, error) {
	chain, exist := solanaChainsBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain not found for selector %d", selector)
	}

	return chain.ChainID, nil
}

func SolanaChainBySelector(selector uint64) (SolanaChain, bool) {
	chain, exists := solanaChainsBySelector[selector]

	return chain, exists
}
