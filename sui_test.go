package chain_selectors

import (
	"fmt"
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
			chainSelector: 124615329519749607,
			chainsId:      1,
			expectErr:     false,
		},
		{
			name:          "sui-testnet",
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
		require.Equal(t, FamilyAptos, family)

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
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyAptos)
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
