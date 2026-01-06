package remote

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

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

// Family constants for blockchain types
const (
	FamilyEVM      = "evm"
	FamilySolana   = "solana"
	FamilyAptos    = "aptos"
	FamilySui      = "sui"
	FamilyTon      = "ton"
	FamilyTron     = "tron"
	FamilyStarknet = "starknet"
)

var (
	// remoteCache stores the parsed remote data to avoid repeated HTTP calls
	remoteCache     *remoteCacheData
	remoteCacheLock sync.RWMutex
)

type remoteCacheData struct {
	// EVM
	evmChainIdToChainSelector map[uint64]ChainDetails
	evmChainsBySelector       map[uint64]Chain
	evmChainsByEvmChainID     map[uint64]Chain
	// Solana
	solanaChainIdToChainSelector map[string]ChainDetails
	solanaChainsBySelector       map[uint64]SolanaChain
	// Aptos
	aptosSelectorsMap     map[uint64]ChainDetails
	aptosChainsBySelector map[uint64]AptosChain
	// Sui
	suiSelectorsMap     map[uint64]ChainDetails
	suiChainsBySelector map[uint64]SuiChain
	// Ton
	tonSelectorsMap      map[int32]ChainDetails
	tonChainIdBySelector map[uint64]int32
	// Tron
	tronSelectorsMap      map[uint64]ChainDetails
	tronChainIdBySelector map[uint64]uint64
	// Starknet
	starknetSelectorsMap     map[string]ChainDetails
	starknetChainsBySelector map[uint64]StarknetChain
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

// ChainDetails represents the basic chain information
type ChainDetails struct {
	ChainSelector uint64 `yaml:"selector" json:"chainSelector"`
	ChainName     string `yaml:"name" json:"chainName"`
}

// Chain represents an EVM chain
type Chain struct {
	EvmChainID uint64
	Selector   uint64
	Name       string
}

// SolanaChain represents a Solana chain
type SolanaChain struct {
	ChainID  string
	Selector uint64
	Name     string
}

// AptosChain represents an Aptos chain
type AptosChain struct {
	ChainID  uint64
	Selector uint64
	Name     string
}

// SuiChain represents a Sui chain
type SuiChain struct {
	ChainID  uint64
	Selector uint64
	Name     string
}

// StarknetChain represents a Starknet chain
type StarknetChain struct {
	ChainID  string
	Selector uint64
	Name     string
}

// ExtraSelectorsData represents the structure of the all_selectors.yml file
type ExtraSelectorsData struct {
	Evm      map[uint64]ChainDetails `yaml:"evm"`
	Solana   map[string]ChainDetails `yaml:"solana"`
	Aptos    map[uint64]ChainDetails `yaml:"aptos"`
	Sui      map[uint64]ChainDetails `yaml:"sui"`
	Ton      map[int32]ChainDetails  `yaml:"ton"`
	Tron     map[uint64]ChainDetails `yaml:"tron"`
	Starknet map[string]ChainDetails `yaml:"starknet"`
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
	var data ExtraSelectorsData
	if err := yaml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse remote selectors YAML: %w", err)
	}

	// Build cache data structure
	cache := &remoteCacheData{
		evmChainIdToChainSelector:    data.Evm,
		evmChainsBySelector:          make(map[uint64]Chain),
		evmChainsByEvmChainID:        make(map[uint64]Chain),
		solanaChainIdToChainSelector: data.Solana,
		solanaChainsBySelector:       make(map[uint64]SolanaChain),
		aptosSelectorsMap:            data.Aptos,
		aptosChainsBySelector:        make(map[uint64]AptosChain),
		suiSelectorsMap:              data.Sui,
		suiChainsBySelector:          make(map[uint64]SuiChain),
		tonSelectorsMap:              data.Ton,
		tonChainIdBySelector:         make(map[uint64]int32),
		tronSelectorsMap:             data.Tron,
		tronChainIdBySelector:        make(map[uint64]uint64),
		starknetSelectorsMap:         data.Starknet,
		starknetChainsBySelector:     make(map[uint64]StarknetChain),
		fetchedAt:                    time.Now(),
	}

	// Build EVM lookup maps
	for chainID, details := range data.Evm {
		chain := Chain{
			EvmChainID: chainID,
			Selector:   details.ChainSelector,
			Name:       details.ChainName,
		}
		cache.evmChainsBySelector[details.ChainSelector] = chain
		cache.evmChainsByEvmChainID[chainID] = chain
	}

	// Build Solana lookup maps
	for chainID, details := range data.Solana {
		chain := SolanaChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}
		cache.solanaChainsBySelector[details.ChainSelector] = chain
	}

	// Build Aptos lookup maps
	for chainID, details := range data.Aptos {
		chain := AptosChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		}
		cache.aptosChainsBySelector[details.ChainSelector] = chain
	}

	// Build Sui lookup maps
	for chainID, details := range data.Sui {
		chain := SuiChain{
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
		chain := StarknetChain{
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
	ChainDetails
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
				Family:       FamilyEVM,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Solana chains
	for chainID, details := range cache.solanaChainIdToChainSelector {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       FamilySolana,
				ChainID:      chainID,
			}, nil
		}
	}

	// Check Aptos chains
	for chainID, details := range cache.aptosSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       FamilyAptos,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Sui chains
	for chainID, details := range cache.suiSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       FamilySui,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Tron chains
	for chainID, details := range cache.tronSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       FamilyTron,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Ton chains
	for chainID, details := range cache.tonSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       FamilyTon,
				ChainID:      fmt.Sprintf("%d", chainID),
			}, nil
		}
	}

	// Check Starknet chains
	for chainID, details := range cache.starknetSelectorsMap {
		if details.ChainSelector == selector {
			return ChainDetailsWithMetadata{
				ChainDetails: details,
				Family:       FamilyStarknet,
				ChainID:      chainID,
			}, nil
		}
	}

	return ChainDetailsWithMetadata{}, fmt.Errorf("unknown chain selector %d", selector)
}

