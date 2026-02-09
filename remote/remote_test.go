package remote

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockYAML contains comprehensive test data for all blockchain families
// Format matches ExtraSelectorsData struct
const mockYAML = `
evm:
  1:
    selector: 5009297550715157269
    name: ethereum-mainnet
  10:
    selector: 3734403246176062136
    name: optimism-mainnet
  56:
    selector: 11344663589394136015
    name: bsc-mainnet
  137:
    selector: 4051577828743386545
    name: polygon-mainnet
  42161:
    selector: 4949039107694359620
    name: arbitrum-mainnet
solana:
  "mainnet":
    selector: 124615329519749607
    name: solana-mainnet
  "testnet":
    selector: 1666700230607807939
    name: solana-testnet
  "devnet":
    selector: 7633325390517157182
    name: solana-devnet
aptos:
  1:
    selector: 5880489174233984516
    name: aptos-mainnet
  2:
    selector: 6433500567565415381
    name: aptos-testnet
sui:
  1:
    selector: 5790810961207155433
    name: sui-mainnet
  2:
    selector: 4419140975832851138
    name: sui-testnet
ton:
  -239:
    selector: 5264266016034146460
    name: ton-mainnet
  -3:
    selector: 3989674961303603008
    name: ton-testnet
tron:
  728126428:
    selector: 10902061574536243337
    name: tron-mainnet
  2494104990:
    selector: 8482416459910711315
    name: tron-testnet-shasta
starknet:
  "SN_MAIN":
    selector: 3919063707296401440
    name: ethereum-mainnet-starknet-1
  "SN_SEPOLIA":
    selector: 1924942427828825923
    name: ethereum-testnet-sepolia-starknet-1
`

// newMockServer creates a test HTTP server that serves the mockYAML
func newMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockYAML))
	}))
}

