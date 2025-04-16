package chain_selectors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BeaconYmlAreValid(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainsId      uint64
		expectErr     bool
	}{
		{
			name:          "ethereum-beacon-mainnet",
			chainSelector: 2007918196257561144,
			chainsId:      1,
			expectErr:     false,
		},
		{
			name:          "not-exist",
			chainSelector: 12345,
			chainsId:      12,
			expectErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, err1 := BeaconNameFromChainId(test.chainsId)
			if test.expectErr {
				require.Error(t, err1)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.name, name)
		})
	}
}

func Test_BeaconChainSelectors(t *testing.T) {
	for selector, chainId := range beaconChainIdBySelector {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"selector %v should be returned as beacon family, but received %v",
			selector, err)
		require.NotEmpty(t, family)
		require.Equal(t, FamilyBeacon, family)

		id, err := BeaconChainIdFromSelector(selector)
		require.Nil(t, err)
		require.Equal(t, chainId, id)
	}
}

func Test_BeaconGetChainDetailsByChainIDAndFamily(t *testing.T) {
	for k, v := range beaconSelectorsMap {
		details, err := GetChainDetailsByChainIDAndFamily(fmt.Sprint(k), FamilyBeacon)
		assert.NoError(t, err)
		assert.Equal(t, v, details)
	}
}

func Test_BeaconGetChainIDByChainSelector(t *testing.T) {
	for k, v := range beaconSelectorsMap {
		chainID, err := GetChainIDFromSelector(v.ChainSelector)
		assert.NoError(t, err)
		assert.Equal(t, chainID, fmt.Sprintf("%v", k))
	}
}
