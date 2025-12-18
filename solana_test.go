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

func Test_YmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      string
		expectErr     bool
	}{
		{
			name:          "solana-mainnet",
			chainSelector: 124615329519749607,
			chainsId:      "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d",
			expectErr:     false,
		},
		{
			name:          "solana-testnet",
			chainSelector: 6302590918974934319,
			chainsId:      "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY",
			expectErr:     false,
		},
		{
			name:          "solana-devnet",
			chainSelector: 16423721717087811551,
			chainsId:      "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG",
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
			name, err1 := SolanaNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_SolanaChainSelectors(t *testing.T) {
	for selector, chain := range solanaChainsBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as solana family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilySolana, family)

		id, err := SolanaChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chain.ChainID, id)

		returnedChain, exists := SolanaChainBySelector(selector)
		require.True(t, exists)
		require.Equal(t, returnedChain.ChainID, id)

		require.Equal(t, id, returnedChain.ChainID)
	}
}

func Test_SolanaGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range solanaSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(k, FamilySolana)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_SolanaGetChainIDByChainSelector(t *testing.T) {
	for k, v := range solanaSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}

func Test_SolanaNoOverlapBetweenRealAndTestChains(t *testing.T) {
	for k, _ := range solanaSelectorsMap {
		_, exist := solanaTestSelectorsMap[k]
		assert.False(t, exist, "Chain %d is duplicated between real and test chains", k)
	}
}

func Test_SolanaRemoteFallback(t *testing.T) {
	// Test data - chain not in embedded data
	testSelector := uint64(5555555555555555555)
	// Use a base58 encoded string that represents exactly 32 bytes
	// This is a slight modification of an existing genesis hash
	testChainID := "2eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"
	testChainName := "solana-test-remote"

	// Create a mock HTTP server
	mockYAML := `
solana:
  2eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d:
    selector: 5555555555555555555
    name: "solana-test-remote"
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

	t.Run("SolanaChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := SolanaChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("SolanaNameFromChainId falls back to remote", func(t *testing.T) {
		name, err := SolanaNameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testChainName, name)
	})

	t.Run("SolanaChainBySelector falls back to remote", func(t *testing.T) {
		chain, exists := SolanaChainBySelector(testSelector)
		require.True(t, exists)
		assert.Equal(t, testSelector, chain.Selector)
		assert.Equal(t, testChainID, chain.ChainID)
		assert.Equal(t, testChainName, chain.Name)
	})
}

func Test_SolanaRemoteDisabled(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use a selector that definitely doesn't exist in embedded data
	nonExistentSelector := uint64(6666666666666666666)
	// Use a slight modification of an existing genesis hash
	nonExistentChainID := "7eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"

	t.Run("SolanaChainIdFromSelector returns error when remote disabled", func(t *testing.T) {
		_, err := SolanaChainIdFromSelector(nonExistentSelector)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain not found")
	})

	t.Run("SolanaNameFromChainId returns error when remote disabled", func(t *testing.T) {
		_, err := SolanaNameFromChainId(nonExistentChainID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain name not found")
	})

	t.Run("SolanaChainBySelector returns false when remote disabled", func(t *testing.T) {
		_, exists := SolanaChainBySelector(nonExistentSelector)
		assert.False(t, exists)
	})
}
