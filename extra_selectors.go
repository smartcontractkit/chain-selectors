package chain_selectors

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type extraSelectorsData struct {
	Evm    map[uint64]ChainDetails `yaml:"evm"`
	Aptos  map[uint64]ChainDetails `yaml:"aptos"`
	Solana map[string]ChainDetails `yaml:"solana"`
	Sui    map[uint64]ChainDetails `yaml:"sui"`
	Ton    map[int32]ChainDetails  `yaml:"ton"`
	Tron   map[uint64]ChainDetails `yaml:"tron"`
}

var extraSelectors extraSelectorsData
var extraSelectorsLoaded bool

func loadAndParseExtraSelectors() (result extraSelectorsData) {
	extraSelectorsFile := os.Getenv("EXTRA_SELECTORS_FILE")
	if extraSelectorsFile == "" {
		return
	}

	fileContent, err := os.ReadFile(extraSelectorsFile)
	if err != nil {
		log.Printf("Error reading extra selectors file %s: %v", extraSelectorsFile, err)
		return
	}

	var data extraSelectorsData
	err = yaml.Unmarshal(fileContent, &data)
	if err != nil {
		log.Printf("Error unmarshaling extra selectors YAML: %v", err)
		return
	}

	log.Printf("Successfully loaded extra selectors from %s", extraSelectorsFile)
	return data
}

func getExtraSelectors() extraSelectorsData {
	if !extraSelectorsLoaded {
		extraSelectors = loadAndParseExtraSelectors()
		extraSelectorsLoaded = true
	}
	return extraSelectors
}