// GetChainDetailsByChainIDAndFamily fetches chain data from GitHub and returns chain details for a given chain ID and family
func GetChainDetailsByChainIDAndFamily(ctx context.Context, chainID string, family string, opts ...Option) (ChainDetails, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return ChainDetails{}, err
	}

	switch family {
	case FamilyEVM:
		evmChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.evmChainIdToChainSelector[evmChainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case FamilySolana:
		details, exist := cache.solanaChainIdToChainSelector[chainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case FamilyAptos:
		aptosChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.aptosSelectorsMap[aptosChainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case FamilySui:
		suiChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.suiSelectorsMap[suiChainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case FamilyTron:
		tronChainID, err := strconv.ParseUint(chainID, 10, 64)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.tronSelectorsMap[tronChainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case FamilyTon:
		tonChainID, err := strconv.ParseInt(chainID, 10, 32)
		if err != nil {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		details, exist := cache.tonSelectorsMap[int32(tonChainID)]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	case FamilyStarknet:
		details, exist := cache.starknetSelectorsMap[chainID]
		if !exist {
			return ChainDetails{}, fmt.Errorf("invalid chain id %s for %s", chainID, family)
		}

		return details, nil

	default:
		return ChainDetails{}, fmt.Errorf("family %s is not supported", family)
	}
}

// ClearCache clears the remote data cache, forcing the next remote call to fetch fresh data
func ClearCache() {
	remoteCacheLock.Lock()
	remoteCache = nil
	remoteCacheLock.Unlock()
}

// TonChain represents a TON chain
type TonChain struct {
	ChainID  int32
	Selector uint64
	Name     string
}

// TronChain represents a TRON chain
type TronChain struct {
	ChainID  uint64
	Selector uint64
	Name     string
}

// EvmGetAllChains fetches chain data from GitHub and returns all EVM chains
func EvmGetAllChains(ctx context.Context, opts ...Option) ([]Chain, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	chains := make([]Chain, 0, len(cache.evmChainsBySelector))
	for _, chain := range cache.evmChainsBySelector {
		chains = append(chains, chain)
	}
	return chains, nil
}

// SolanaGetAllChains fetches chain data from GitHub and returns all Solana chains
func SolanaGetAllChains(ctx context.Context, opts ...Option) ([]SolanaChain, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	chains := make([]SolanaChain, 0, len(cache.solanaChainsBySelector))
	for _, chain := range cache.solanaChainsBySelector {
		chains = append(chains, chain)
	}
	return chains, nil
}

// AptosGetAllChains fetches chain data from GitHub and returns all Aptos chains
func AptosGetAllChains(ctx context.Context, opts ...Option) ([]AptosChain, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	chains := make([]AptosChain, 0, len(cache.aptosChainsBySelector))
	for _, chain := range cache.aptosChainsBySelector {
		chains = append(chains, chain)
	}
	return chains, nil
}

// SuiGetAllChains fetches chain data from GitHub and returns all Sui chains
func SuiGetAllChains(ctx context.Context, opts ...Option) ([]SuiChain, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	chains := make([]SuiChain, 0, len(cache.suiChainsBySelector))
	for _, chain := range cache.suiChainsBySelector {
		chains = append(chains, chain)
	}
	return chains, nil
}

// TonGetAllChains fetches chain data from GitHub and returns all TON chains
func TonGetAllChains(ctx context.Context, opts ...Option) ([]TonChain, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	chains := make([]TonChain, 0, len(cache.tonSelectorsMap))
	for chainID, details := range cache.tonSelectorsMap {
		chains = append(chains, TonChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		})
	}
	return chains, nil
}

// TronGetAllChains fetches chain data from GitHub and returns all TRON chains
func TronGetAllChains(ctx context.Context, opts ...Option) ([]TronChain, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	chains := make([]TronChain, 0, len(cache.tronSelectorsMap))
	for chainID, details := range cache.tronSelectorsMap {
		chains = append(chains, TronChain{
			ChainID:  chainID,
			Selector: details.ChainSelector,
			Name:     details.ChainName,
		})
	}
	return chains, nil
}

// TonChainIdToChainSelector fetches chain data from GitHub and returns a map of TON chain ID to chain selector
func TonChainIdToChainSelector(ctx context.Context, opts ...Option) (map[int32]uint64, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	result := make(map[int32]uint64, len(cache.tonSelectorsMap))
	for k, v := range cache.tonSelectorsMap {
		result[k] = v.ChainSelector
	}
	return result, nil
}

// SolanaChainIdToChainSelector fetches chain data from GitHub and returns a map of Solana chain ID to chain selector
func SolanaChainIdToChainSelector(ctx context.Context, opts ...Option) (map[string]uint64, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	result := make(map[string]uint64, len(cache.solanaChainIdToChainSelector))
	for k, v := range cache.solanaChainIdToChainSelector {
		result[k] = v.ChainSelector
	}
	return result, nil
}

