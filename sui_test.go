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

func Test_SuiYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      uint64
		expectErr     bool
	}{
		{
			name:          "sui-mainnet",
			chainSelector: 17529533435026248318,
			chainsId:      1,
			expectErr:     false,
		},
		{
			name:          "sui-testnet",
			chainSelector: 9762610643973837292,
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
			name, err1 := SuiNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_SuiChainSelectors(t *testing.T) {
	for selector, chain := range suiChainsBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as sui family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilySui, family)

		id, err := SuiChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chain.ChainID, id)

		returnedChain, exists := SuiChainBySelector(selector)
		require.True(t, exists)
		require.Equal(t, returnedChain.ChainID, id)
		require.Equal(t, id, returnedChain.ChainID)
	}
}

func Test_SuiGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range suiSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilySui)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_SuiGetChainIDByChainSelector(t *testing.T) {
	for k, v := range suiSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}

func Test_SuiRemoteFallback(t *testing.T) {
	// Test data - chain not in embedded data
	testSelector := uint64(3333333333333333333)
	testChainID := uint64(999666333)
	testChainName := "sui-test-remote"

	// Create a mock HTTP server
	mockYAML := `
sui:
  999666333:
    selector: 3333333333333333333
    name: "sui-test-remote"
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

	t.Run("SuiChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := SuiChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("SuiNameFromChainId falls back to remote", func(t *testing.T) {
		name, err := SuiNameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testChainName, name)
	})

	t.Run("SuiChainBySelector falls back to remote", func(t *testing.T) {
		chain, exists := SuiChainBySelector(testSelector)
		require.True(t, exists)
		assert.Equal(t, testSelector, chain.Selector)
		assert.Equal(t, testChainID, chain.ChainID)
		assert.Equal(t, testChainName, chain.Name)
	})
}

func Test_SuiRemoteDisabled(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use a selector that definitely doesn't exist in embedded data
	nonExistentSelector := uint64(2222222222222222222)
	nonExistentChainID := uint64(777555333)

	t.Run("SuiChainIdFromSelector returns error when remote disabled", func(t *testing.T) {
		_, err := SuiChainIdFromSelector(nonExistentSelector)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain id not found")
	})

	t.Run("SuiNameFromChainId returns error when remote disabled", func(t *testing.T) {
		_, err := SuiNameFromChainId(nonExistentChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain name not found")
	})

	t.Run("SuiChainBySelector returns false when remote disabled", func(t *testing.T) {
		_, exists := SuiChainBySelector(nonExistentSelector)
		assert.False(t, exists)
	})
}
