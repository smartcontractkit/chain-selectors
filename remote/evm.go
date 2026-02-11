package remote

import (
	"context"
	"fmt"
	"strconv"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
)

// EvmChainIdToChainSelector fetches chain data from GitHub and returns a map of EVM chain ID to chain selector
func EvmChainIdToChainSelector(ctx context.Context, opts ...Option) (map[uint64]uint64, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return nil, err
	}

	result := make(map[uint64]uint64, len(cache.evmChainIdToChainSelector))
	for k, v := range cache.evmChainIdToChainSelector {
		result[k] = v.ChainSelector
	}
	return result, nil
}

// EvmChainIdFromName returns the EVM chain ID for a given chain name.
// It first checks local embedded data, then falls back to remote if not found.
func EvmChainIdFromName(ctx context.Context, name string, opts ...Option) (uint64, error) {
	config := applyOptions(opts)
	
	// Try local data first
	if chainId, err := chain_selectors.ChainIdFromName(name); err == nil {
		return chainId, nil
	}
	// If not found locally, try remote
	
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return 0, err
	}

	for chainId, details := range cache.evmChainIdToChainSelector {
		if details.ChainName == name {
			return chainId, nil
		}
	}

	// Before returning error, check if name is actually a chain ID (for chains without a name)
	chainId, err := strconv.ParseUint(name, 10, 64)
	if err == nil {
		if details, exist := cache.evmChainIdToChainSelector[chainId]; exist && details.ChainName == "" {
			return chainId, nil
		}
	}

	return 0, fmt.Errorf("chain not found for name %s", name)
}

// EvmChainBySelector returns the EVM chain for a given selector.
// It first checks local embedded data, then falls back to remote if not found.
func EvmChainBySelector(ctx context.Context, sel uint64, opts ...Option) (chain_selectors.Chain, bool, error) {
	config := applyOptions(opts)
	
	// Try local data first
	if ch, exists := chain_selectors.ChainBySelector(sel); exists {
		return ch, true, nil
	}
	// If not found locally, try remote
	
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return chain_selectors.Chain{}, false, err
	}

	ch, exists := cache.evmChainsBySelector[sel]
	return ch, exists, nil
}

// EvmChainByEvmChainID returns the EVM chain for a given EVM chain ID.
// It first checks local embedded data, then falls back to remote if not found.
func EvmChainByEvmChainID(ctx context.Context, evmChainID uint64, opts ...Option) (chain_selectors.Chain, bool, error) {
	config := applyOptions(opts)
	
	// Try local data first
	if ch, exists := chain_selectors.ChainByEvmChainID(evmChainID); exists {
		return ch, true, nil
	}
	// If not found locally, try remote
	
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return chain_selectors.Chain{}, false, err
	}

	ch, exists := cache.evmChainsByEvmChainID[evmChainID]
	return ch, exists, nil
}

// IsEvm checks if a chain selector is for an EVM chain.
// It first checks local embedded data, then falls back to remote if not found.
func IsEvm(ctx context.Context, chainSel uint64, opts ...Option) (bool, error) {
	config := applyOptions(opts)
	
	// Try local data first
	isEvm, err := chain_selectors.IsEvm(chainSel)
	if err == nil {
		return isEvm, nil
	}
	// If not found locally, try remote
	
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return false, err
	}

	_, exists := cache.evmChainsBySelector[chainSel]
	if !exists {
		return false, fmt.Errorf("chain %d not found", chainSel)
	}
	return true, nil
}
