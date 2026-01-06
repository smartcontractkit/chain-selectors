package chain_selectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CantonChainSelectors(t *testing.T) {
	for selector, chain := range cantonChainsBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as canton family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyCanton, family)

		id, err := CantonChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chain.ChainID, id)

		returnedChain, exists := CantonChainBySelector(selector)
		require.True(t, exists)
		require.Equal(t, returnedChain.ChainID, id)
		require.Equal(t, id, returnedChain.ChainID)
	}
}

func Test_CantonGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range cantonChainsByChainId {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyCanton)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_CantonGetChainIDByChainSelector(t *testing.T) {
	for k, v := range cantonChainsByChainId {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