func TestGetChainDetailsBySelector(t *testing.T) {
	ClearCache()
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test with Ethereum Mainnet selector
	ethereumMainnetSelector := uint64(5009297550715157269)

	details, err := GetChainDetailsBySelector(ctx, ethereumMainnetSelector,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.Equal(t, chain_selectors.FamilyEVM, details.Family)
	assert.Equal(t, "1", details.ChainID)
	assert.Equal(t, "ethereum-mainnet", details.ChainName)
	assert.Equal(t, ethereumMainnetSelector, details.ChainSelector)

	// Test with non-existent selector
	_, err = GetChainDetailsBySelector(ctx, uint64(999999999999999999),
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	assert.Error(t, err)
}

func TestGetChainDetailsByChainIDAndFamily(t *testing.T) {
	ClearCache()
	server := newMockServer()
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test with Ethereum Mainnet
	details, err := GetChainDetailsByChainIDAndFamily(ctx, "1", chain_selectors.FamilyEVM,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.Equal(t, uint64(5009297550715157269), details.ChainSelector)
	assert.Equal(t, "ethereum-mainnet", details.ChainName)

	// Test with non-existent chain ID
	_, err = GetChainDetailsByChainIDAndFamily(ctx, "999999999", chain_selectors.FamilyEVM,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	assert.Error(t, err)

	// Test with Solana chain
	details, err = GetChainDetailsByChainIDAndFamily(ctx, "mainnet", chain_selectors.FamilySolana,
		WithURL(server.URL),
		WithTimeout(5*time.Second),
	)
	require.NoError(t, err)
	assert.Equal(t, uint64(124615329519749607), details.ChainSelector)
	assert.Equal(t, "solana-mainnet", details.ChainName)
}

func TestRemoteWithMockServer(t *testing.T) {
	ClearCache()
	// Create a mock server that returns test data
	mockYAML := `evm:
  1:
    selector: 5009297550715157269
    name: ethereum-mainnet
  137:
    selector: 4051577828743386545
    name: polygon-mainnet
solana:
  "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d":
    selector: 124615329519749607
    name: solana-mainnet
aptos:
  1:
    selector: 4741433654867091981
    name: aptos-mainnet
canton:
  MainNet:
    selector: 2199546568103630433
    name: canton-mainnet
  TestNet:
    selector: 13503176106905080262
    name: canton-testnet
`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockYAML))
	}))
	t.Cleanup(server.Close)

	ctx := context.Background()

	// Test GetChainDetailsBySelector with mock server
	t.Run("GetChainDetailsBySelector", func(t *testing.T) {
		details, err := GetChainDetailsBySelector(ctx, 5009297550715157269,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
		)
		require.NoError(t, err)
		assert.Equal(t, chain_selectors.FamilyEVM, details.Family)
		assert.Equal(t, "1", details.ChainID)
		assert.Equal(t, "ethereum-mainnet", details.ChainName)
	})

	// Test GetChainDetailsByChainIDAndFamily with mock server
	t.Run("GetChainDetailsByChainIDAndFamily", func(t *testing.T) {
		details, err := GetChainDetailsByChainIDAndFamily(ctx, "137", chain_selectors.FamilyEVM,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
		)
		require.NoError(t, err)
		assert.Equal(t, uint64(4051577828743386545), details.ChainSelector)
		assert.Equal(t, "polygon-mainnet", details.ChainName)
	})

	// Test with Solana chain
	t.Run("GetSolanaChainDetails", func(t *testing.T) {
		details, err := GetChainDetailsBySelector(ctx, 124615329519749607,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
		)
		require.NoError(t, err)
		assert.Equal(t, chain_selectors.FamilySolana, details.Family)
		// Solana mainnet returns the actual on-chain ID, not just "mainnet"
		assert.NotEmpty(t, details.ChainID, "Chain ID should not be empty")
		assert.Equal(t, "solana-mainnet", details.ChainName)
	})

	// Test with Canton chain
	t.Run("GetCantonChainDetails", func(t *testing.T) {
		details, err := GetChainDetailsBySelector(ctx, 2199546568103630433,
			WithURL(server.URL),
		)
		require.NoError(t, err)
		assert.Equal(t, chain_selectors.FamilyCanton, details.Family)
		assert.Equal(t, "MainNet", details.ChainID)
		assert.Equal(t, "canton-mainnet", details.ChainName)
	})
	t.Run("GetCantonChainDetailsByChainIDAndFamily", func(t *testing.T) {
		details, err := GetChainDetailsByChainIDAndFamily(ctx, "TestNet", chain_selectors.FamilyCanton,
			WithURL(server.URL),
		)
		require.NoError(t, err)
		assert.Equal(t, uint64(13503176106905080262), details.ChainSelector)
		assert.Equal(t, "canton-testnet", details.ChainName)
	})

	// Test EVM-specific functions
	t.Run("EvmChainByEvmChainID", func(t *testing.T) {
		chain, exists, err := EvmChainByEvmChainID(ctx, 1,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
		)
		require.NoError(t, err)
		assert.True(t, exists, "Expected chain to exist")
		assert.Equal(t, uint64(1), chain.EvmChainID)
		assert.Equal(t, uint64(5009297550715157269), chain.Selector)
		assert.Equal(t, "ethereum-mainnet", chain.Name)
	})

	// Test caching with mock server
	t.Run("Caching", func(t *testing.T) {
		callCount := 0
		cachingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/x-yaml")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockYAML))
		}))
		t.Cleanup(cachingServer.Close)

		// Clear cache before test
		ClearCache()

		// First call - should hit server
		_, err := GetChainDetailsBySelector(ctx, 5009297550715157269,
			WithURL(cachingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(1*time.Minute),
		)
		require.NoError(t, err, "First call failed")
		assert.Equal(t, 1, callCount, "Expected 1 server call")

		// Second call - should use cache
		_, err = GetChainDetailsBySelector(ctx, 5009297550715157269,
			WithURL(cachingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(1*time.Minute),
		)
		require.NoError(t, err, "Second call failed")
		assert.Equal(t, 1, callCount, "Expected still 1 server call (cached)")

		// Clear cache
		ClearCache()

		// Third call - should hit server again
		_, err = GetChainDetailsBySelector(ctx, 5009297550715157269,
			WithURL(cachingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(1*time.Minute),
		)
		require.NoError(t, err, "Third call failed")
		assert.Equal(t, 2, callCount, "Expected 2 server calls after cache clear")
	})

	// Test error handling
	t.Run("ServerError", func(t *testing.T) {
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(errorServer.Close)

		_, err := GetChainDetailsBySelector(ctx, 5009297550715157269,
			WithURL(errorServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache to ensure we hit the error server
		)
		assert.Error(t, err, "Expected error for server error")
	})

	// Test invalid YAML
	t.Run("InvalidYAML", func(t *testing.T) {
		invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid: yaml: content: ["))
		}))
		t.Cleanup(invalidServer.Close)

		_, err := GetChainDetailsBySelector(ctx, 5009297550715157269,
			WithURL(invalidServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache to ensure we hit the invalid YAML server
		)
		assert.Error(t, err, "Expected error for invalid YAML")
	})
}

