package chain_selectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_DogecoinYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      string
		expectErr     bool
	}{
		{
			name:          "dogecoin-mainnet",
			chainSelector: 16302150171372387475,
			chainsId:      "dogecoin_mainnet",
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
			name, err1 := DogecoinNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_DogecoinChainSelectors(t *testing.T) {
	for selector, chainId := range dogecoinChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as dogecoin family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyDogecoin, family)

		id, err := DogecoinChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_DogecoinGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range dogecoinSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyDogecoin)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_DogecoinGetChainIDByChainSelector(t *testing.T) {
	for k, v := range dogecoinSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
