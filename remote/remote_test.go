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

	// Test with non-existent chain ID (not in local or remote)
	_, err = GetChainDetailsByChainIDAndFamily(ctx, "7777777777", chain_selectors.FamilyEVM,
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
		// May find from local or remote, both are valid
		assert.NotEmpty(t, details.ChainSelector)
		assert.Contains(t, details.ChainName, "testnet")
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

	// Test caching with mock server (using a remote-only chain)
	t.Run("Caching", func(t *testing.T) {
		callCount := 0
		// Create YAML with a chain that doesn't exist in local
		remoteCacheYAML := `
evm:
  7777777777:
    selector: 8888888888888888888
    name: test-cache-chain
    network_type: testnet
`
		cachingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/x-yaml")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(remoteCacheYAML))
		}))
		t.Cleanup(cachingServer.Close)

		// Clear cache before test
		ClearCache()

		// First call - should hit server (not in local, so must fetch remote)
		_, err := GetChainDetailsBySelector(ctx, 8888888888888888888,
			WithURL(cachingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(1*time.Minute),
		)
		require.NoError(t, err, "First call failed")
		assert.Equal(t, 1, callCount, "Expected 1 server call")

		// Second call - should use cache
		_, err = GetChainDetailsBySelector(ctx, 8888888888888888888,
			WithURL(cachingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(1*time.Minute),
		)
		require.NoError(t, err, "Second call failed")
		assert.Equal(t, 1, callCount, "Expected still 1 server call (cached)")

		// Clear cache
		ClearCache()

		// Third call - should hit server again
		_, err = GetChainDetailsBySelector(ctx, 8888888888888888888,
			WithURL(cachingServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(1*time.Minute),
		)
		require.NoError(t, err, "Third call failed")
		assert.Equal(t, 2, callCount, "Expected 2 server calls after cache clear")
	})

	// Test error handling (using a selector that doesn't exist in local)
	t.Run("ServerError", func(t *testing.T) {
		errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(errorServer.Close)

		// Use a selector that doesn't exist in local, so it tries remote and gets error
		_, err := GetChainDetailsBySelector(ctx, 8888888888888888888,
			WithURL(errorServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache to ensure we hit the error server
		)
		assert.Error(t, err, "Expected error for server error")
	})

	// Test invalid YAML (using a selector that doesn't exist in local)
	t.Run("InvalidYAML", func(t *testing.T) {
		invalidServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid: yaml: content: ["))
		}))
		t.Cleanup(invalidServer.Close)

		// Use a selector that doesn't exist in local, so it tries remote and gets error
		_, err := GetChainDetailsBySelector(ctx, 8888888888888888888,
			WithURL(invalidServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache to ensure we hit the invalid YAML server
		)
		assert.Error(t, err, "Expected error for invalid YAML")
	})
}

func TestLocalFirstWithRemoteFallback(t *testing.T) {
	ClearCache()
	
	// Create a mock server with limited chain data (missing some chains from local package)
	limitedMockYAML := `
evm:
  1:
    selector: 5009297550715157269
    name: ethereum-mainnet
solana:
  "mainnet":
    selector: 124615329519749607
    name: solana-mainnet
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(limitedMockYAML))
	}))
	t.Cleanup(server.Close)

	ctx := context.Background()

	t.Run("GetChainDetailsBySelector_LocalOnly", func(t *testing.T) {
		// Chain that exists in local but not in remote should be found from local
		// Using Ethereum Sepolia testnet selector which exists in local package but not in remote
		sepoliaSelector := uint64(16015286601757825753)
		
		details, err := GetChainDetailsBySelector(ctx, sepoliaSelector,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		require.NoError(t, err, "Expected no error, should find from local")
		assert.Equal(t, chain_selectors.FamilyEVM, details.Family)
		assert.Equal(t, sepoliaSelector, details.ChainSelector)
		assert.NotEmpty(t, details.ChainName)
	})

	t.Run("GetChainDetailsBySelector_FallsBackToRemote", func(t *testing.T) {
		// Create a mock server with a chain that exists only in remote, not in local
		remoteMockYAML := `
evm:
  999999999:
    selector: 8888888888888888888
    name: test-remote-only-chain
    network_type: testnet
`
		remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(remoteMockYAML))
		}))
		t.Cleanup(remoteServer.Close)
		
		// Try to get a chain that doesn't exist in local but exists in remote
		remoteSelector := uint64(8888888888888888888)
		
		details, err := GetChainDetailsBySelector(ctx, remoteSelector,
			WithURL(remoteServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		require.NoError(t, err, "Expected no error - should fall back to remote")
		assert.Equal(t, chain_selectors.FamilyEVM, details.Family)
		assert.Equal(t, "test-remote-only-chain", details.ChainName)
	})

	t.Run("GetChainDetailsByChainIDAndFamily_LocalOnly", func(t *testing.T) {
		// Chain that exists in local but not in remote should be found from local
		// Try to get BSC chain which exists in local package but not in limited remote data
		details, err := GetChainDetailsByChainIDAndFamily(ctx, "56", chain_selectors.FamilyEVM,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		require.NoError(t, err, "Expected no error, should find from local")
		assert.NotEmpty(t, details.ChainSelector)
		assert.NotEmpty(t, details.ChainName)
	})

	t.Run("GetChainDetailsByChainIDAndFamily_FallsBackToRemote", func(t *testing.T) {
		// Create a mock server with a chain that exists only in remote
		remoteMockYAML := `
evm:
  7777777777:
    selector: 8888888888888888888
    name: test-remote-only-chain
    network_type: testnet
`
		remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(remoteMockYAML))
		}))
		t.Cleanup(remoteServer.Close)
		
		// Try to get a chain that doesn't exist in local but exists in remote
		details, err := GetChainDetailsByChainIDAndFamily(ctx, "7777777777", chain_selectors.FamilyEVM,
			WithURL(remoteServer.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		require.NoError(t, err, "Expected no error - should fall back to remote")
		assert.Equal(t, uint64(8888888888888888888), details.ChainSelector)
		assert.Equal(t, "test-remote-only-chain", details.ChainName)
	})

	t.Run("GetChainDetailsByChainIDAndFamily_Solana_LocalOnly", func(t *testing.T) {
		// Try to get Solana devnet which exists in local package but not in limited remote data
		// Using the actual base58 encoded genesis hash
		details, err := GetChainDetailsByChainIDAndFamily(ctx, "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG", chain_selectors.FamilySolana,
			WithURL(server.URL),
			WithTimeout(5*time.Second),
			WithCacheTTL(0), // Disable cache
		)
		require.NoError(t, err, "Expected no error, should find from local")
		assert.NotEmpty(t, details.ChainSelector)
		assert.Contains(t, details.ChainName, "devnet")
	})
}
