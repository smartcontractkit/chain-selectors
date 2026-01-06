package remote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultGitHubRawURL is the default URL for the all_selectors.yml file on GitHub main branch
	DefaultGitHubRawURL = "https://raw.githubusercontent.com/smartcontractkit/chain-selectors/main/all_selectors.yml"
	// DefaultRemoteFetchTimeout is the default timeout for remote HTTP requests
	DefaultRemoteFetchTimeout = 10 * time.Second
	// DefaultCacheTTL is the default time-to-live for cached remote data
	DefaultCacheTTL = 5 * time.Minute
)

var (
	// remoteCache stores the parsed remote data to avoid repeated HTTP calls
	remoteCache     *remoteCacheData
	remoteCacheLock sync.RWMutex
)

type remoteCacheData struct {
	// EVM
	evmChainIdToChainSelector map[uint64]chain_selectors.ChainDetails
	evmChainsBySelector       map[uint64]chain_selectors.Chain
	evmChainsByEvmChainID     map[uint64]chain_selectors.Chain
	// Solana
	solanaChainIdToChainSelector map[string]chain_selectors.ChainDetails
	solanaChainsBySelector       map[uint64]chain_selectors.SolanaChain
	// Aptos
	aptosSelectorsMap     map[uint64]chain_selectors.ChainDetails
	aptosChainsBySelector map[uint64]chain_selectors.AptosChain
	// Sui
	suiSelectorsMap     map[uint64]chain_selectors.ChainDetails
	suiChainsBySelector map[uint64]chain_selectors.SuiChain
	// Ton
	tonSelectorsMap      map[int32]chain_selectors.ChainDetails
	tonChainIdBySelector map[uint64]int32
	// Tron
	tronSelectorsMap      map[uint64]chain_selectors.ChainDetails
	tronChainIdBySelector map[uint64]uint64
	// Starknet
	starknetSelectorsMap     map[string]chain_selectors.ChainDetails
	starknetChainsBySelector map[uint64]chain_selectors.StarknetChain
	// Metadata
	fetchedAt time.Time
}

// Config holds configuration for remote API calls
type Config struct {
	// URL is the raw GitHub URL to fetch the all_selectors.yml file
	// If empty, DefaultGitHubRawURL will be used
	URL string
	// Timeout for the HTTP request
	// If zero, DefaultRemoteFetchTimeout will be used
	Timeout time.Duration
	// CacheTTL is the time-to-live for cached remote data
	// If zero, no caching will be used (always fetch fresh data)
	CacheTTL time.Duration
}

// Option is a functional option for configuring remote API calls
type Option func(*Config)

// WithURL sets a custom URL for fetching the all_selectors.yml file.
// If not provided, DefaultGitHubRawURL will be used.
func WithURL(url string) Option {
	return func(c *Config) {
		c.URL = url
	}
}

// WithTimeout sets the HTTP request timeout.
// If not provided, DefaultRemoteFetchTimeout will be used.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithCacheTTL sets the cache time-to-live.
// Default is 5 minutes. Set to 0 to disable caching (always fetch fresh data).
func WithCacheTTL(ttl time.Duration) Option {
	return func(c *Config) {
		c.CacheTTL = ttl
	}
}

// fetchRemoteSelectors fetches and parses the all_selectors.yml file from GitHub
func fetchRemoteSelectors(ctx context.Context, config *Config) (*remoteCacheData, error) {
	// Check cache first if TTL is set
	if config.CacheTTL > 0 {
		remoteCacheLock.RLock()
		if remoteCache != nil && time.Since(remoteCache.fetchedAt) < config.CacheTTL {
			cached := remoteCache
			remoteCacheLock.RUnlock()
			return cached, nil
		}
		remoteCacheLock.RUnlock()
	}

	url := config.URL
	timeout := config.Timeout

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: timeout,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote selectors from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch remote selectors, status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse YAML
	var data chain_selectors.ExtraSelectorsData
	if err := yaml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse remote selectors YAML: %w", err)
	}

	// Build cache data structure
	cache := &remoteCacheData{
		evmChainIdToChainSelector:    data.Evm,
		evmChainsBySelector:          make(map[uint64]chain_selectors.Chain),
		evmChainsByEvmChainID:        make(map[uint64]chain_selectors.Chain),
		solanaChainIdToChainSelector: data.Solana,
		solanaChainsBySelector:       make(map[uint64]chain_selectors.SolanaChain),
		aptosSelectorsMap:            data.Aptos,
		aptosChainsBySelector:        make(map[uint64]chain_selectors.AptosChain),
		suiSelectorsMap:              data.Sui,
		suiChainsBySelector:          make(map[uint64]chain_selectors.SuiChain),
		tonSelectorsMap:              data.Ton,
		tonChainIdBySelector:         make(map[uint64]int32),
		tronSelectorsMap:             data.Tron,
		tronChainIdBySelector:        make(map[uint64]uint64),
		starknetSelectorsMap:         data.Starknet,
		starknetChainsBySelector:     make(map[uint64]chain_selectors.StarknetChain),
		fetchedAt:                    time.Now(),
	}

	// Build EVM lookup maps
	for chainID, details := range data.Evm {
		chain := chain_selectors.Chain{
			EvmChainID: chainID,
			Selector:   details.ChainSelector,
			Name:       details.ChainName,
		}
		cache.evmChainsBySelector[details.ChainSelector] = chain
		cache.evmChainsByEvmChainID[chainID] = chain
	}

	// Build Solana lookup maps
	for chainID, details := range data.Solana {
		chain := chain_selectors.SolanaChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}
		cache.solanaChainsBySelector[details.ChainSelector] = chain
	}

	// Build Aptos lookup maps
	for chainID, details := range data.Aptos {
		chain := chain_selectors.AptosChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}
		cache.aptosChainsBySelector[details.ChainSelector] = chain
	}

	// Build Sui lookup maps
	for chainID, details := range data.Sui {
		chain := chain_selectors.SuiChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}
		cache.suiChainsBySelector[details.ChainSelector] = chain
	}

	// Build Ton lookup maps
	for chainID, details := range data.Ton {
		cache.tonChainIdBySelector[details.ChainSelector] = chainID
	}

	// Build Tron lookup maps
	for chainID, details := range data.Tron {
		cache.tronChainIdBySelector[details.ChainSelector] = chainID
	}

	// Build Starknet lookup maps
	for chainID, details := range data.Starknet {
		chain := chain_selectors.StarknetChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}
		cache.starknetChainsBySelector[details.ChainSelector] = chain
	}

	// Update cache if TTL is set
	if config.CacheTTL > 0 {
		remoteCacheLock.Lock()
		remoteCache = cache
		remoteCacheLock.Unlock()
	}

	return cache, nil
}

