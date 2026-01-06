package chain_selectors

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// ExtraSelectorsData is a format expected when loading extra selectors from a YAML file.
type ExtraSelectorsData struct {
	Evm      map[uint64]ChainDetails `yaml:"evm,omitempty"`
	Aptos    map[uint64]ChainDetails `yaml:"aptos,omitempty"`
	Solana   map[string]ChainDetails `yaml:"solana,omitempty"`
	Sui      map[uint64]ChainDetails `yaml:"sui,omitempty"`
	Ton      map[int32]ChainDetails  `yaml:"ton,omitempty"`
	Tron     map[uint64]ChainDetails `yaml:"tron,omitempty"`
	Starknet map[string]ChainDetails `yaml:"starknet,omitempty"`
	Canton   map[string]ChainDetails `yaml:"canton,omitempty"`
}

var (
	extraSelectors       ExtraSelectorsData
	extraSelectorsLoaded bool
)

func loadAndParseExtraSelectors() (result ExtraSelectorsData) {
	extraSelectorsFile := os.Getenv("EXTRA_SELECTORS_FILE")
	if extraSelectorsFile == "" {
		return
	}

	fileContent, err := os.ReadFile(extraSelectorsFile)
	if err != nil {
		log.Printf("Error reading extra selectors file %s: %v", extraSelectorsFile, err)
		panic(err)
	}

	var data ExtraSelectorsData
	err = yaml.Unmarshal(fileContent, &data)
	if err != nil {
		log.Printf("Error unmarshaling extra selectors YAML: %v", err)
		panic(err)
	}

	// Validate individual chain formats
	err = validateSolanaChainID(data.Solana)
	if err != nil {
		log.Println(data.Solana)
		log.Printf("Error parsing extra selectors for Solana: %v", err)
		panic(err)
	}

	err = validateSuiChainID(data.Sui)
	if err != nil {
		log.Printf("Error parsing extra selectors for Sui: %v", err)
		panic(err)
	}
	err = validateAptosChainID(data.Aptos)
	if err != nil {
		log.Printf("Error parsing extra selectors for Aptos: %v", err)
		panic(err)
	}

	if err := validateCantonChainID(data.Canton); err != nil {
		log.Printf("Error parsing extra selectors for Canton: %v", err)
		panic(err)
	}

	log.Printf("Successfully loaded extra selectors from %s", extraSelectorsFile)
	return data
}

func getExtraSelectors() ExtraSelectorsData {
	if !extraSelectorsLoaded {
		extraSelectors = loadAndParseExtraSelectors()
		extraSelectorsLoaded = true
	}
	return extraSelectors
}
