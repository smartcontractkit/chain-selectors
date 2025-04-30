package chain_selectors

import (
	"fmt"
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
