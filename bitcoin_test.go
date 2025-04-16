package chain_selectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BitcoinYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      string
		expectErr     bool
	}{
		{
			name:          "bitcoin-mainnet",
			chainSelector: 5295418179606748534,
			chainsId:      "bitcoin_mainnet",
			expectErr:     false,
		},
		{
			name:          "bitcoin-testnet",
			chainSelector: 8997884582513251897,
			chainsId:      "bitcoin_testnet",
			expectErr:     false,
		},
		{
			name:          "not-exist",
			chainSelector: 12345,
			chainsId:      "not_exist",
			expectErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err1 := BitcoinNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_BitcoinChainSelectors(t *testing.T) {
	for selector, chainId := range bitcoinChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as bitcoin family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyBitcoin, family)

		id, err := BitcoinChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_BitcoinGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range bitcoinSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyBitcoin)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_BitcoinGetChainIDByChainSelector(t *testing.T) {
	for k, v := range bitcoinSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
