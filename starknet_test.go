package chain_selectors

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StarknetYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      string
		expectErr     bool
	}{
		{
			name:          "ethereum-mainnet-starknet-1",
			chainSelector: 511843109281680063,
			chainsId:      "SN_MAIN",
			expectErr:     false,
		},
		{
			name:          "ethereum-testnet-sepolia-starknet-1",
			chainSelector: 4115550741429562104,
			chainsId:      "SN_SEPOLIA",
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
			name, err1 := StarknetNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_StarknetChainSelectors(t *testing.T) {
	for selector, chain := range starknetChainsBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as starknet family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyStarknet, family)

		id, err := StarknetChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chain.ChainID, id)

		returnedChain, exists := StarknetChainBySelector(selector)
		require.True(t, exists)
		require.Equal(t, returnedChain.ChainID, id)

		require.Equal(t, id, returnedChain.ChainID)
	}
}

func Test_StarknetGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range starknetSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(k, FamilyStarknet)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_StarknetGetChainIDByChainSelector(t *testing.T) {
	for k, v := range starknetSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
