package chain_selectors

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type extraSelectorsData struct {
	Evm      map[uint64]ChainDetails `yaml:"evm"`
	Aptos    map[uint64]ChainDetails `yaml:"aptos"`
	Solana   map[string]ChainDetails `yaml:"solana"`
	Sui      map[uint64]ChainDetails `yaml:"sui"`
	Ton      map[int32]ChainDetails  `yaml:"ton"`
	Tron     map[uint64]ChainDetails `yaml:"tron"`
	Starknet map[string]ChainDetails `yaml:"starknet"`
}

var (
	extraSelectors       extraSelectorsData
	extraSelectorsLoaded bool
)

func loadAndParseExtraSelectors() (result extraSelectorsData) {
	extraSelectorsFile := os.Getenv("EXTRA_SELECTORS_FILE")
	if extraSelectorsFile == "" {
		return
	}

	var fileContent []byte
	var err error

	// Check if it's a URL
	if strings.HasPrefix(extraSelectorsFile, "http://") || strings.HasPrefix(extraSelectorsFile, "https://") {
		log.Printf("Fetching extra selectors from URL: %s", extraSelectorsFile)
		fileContent, err = fetchFromURL(extraSelectorsFile)
		if err != nil {
			log.Printf("Error fetching extra selectors from URL %s: %v", extraSelectorsFile, err)
			panic(err)
		}
		log.Printf("Successfully fetched extra selectors from URL")
	} else {
		// File path (existing behavior)
		fileContent, err = os.ReadFile(extraSelectorsFile)
		if err != nil {
			log.Printf("Error reading extra selectors file %s: %v", extraSelectorsFile, err)
			panic(err)
		}
	}

	var data extraSelectorsData
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

// fetchFromURL fetches content from an HTTP/HTTPS URL with a timeout
func fetchFromURL(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}
