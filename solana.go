package chain_selectors

import (
	_ "embed"
	"fmt"

	"github.com/mr-tron/base58"
	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_solana.go

//go:embed selectors_solana.yml
var solanaSelectorsYml []byte

var (
	solanaSelectorsMap      = parseSolanaYml(solanaSelectorsYml)
	solanaChainIdBySelector = make(map[uint64]string)
)

func init() {
	for k, v := range solanaSelectorsMap {
		solanaChainIdBySelector[v.ChainSelector] = k
	}
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

	validateSolanaChainID(data.SelectorsBySolanaChainId)
	return data.SelectorsBySolanaChainId
}

func validateSolanaChainID(data map[string]ChainDetails) {
	for genesisHash := range data {
		b, err := base58.Decode(genesisHash)
		if err != nil {
			panic(fmt.Errorf("failed to decode base58 genesis hash %s: %w", genesisHash, err))
		}
		if len(b) != 32 {
			panic(fmt.Errorf("decoded genesis hash %s is not 32 bytes long", genesisHash))
		}
	}
}

func SolanaChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(solanaSelectorsMap))
	for k, v := range solanaSelectorsMap {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func SolanaNameFromChainId(chainId string) (string, error) {
	details, exist := solanaSelectorsMap[chainId]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain %v", chainId)
	}
	if details.ChainName == "" {
		return chainId, nil
	}
	return details.ChainName, nil
}

func SolanaChainIdFromSelector(selector uint64) (string, error) {
	chainId, exist := solanaChainIdBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain id not found for selector %d", selector)
	}

	return chainId, nil
}
