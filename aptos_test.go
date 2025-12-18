package chain_selectors

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AptosYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      uint64
		expectErr     bool
	}{
		{
			name:          "aptos-mainnet",
			chainSelector: 124615329519749607,
			chainsId:      1,
			expectErr:     false,
		},
		{
			name:          "aptos-testnet",
			chainSelector: 6302590918974934319,
			chainsId:      2,
			expectErr:     false,
		},
		{
			name:          "non-existing",
			chainSelector: 81923186267,
			chainsId:      471,
			expectErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err1 := AptosNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_AptosChainSelectors(t *testing.T) {
	for selector, chain := range aptosChainsBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as aptos family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyAptos, family)

		id, err := AptosChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chain.ChainID, id)

		returnedChain, exists := AptosChainBySelector(selector)
		require.True(t, exists)
		require.Equal(t, returnedChain.ChainID, id)
		require.Equal(t, id, returnedChain.ChainID)
	}
}

func Test_AptosGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range aptosSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyAptos)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_AptosGetChainIDByChainSelector(t *testing.T) {
	for k, v := range aptosSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}

func Test_AptosRemoteFallback(t *testing.T) {
	// Test data - chain not in embedded data
	testSelector := uint64(4444444444444444444)
	testChainID := uint64(999888777)
	testChainName := "aptos-test-remote"

	// Create a mock HTTP server
	mockYAML := `
aptos:
  999888777:
    selector: 4444444444444444444
    name: "aptos-test-remote"
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

	t.Run("AptosChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := AptosChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("AptosNameFromChainId falls back to remote", func(t *testing.T) {
		name, err := AptosNameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testChainName, name)
	})

	t.Run("AptosChainBySelector falls back to remote", func(t *testing.T) {
		chain, exists := AptosChainBySelector(testSelector)
		require.True(t, exists)
		assert.Equal(t, testSelector, chain.Selector)
		assert.Equal(t, testChainID, chain.ChainID)
		assert.Equal(t, testChainName, chain.Name)
	})
}

func Test_AptosRemoteDisabled(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use a selector that definitely doesn't exist in embedded data
	nonExistentSelector := uint64(3333333333333333333)
	nonExistentChainID := uint64(888777666)

	t.Run("AptosChainIdFromSelector returns error when remote disabled", func(t *testing.T) {
		_, err := AptosChainIdFromSelector(nonExistentSelector)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain id not found")
	})

	t.Run("AptosNameFromChainId returns error when remote disabled", func(t *testing.T) {
		_, err := AptosNameFromChainId(nonExistentChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain name not found")
	})

	t.Run("AptosChainBySelector returns false when remote disabled", func(t *testing.T) {
		_, exists := AptosChainBySelector(nonExistentSelector)
		assert.False(t, exists)
	})
}