// applyOptions applies functional options and returns a config with defaults
func applyOptions(opts []Option) *Config {
	config := &Config{
		URL:      DefaultGitHubRawURL,
		Timeout:  DefaultRemoteFetchTimeout,
		CacheTTL: DefaultCacheTTL,
	}
	for _, opt := range opts {
		opt(config)
	}
	// Apply defaults if not set
	if config.URL == "" {
		config.URL = DefaultGitHubRawURL
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultRemoteFetchTimeout
	}
	return config
}

// ChainDetailsWithMetadata extends ChainDetails with additional metadata
type ChainDetailsWithMetadata struct {
	chain_selectors.ChainDetails
	Family  string
	ChainID string
}

// GetChainDetailsBySelector fetches chain data from GitHub and returns chain details for a given selector
func GetChainDetailsBySelector(ctx context.Context, selector uint64, opts ...Option) (ChainDetailsWithMetadata, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return ChainDetailsWithMetadata{}, err
	}

	// Check EVM chains
	for chainID, details := range cache.evmChainIdToChainSelector {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilyEVM,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Solana chains
	for chainID, details := range cache.solanaChainIdToChainSelector {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilySolana,
				ChainID:      chainID,
			}, nil
		}
	}

	// Check Aptos chains
	for chainID, details := range cache.aptosSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilyAptos,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Sui chains
	for chainID, details := range cache.suiSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilySui,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Tron chains
	for chainID, details := range cache.tronSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilyTron,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Ton chains
	for chainID, details := range cache.tonSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilyTon,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Starknet chains
	for chainID, details := range cache.starknetSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       chain_selectors.FamilyStarknet,
				ChainID:      chainID,
			}, nil
		}
	}

	return ChainDetailsWithMetadata{}, fmt.Errorf("unknown chain selector %d", selector)
}

// GetChainDetailsByChainIDAndFamily fetches chain data from GitHub and returns chain details for a given chain ID and family
func GetChainDetailsByChainIDAndFamily(ctx context.Context, chainID string, family string, opts ...Option) (chain_selectors.ChainDetails, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return chain_selectors.ChainDetails{}, err
	}

	switch family {
	case chain_selectors.FamilyEVM:
		evmChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.evmChainIdToChainSelector[evmChainID]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case chain_selectors.FamilySolana:
		details, exist := cache.solanaChainIdToChainSelector[chainID]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case chain_selectors.FamilyAptos:
		aptosChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.aptosSelectorsMap[aptosChainID]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case chain_selectors.FamilySui:
		suiChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.suiSelectorsMap[suiChainID]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case chain_selectors.FamilyTron:
		tronChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.tronSelectorsMap[tronChainID]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case chain_selectors.FamilyTon:
		tonChainID, err := strconv.ParseInt(chainID, 10, 32)
		if err != nil {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.tonSelectorsMap[int32(tonChainID)]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case chain_selectors.FamilyStarknet:
		details, exist := cache.starknetSelectorsMap[chainID]
		if !exist {
			return chain_selectors.ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	default:
		return chain_selectors.ChainDetails{}, fmt.Errorf("family %s is not supported", family)
	}
}

// ClearCache clears the remote data cache, forcing the next remote call to fetch fresh data
func ClearCache() {
	remoteCacheLock.Lock()
	remoteCache = nil
	remoteCacheLock.Unlock()
}
