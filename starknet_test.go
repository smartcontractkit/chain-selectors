package chain_selectors

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StarknetYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      string
		expectErr     bool
	}{
		{
			name:          "ethereum-mainnet-starknet-1",
			chainSelector: 511843109281680063,
			chainsId:      "SN_MAIN",
			expectErr:     false,
		},
		{
			name:          "ethereum-testnet-sepolia-starknet-1",
			chainSelector: 4115550741429562104,
			chainsId:      "SN_SEPOLIA",
			expectErr:     false,
		},
		{
			name:          "non-existing",
			chainSelector: rand.Uint64(),
			chainsId:      "non-existing",
			expectErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err1 := StarknetNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_StarknetChainSelectors(t *testing.T) {
	for selector, chain := range starknetChainsBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as starknet family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyStarknet, family)

		id, err := StarknetChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chain.ChainID, id)

		returnedChain, exists := StarknetChainBySelector(selector)
		require.True(t, exists)
		require.Equal(t, returnedChain.ChainID, id)

		require.Equal(t, id, returnedChain.ChainID)
	}
}

func Test_StarknetGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range starknetSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(k, FamilyStarknet)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_StarknetGetChainIDByChainSelector(t *testing.T) {
	for k, v := range starknetSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}

func Test_StarknetRemoteFallback(t *testing.T) {
	// Test data - chain not in embedded data
	testSelector := uint64(1111111111111111111)
	testChainID := "SN_TEST_REMOTE"
	testChainName := "starknet-test-remote"

	// Create a mock HTTP server
	mockYAML := `
starknet:
  SN_TEST_REMOTE:
    selector: 1111111111111111111
    name: "starknet-test-remote"
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockYAML))
	}))
	t.Cleanup(server.Close)

	// Override the remote datasource URL
	setRemoteDatasourceURL(t, server.URL)

	// Reset remote state for this test
	originalRemoteSelectors := remoteSelectors
	originalRemoteSelectorsFetched := remoteSelectorsFetched
	originalOnce := remoteSelectorsOnce
	t.Cleanup(func() {
		remoteSelectors = originalRemoteSelectors
		remoteSelectorsFetched = originalRemoteSelectorsFetched
		remoteSelectorsOnce = originalOnce
	})

	remoteSelectorsOnce = sync.Once{}
	remoteSelectorsFetched = false

	// Enable remote datasource
	t.Setenv("ENABLE_REMOTE_DATASOURCE", "true")

	t.Run("StarknetChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := StarknetChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("StarknetNameFromChainId falls back to remote", func(t *testing.T) {
		name, err := StarknetNameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testChainName, name)
	})

	t.Run("StarknetChainBySelector falls back to remote", func(t *testing.T) {
		chain, exists := StarknetChainBySelector(testSelector)
		require.True(t, exists)
		assert.Equal(t, testSelector, chain.Selector)
		assert.Equal(t, testChainID, chain.ChainID)
		assert.Equal(t, testChainName, chain.Name)
	})
}

func Test_StarknetRemoteDisabled(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use a selector that definitely doesn't exist in embedded data
	nonExistentSelector := uint64(2222222222222222222)
	nonExistentChainID := "SN_NONEXISTENT"

	t.Run("StarknetChainIdFromSelector returns error when remote disabled", func(t *testing.T) {
		_, err := StarknetChainIdFromSelector(nonExistentSelector)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain not found")
	})

	t.Run("StarknetNameFromChainId returns error when remote disabled", func(t *testing.T) {
		_, err := StarknetNameFromChainId(nonExistentChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain name not found")
	})

	t.Run("StarknetChainBySelector returns false when remote disabled", func(t *testing.T) {
		_, exists := StarknetChainBySelector(nonExistentSelector)
		assert.False(t, exists)
	})
}
