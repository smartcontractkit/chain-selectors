package chain_selectors

import (
	_ "embed"
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

//go:generate go run genchains_stellar.go
//go:generate go run generate_all_selectors.go

//go:embed selectors_stellar.yml
var stellarSelectorsYml []byte

var (
	stellarChainsByChainId  = parseStellarYml(stellarSelectorsYml)
	stellarChainsBySelector = loadAllStellarSelectors(stellarChainsByChainId)
	// Network passphrases are Stellar-specific (not part of the shared ChainDetails),
	// so they are parsed separately. A Stellar network is defined by its passphrase and
	// its network ID is SHA-256(passphrase); tx signing requires the passphrase string,
	// which cannot be derived from the chain ID (the hash), hence it lives here.
	stellarChainIdToPassphrase = parseStellarPassphrases(stellarSelectorsYml)
)

func init() {
	// Load extra selectors
	for chainID, chainDetails := range getExtraSelectors().Stellar {
		if _, exists := stellarChainsByChainId[chainID]; exists {
			log.Printf("WARN: Skipping extra selector for Stellar chain %s because it already exists", chainID)
			continue
		}
		stellarChainsByChainId[chainID] = chainDetails
		stellarChainsBySelector[chainDetails.ChainSelector] = StellarChain{
			ChainID:     chainID,
			Selector:    chainDetails.ChainSelector,
			Name:        chainDetails.ChainName,
			NetworkType: chainDetails.NetworkType,
			Passphrase:  stellarChainIdToPassphrase[chainID],
		}
	}

	err := validateStellarChainID(stellarChainsByChainId)
	if err != nil {
		panic(err)
	}
}

func parseStellarYml(ymlFile []byte) map[string]ChainDetails {
	type ymlData struct {
		SelectorsByNetworkId map[string]ChainDetails `yaml:"selectors"`
	}

	var data ymlData
	if err := yaml.Unmarshal(ymlFile, &data); err != nil {
		panic(err)
	}

	return data.SelectorsByNetworkId
}

func parseStellarPassphrases(ymlFile []byte) map[string]string {
	type stellarDetails struct {
		Passphrase string `yaml:"passphrase"`
	}
	type ymlData struct {
		SelectorsByNetworkId map[string]stellarDetails `yaml:"selectors"`
	}

	var data ymlData
	if err := yaml.Unmarshal(ymlFile, &data); err != nil {
		panic(err)
	}

	out := make(map[string]string, len(data.SelectorsByNetworkId))
	for chainID, v := range data.SelectorsByNetworkId {
		out[chainID] = v.Passphrase
	}
	return out
}

func loadAllStellarSelectors(in map[string]ChainDetails) map[uint64]StellarChain {
	output := make(map[uint64]StellarChain, len(stellarChainsByChainId))
	for chainID, v := range in {
		output[v.ChainSelector] = StellarChain{
			ChainID:     chainID,
			Selector:    v.ChainSelector,
			Name:        v.ChainName,
			NetworkType: v.NetworkType,
			Passphrase:  stellarChainIdToPassphrase[chainID],
		}
	}
	return output
}

func validateStellarChainID(data map[string]ChainDetails) error {
	// Chain IDs are SHA-256 hashes of network passphrases
	// Add validation logic if needed
	return nil
}

func StellarChainIdToChainSelector() map[string]uint64 {
	copyMap := make(map[string]uint64, len(stellarChainsByChainId))
	for k, v := range stellarChainsByChainId {
		copyMap[k] = v.ChainSelector
	}
	return copyMap
}

func StellarNameFromChainId(chainID string) (string, error) {
	details, exist := stellarChainsByChainId[chainID]
	if !exist {
		return "", fmt.Errorf("chain name not found for chain: %v", chainID)
	}
	if details.ChainName == "" {
		return "", fmt.Errorf("chain name is empty for chain: %v", chainID)
	}
	return details.ChainName, nil
}

func StellarChainIdFromSelector(selector uint64) (string, error) {
	chain, exist := stellarChainsBySelector[selector]
	if !exist {
		return "", fmt.Errorf("chain not found for chain selector %d", selector)
	}

	return chain.ChainID, nil
}

func StellarChainBySelector(selector uint64) (StellarChain, bool) {
	chain, exists := stellarChainsBySelector[selector]
	return chain, exists
}

func StellarNetworkTypeFromChainId(chainId string) (NetworkType, error) {
	if chainDetails, exist := stellarChainsByChainId[chainId]; exist {
		return chainDetails.NetworkType, nil
	}
	return "", fmt.Errorf("chain network type not found for chain %v", chainId)
}

// StellarPassphraseFromChainId returns the network passphrase for a Stellar chain
// ID (network ID). The network ID is SHA-256(passphrase); signing requires the
// passphrase string, which cannot be derived from the ID.
func StellarPassphraseFromChainId(chainID string) (string, error) {
	passphrase, exist := stellarChainIdToPassphrase[chainID]
	if !exist || passphrase == "" {
		return "", fmt.Errorf("network passphrase not found for chain: %v", chainID)
	}
	return passphrase, nil
}
