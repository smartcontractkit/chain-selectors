package remote

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvmChainIdToChainSelector(t *testing.T) {
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	mapping, err := EvmChainIdToChainSelector(ctx,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.NotEmpty(t, mapping)

	// Test against a known chain (Ethereum Mainnet)
	ethereumMainnetChainID := uint64(1)
	expectedSelector := uint64(5009297550715157269)

	selector, exists := mapping[ethereumMainnetChainID]
	assert.True(t, exists, "Expected Ethereum Mainnet (chain ID %d) to exist", ethereumMainnetChainID)
	assert.Equal(t, expectedSelector, selector)
}

func TestEvmChainIdFromName(t *testing.T) {
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test with Ethereum Mainnet name
	name := "ethereum-mainnet"
	expectedChainID := uint64(1)

	chainID, err := EvmChainIdFromName(ctx, name,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.Equal(t, expectedChainID, chainID)

	// Test with non-existent name
	_, err = EvmChainIdFromName(ctx, "non-existent-chain",
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	assert.Error(t, err)
}

func TestEvmChainBySelector(t *testing.T) {
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test with Ethereum Mainnet selector
	ethereumMainnetSelector := uint64(5009297550715157269)

	chain, exists, err := EvmChainBySelector(ctx, ethereumMainnetSelector,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.True(t, exists, "Expected Ethereum Mainnet to exist")
	assert.Equal(t, uint64(1), chain.EvmChainID)
	assert.Equal(t, ethereumMainnetSelector, chain.Selector)
	assert.Equal(t, "ethereum-mainnet", chain.Name)

	// Test with non-existent selector
	_, exists, err = EvmChainBySelector(ctx, uint64(999999999999999999),
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.False(t, exists, "Expected chain to not exist for invalid selector")
}

func TestEvmChainByEvmChainID(t *testing.T) {
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test with Ethereum Mainnet chain ID
	ethereumMainnetChainID := uint64(1)

	chain, exists, err := EvmChainByEvmChainID(ctx, ethereumMainnetChainID,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.True(t, exists, "Expected Ethereum Mainnet to exist")
	assert.Equal(t, ethereumMainnetChainID, chain.EvmChainID)
	assert.Equal(t, "ethereum-mainnet", chain.Name)

	// Test with non-existent chain ID (using a very large number that won't exist in local or remote)
	_, exists, err = EvmChainByEvmChainID(ctx, uint64(7777777777),
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.False(t, exists, "Expected chain to not exist for invalid chain ID")
}

func TestIsEvm(t *testing.T) {
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test with valid EVM selector
	ethereumMainnetSelector := uint64(5009297550715157269)

	isEvm, err := IsEvm(ctx, ethereumMainnetSelector,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.True(t, isEvm, "Expected selector to be EVM")

	// Test with non-existent selector
	_, err = IsEvm(ctx, uint64(999999999999999999),
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	assert.Error(t, err)
}

func TestEvmMultipleChainsConsistency(t *testing.T) {
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test a few well-known chains for consistency
	// Note: Since we check local first, the names may come from local data
	testCases := []struct {
		chainID  uint64
		selector uint64
		name     string
	}{
		{1, 5009297550715157269, "ethereum-mainnet"},
		{10, 3734403246176062136, "optimism-mainnet"},
		{56, 11344663589394136015, "bsc-mainnet"},
		{137, 4051577828743386545, "polygon-mainnet"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test chain by selector - check IDs match
			chain, exists, err := EvmChainBySelector(ctx, tc.selector,
				WithURL(server.URL),
				WithTimeout(5*time.Second),
			)
			require.NoError(t, err)
			assert.True(t, exists, "Expected chain with selector %d to exist", tc.selector)
			assert.Equal(t, tc.chainID, chain.EvmChainID)
			assert.NotEmpty(t, chain.Name, "Chain name should not be empty")

			// Test chain by EVM chain ID - check selector matches
			chain, exists, err = EvmChainByEvmChainID(ctx, tc.chainID,
				WithURL(server.URL),
				WithTimeout(5*time.Second),
			)
			require.NoError(t, err)
			assert.True(t, exists, "Expected chain with chain ID %d to exist", tc.chainID)
			assert.Equal(t, tc.selector, chain.Selector)
			assert.NotEmpty(t, chain.Name, "Chain name should not be empty")
		})
	}
}

