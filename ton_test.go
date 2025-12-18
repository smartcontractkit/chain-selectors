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

func Test_TonYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      int32
		expectErr     bool
	}{
		{
			name:          "ton-mainnet",
			chainSelector: 16448340667252469081,
			chainsId:      -239,
			expectErr:     false,
		},
		{
			name:          "ton-testnet",
			chainSelector: 1399300952838017768,
			chainsId:      -3,
			expectErr:     false,
		},
		{
			name:          "ton-localnet",
			chainSelector: 13879075125137744094,
			chainsId:      -217,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err1 := TonNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_TonChainSelectors(t *testing.T) {
	for selector, chainId := range tonChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as ton family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyTon, family)

		id, err := TonChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_TonGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range tonSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyTon)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_TonGetChainIDByChainSelector(t *testing.T) {
	for k, v := range tonSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}

func Test_TonRemoteFallback(t *testing.T) {
	// Test data - chain not in embedded data
	testSelector := uint64(8888888888888888888)
	testChainID := int32(-999)
	testChainName := "ton-test-remote"

	// Create a mock HTTP server
	mockYAML := `
ton:
  -999:
    selector: 8888888888888888888
    name: "ton-test-remote"
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

	t.Run("TonChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := TonChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("TonNameFromChainId falls back to remote", func(t *testing.T) {
		name, err := TonNameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testChainName, name)
	})
}

func Test_TonRemoteDisabled(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use a selector that definitely doesn't exist in embedded data
	nonExistentSelector := uint64(9999999999999999999)
	nonExistentChainID := int32(99999)

	t.Run("TonChainIdFromSelector returns error when remote disabled", func(t *testing.T) {
		_, err := TonChainIdFromSelector(nonExistentSelector)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain id not found")
	})

	t.Run("TonNameFromChainId returns error when remote disabled", func(t *testing.T) {
		_, err := TonNameFromChainId(nonExistentChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain name not found")
	})
}
