package chain_selectors

import (
	"math/rand"
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
	for selector, chainId := range solanaChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as solana family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilySolana, family)

		id, err := SolanaChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_SolanaGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range solanaSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(k, FamilySolana)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}
