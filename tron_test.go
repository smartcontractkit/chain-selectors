package chain_selectors

import (
	"fmt"
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
			chainSelector: 1546563616611573946,
			chainsId:      728126428,
			expectErr:     false,
		},
		{
			name:          "tron-testnet-nile",
			chainSelector: 2052925811360307749,
			chainsId:      3448148188,
			expectErr:     false,
		},
		{
			name:          "tron-testnet-shasta",
			chainSelector: 13231703482326770598,
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
		csObj, err := NewChainSelectorsObj(ChainInfo{})
		assert.NoError(t, err)

		family, err := csObj.GetSelectorFamily(selector)
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
		csObj, err := NewChainSelectorsObj(ChainInfo{})
		assert.NoError(t, err)

		details, err := csObj.GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyTron)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_TronGetChainIDByChainSelector(t *testing.T) {
	for k, v := range tronSelectorsMap {
		csObj, err := NewChainSelectorsObj(ChainInfo{})
		assert.NoError(t, err)

		chainID, err := csObj.GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
