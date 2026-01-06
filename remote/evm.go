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

// EvmChainIdFromName fetches chain data from GitHub and returns the EVM chain ID for a given chain name
func EvmChainIdFromName(ctx context.Context, name string, opts ...Option) (uint64, error) {
	config := applyOptions(opts)
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

// EvmChainBySelector fetches chain data from GitHub and returns the EVM chain for a given selector
func EvmChainBySelector(ctx context.Context, sel uint64, opts ...Option) (chain_selectors.Chain, bool, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return chain_selectors.Chain{}, false, err
	}

	ch, exists := cache.evmChainsBySelector[sel]
	return ch, exists, nil
}

// EvmChainByEvmChainID fetches chain data from GitHub and returns the EVM chain for a given EVM chain ID
func EvmChainByEvmChainID(ctx context.Context, evmChainID uint64, opts ...Option) (chain_selectors.Chain, bool, error) {
	config := applyOptions(opts)
	cache, err := fetchRemoteSelectors(ctx, config)
	if err != nil {
		return chain_selectors.Chain{}, false, err
	}

	ch, exists := cache.evmChainsByEvmChainID[evmChainID]
	return ch, exists, nil
}

// IsEvm fetches chain data from GitHub and checks if a chain selector is for an EVM chain
func IsEvm(ctx context.Context, chainSel uint64, opts ...Option) (bool, error) {
	config := applyOptions(opts)
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
