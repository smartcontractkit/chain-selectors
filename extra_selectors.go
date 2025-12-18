package chain_selectors

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// remoteDatasourceURL is the URL for fetching all chains from the main branch
// Can be overridden in tests to point to a mock server
var remoteDatasourceURL = "https://raw.githubusercontent.com/smartcontractkit/chain-selectors/main/all_selectors.yml"

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

	// Lazy loading of remote selectors
	remoteSelectorsOnce    sync.Once
	remoteSelectors        extraSelectorsData
	remoteSelectorsFetched bool
)

// isRemoteDatasourceEnabled checks if ENABLE_REMOTE_DATASOURCE is set to true
func isRemoteDatasourceEnabled() bool {
	enabled, _ := strconv.ParseBool(os.Getenv("ENABLE_REMOTE_DATASOURCE"))
	return enabled
}

// tryLazyFetchRemoteSelectors fetches remote selectors on first call (thread-safe)
// Returns true if remote selectors are available
func tryLazyFetchRemoteSelectors() bool {
	if !isRemoteDatasourceEnabled() {
		return false
	}

	remoteSelectorsOnce.Do(func() {
		log.Printf("Lazy loading: Chain not found in embedded, fetching from remote: %s", remoteDatasourceURL)
		remoteSelectors = loadRemoteDatasource()
		remoteSelectorsFetched = true
	})

	return remoteSelectorsFetched
}

// getRemoteChainByID looks up a chain from the remote datasource by family and chain ID
// chainID should be the string representation of the chain ID
func getRemoteChainByID(family string, chainID string) (ChainDetails, bool) {
	if !tryLazyFetchRemoteSelectors() {
		return ChainDetails{}, false
	}

	switch family {
	case FamilyEVM:
		id, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, false
		}
		details, ok := remoteSelectors.Evm[id]
		return details, ok
	case FamilySolana:
		details, ok := remoteSelectors.Solana[chainID]
		return details, ok
	case FamilyAptos:
		id, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, false
		}
		details, ok := remoteSelectors.Aptos[id]
		return details, ok
	case FamilySui:
		id, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, false
		}
		details, ok := remoteSelectors.Sui[id]
		return details, ok
	case FamilyTron:
		id, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, false
		}
		details, ok := remoteSelectors.Tron[id]
		return details, ok
	case FamilyTon:
		id, err := strconv.ParseInt(chainID, 10, 32)
		if err != nil {
			return ChainDetails{}, false
		}
		details, ok := remoteSelectors.Ton[int32(id)]
		return details, ok
	case FamilyStarknet:
		details, ok := remoteSelectors.Starknet[chainID]
		return details, ok
	default:
		return ChainDetails{}, false
	}
}

// getRemoteChainBySelector looks up a chain from the remote datasource by selector
// Returns the family, chain ID (as string), chain details, and whether it was found
//
// Note: This function performs a linear search across all families and chains.
// We could optimize this with a reverse map (selector -> chain info) for O(1) lookups,
// but given that this is a fallback mechanism only called when embedded data doesn't
// have the chain, and the remote datasource is lazy-loaded only once, the additional
// memory and complexity may not be worth it. If this becomes a bottleneck, consider
// building a selector index map after loading the remote datasource.
func getRemoteChainBySelector(selector uint64) (family string, chainID string, details ChainDetails, ok bool) {
	if !tryLazyFetchRemoteSelectors() {
		return "", "", ChainDetails{}, false
	}

	// Check EVM
	for id, d := range remoteSelectors.Evm {
		if d.ChainSelector == selector {
			return FamilyEVM, fmt.Sprintf("%d", id), d, true
		}
	}

	// Check Solana
	for id, d := range remoteSelectors.Solana {
		if d.ChainSelector == selector {
			return FamilySolana, id, d, true
		}
	}

	// Check Aptos
	for id, d := range remoteSelectors.Aptos {
		if d.ChainSelector == selector {
			return FamilyAptos, fmt.Sprintf("%d", id), d, true
		}
	}

	// Check Sui
	for id, d := range remoteSelectors.Sui {
		if d.ChainSelector == selector {
			return FamilySui, fmt.Sprintf("%d", id), d, true
		}
	}

	// Check Tron
	for id, d := range remoteSelectors.Tron {
		if d.ChainSelector == selector {
			return FamilyTron, fmt.Sprintf("%d", id), d, true
		}
	}

	// Check Ton
	for id, d := range remoteSelectors.Ton {
		if d.ChainSelector == selector {
			return FamilyTon, fmt.Sprintf("%d", id), d, true
		}
	}

	// Check Starknet
	for id, d := range remoteSelectors.Starknet {
		if d.ChainSelector == selector {
			return FamilyStarknet, id, d, true
		}
	}

	return "", "", ChainDetails{}, false
}

func loadAndParseExtraSelectors() (result extraSelectorsData) {
	extraSelectorsFile := os.Getenv("EXTRA_SELECTORS_FILE")
	if extraSelectorsFile == "" {
		return
	}

	fileContent, err := os.ReadFile(extraSelectorsFile)
	if err != nil {
		log.Printf("Error reading extra selectors file %s: %v", extraSelectorsFile, err)
		panic(err)
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
		// Only load EXTRA_SELECTORS_FILE at startup
		// Remote datasource is loaded lazily on first miss
		extraSelectors = loadAndParseExtraSelectors()
		extraSelectorsLoaded = true
	}
	return extraSelectors
}

// loadRemoteDatasource fetches chain data from the hardcoded remote URL
// Called lazily when a chain is not found in embedded chains
// Gracefully handles fetch/parse errors (returns empty data to allow embedded chains to work)
// but panics on validation errors (data corruption that should be fixed in the remote datasource)
func loadRemoteDatasource() extraSelectorsData {
	log.Printf("Fetching remote chain data from: %s", remoteDatasourceURL)

	fileContent, err := fetchFromURL(remoteDatasourceURL)
	if err != nil {
		log.Printf("Warning: Failed to fetch remote datasource from %s: %v. Continuing with embedded chains only.", remoteDatasourceURL, err)
		return extraSelectorsData{}
	}

	var data extraSelectorsData
	err = yaml.Unmarshal(fileContent, &data)
	if err != nil {
		log.Printf("Warning: Failed to parse remote datasource YAML: %v. Continuing with embedded chains only.", err)
		return extraSelectorsData{}
	}

	// Validate chain formats - panic on validation errors (data corruption)
	if err := validateSolanaChainID(data.Solana); err != nil {
		log.Printf("Error: Invalid Solana chain IDs in remote datasource: %v", err)
		panic(err)
	}

	if err := validateSuiChainID(data.Sui); err != nil {
		log.Printf("Error: Invalid Sui chain IDs in remote datasource: %v", err)
		panic(err)
	}

	if err := validateAptosChainID(data.Aptos); err != nil {
		log.Printf("Error: Invalid Aptos chain IDs in remote datasource: %v", err)
		panic(err)
	}

	log.Printf("Successfully loaded remote datasource from %s", remoteDatasourceURL)
	return data
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
