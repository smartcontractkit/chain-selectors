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

func Test_TronYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      uint64
		expectErr     bool
	}{
		{
			name:          "tron-mainnet",
			chainSelector: 1546563616611573945,
			chainsId:      728126428,
			expectErr:     false,
		},
		{
			name:          "tron-testnet-nile",
			chainSelector: 2052925811360307740,
			chainsId:      3448148188,
			expectErr:     false,
		},
		{
			name:          "tron-testnet-shasta",
			chainSelector: 13231703482326770597,
			chainsId:      2494104990,
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
			name, err1 := TronNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_TronChainSelectors(t *testing.T) {
	for selector, chainId := range tronChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as tron family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyTron, family)

		id, err := TronChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_TronGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range tronSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyTron)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_TronGetChainIDByChainSelector(t *testing.T) {
	for k, v := range tronSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}

func Test_TronRemoteFallback(t *testing.T) {
	// Test data - chain not in embedded data
	testSelector := uint64(7777777777777777777)
	testChainID := uint64(999999999)
	testChainName := "tron-test-remote"

	// Create a mock HTTP server
	mockYAML := `
tron:
  999999999:
    selector: 7777777777777777777
    name: "tron-test-remote"
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

	t.Run("TronChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := TronChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("TronNameFromChainId falls back to remote", func(t *testing.T) {
		name, err := TronNameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testChainName, name)
	})
}

func Test_TronRemoteDisabled(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use a selector that definitely doesn't exist in embedded data
	nonExistentSelector := uint64(8888888888888888888)
	nonExistentChainID := uint64(888888888)

	t.Run("TronChainIdFromSelector returns error when remote disabled", func(t *testing.T) {
		_, err := TronChainIdFromSelector(nonExistentSelector)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain id not found")
	})

	t.Run("TronNameFromChainId returns error when remote disabled", func(t *testing.T) {
		_, err := TronNameFromChainId(nonExistentChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain name not found")
	})
}
