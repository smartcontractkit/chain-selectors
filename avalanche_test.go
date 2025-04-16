package chain_selectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AvalancheYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      string
		expectErr     bool
	}{
		{
			name:          "avalanche-testnet-pchain",
			chainSelector: 14538442734262914677,
			chainsId:      "P_CHAIN",
			expectErr:     false,
		},
		{
			name:          "not-exist",
			chainSelector: 12345,
			chainsId:      "NOT_EXIST",
			expectErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err1 := AvalancheNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_AvalancheChainSelectors(t *testing.T) {
	for selector, chainId := range avalancheChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as avalanche family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyAvalanche, family)

		id, err := AvalancheChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_AvalancheGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range avalancheSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyAvalanche)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_AvalancheGetChainIDByChainSelector(t *testing.T) {
	for k, v := range avalancheSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