func TestFallbackToLocal(t *testing.T) {
	ClearCache()
	ctx := context.Background()

	// Create a server that always fails to simulate network failure
	failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate network failure by returning 500
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(failingServer.Close)

	// Test GetChainDetailsBySelector with fallback disabled (default)
	t.Run("NoFallback_ShouldReturnError", func(t *testing.T) {
		_, err := GetChainDetailsBySelector(ctx, chain_selectors.ETHEREUM_MAINNET.Selector,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		assert.Error(t, err, "Expected error when remote fetch fails without fallback")
	})

	// Test GetChainDetailsBySelector with fallback enabled
	t.Run("WithFallback_ShouldReturnLocalData", func(t *testing.T) {
		details, err := GetChainDetailsBySelector(ctx, chain_selectors.ETHEREUM_MAINNET.Selector,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
			WithFallbackToLocal(true),
		)
		require.NoError(t, err, "Expected no error when fallback is enabled")
		assert.Equal(t, chain_selectors.FamilyEVM, details.Family)
		assert.Equal(t, "1", details.ChainID)
		assert.Equal(t, chain_selectors.ETHEREUM_MAINNET.Selector, details.ChainSelector)
	})

	// Test GetChainDetailsByChainIDAndFamily with fallback disabled
	t.Run("NoFallback_ByChainID_ShouldReturnError", func(t *testing.T) {
		_, err := GetChainDetailsByChainIDAndFamily(ctx, "1", chain_selectors.FamilyEVM,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		assert.Error(t, err, "Expected error when remote fetch fails without fallback")
	})

	// Test GetChainDetailsByChainIDAndFamily with fallback enabled
	t.Run("WithFallback_ByChainID_ShouldReturnLocalData", func(t *testing.T) {
		details, err := GetChainDetailsByChainIDAndFamily(ctx, "1", chain_selectors.FamilyEVM,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
			WithFallbackToLocal(true),
		)
		require.NoError(t, err, "Expected no error when fallback is enabled")
		assert.Equal(t, chain_selectors.ETHEREUM_MAINNET.Selector, details.ChainSelector)
	})

	// Test with timeout to simulate network timeout
	t.Run("NetworkTimeout_WithFallback", func(t *testing.T) {
		// Create a server that never responds
		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second) // Sleep longer than timeout
		}))
		t.Cleanup(slowServer.Close)

		details, err := GetChainDetailsBySelector(ctx, chain_selectors.ETHEREUM_MAINNET.Selector,
			WithURL(slowServer.URL),
			WithTimeout(100*time.Millisecond), // Very short timeout
			WithCacheTTL(0),
			WithFallbackToLocal(true),
		)
		require.NoError(t, err, "Expected no error with fallback on timeout")
		assert.Equal(t, chain_selectors.FamilyEVM, details.Family)
		assert.Equal(t, chain_selectors.ETHEREUM_MAINNET.Selector, details.ChainSelector)
	})

	// Test fallback with non-EVM chains
	t.Run("Fallback_Solana", func(t *testing.T) {
		// Get a known Solana selector from local data
		solanaChains := chain_selectors.SolanaChainIdToChainSelector()
		if len(solanaChains) == 0 {
			t.Skip("No Solana chains available in local data")
		}

		// Get the first available Solana chain
		var testSelector uint64
		for _, selector := range solanaChains {
			testSelector = selector
			break
		}

		details, err := GetChainDetailsBySelector(ctx, testSelector,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0),
			WithFallbackToLocal(true),
		)
		require.NoError(t, err, "Expected no error with fallback for Solana")
		assert.Equal(t, chain_selectors.FamilySolana, details.Family)
		assert.Equal(t, testSelector, details.ChainSelector)
	})

	// Test fallback returns original error when chain not found in local
	t.Run("Fallback_ChainNotFoundInLocal_ReturnsOriginalError", func(t *testing.T) {
		// Use a selector that doesn't exist in local data
		unknownSelector := uint64(9999999999999999)

		_, err := GetChainDetailsBySelector(ctx, unknownSelector,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0),
			WithFallbackToLocal(true),
		)
		// Should return the original remote fetch error, not "unknown chain selector" error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch remote selectors")
	})

	// Test fallback returns original error for GetChainDetailsByChainIDAndFamily
	t.Run("Fallback_ByChainID_NotFoundInLocal_ReturnsOriginalError", func(t *testing.T) {
		// Use a chain ID that doesn't exist in local data
		_, err := GetChainDetailsByChainIDAndFamily(ctx, "9999999999", chain_selectors.FamilyEVM,
			WithURL(failingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0),
			WithFallbackToLocal(true),
		)
		// Should return the original remote fetch error, not "invalid chain id" error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch remote selectors")
	})
}
