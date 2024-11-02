package chain_selectors

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoSameChainSelectorsAreGenerated(t *testing.T) {
	chainSelectors := map[uint64]struct{}{}

	for selector := range selectorToChainDetails {
		_, exist := chainSelectors[selector]
		assert.False(t, exist, "Chain Selectors should be unique. Selector %d is duplicated for chain %d", selector)
		chainSelectors[selector] = struct{}{}
	}
}

func TestNoOverlapBetweenRealAndTestChains(t *testing.T) {
	for k, _ := range selectorToChainDetails {
		_, exist := testSelectorsMap[k]
		assert.False(t, exist, "Chain %d is duplicated between real and test chains", k)
	}
}

func TestBothSelectorsYmlAndTestSelectorsYmlAreValid(t *testing.T) {
	optimismGoerliSelector, err := SelectorFromChainIdAndFamily("420", "")
	require.NoError(t, err)
	assert.Equal(t, uint64(2664363617261496610), optimismGoerliSelector)

	testChainSelector, exist := TestChainBySelector(17810359353458878177)
	require.True(t, exist)
	assert.Equal(t, "90000020", testChainSelector.ChainID)
}

func TestEvmChainIdToChainSelectorReturningCopiedMap(t *testing.T) {
	selectors := ChainSelectorToChainDetails()
	tmp := selectors[5009297550715157269]
	tmp.ChainID = "2"
	selectors[5009297550715157269] = tmp

	chainID, err := ChainIdFromSelector(5009297550715157269)
	assert.NoError(t, err)
	assert.NotEqual(t, chainID, tmp)
}

func TestAllChainSelectorsHaveFamilies(t *testing.T) {
	for selector, details := range selectorToChainDetails {
		family, err := GetSelectorFamily(selector)
		require.NoError(t, err,
			"Family not found for selector %d (chain id %d, name %s), please update selector_families.yml with the appropriate chain family for this chain",
			selector, details.ChainID, details.Name)
		require.NotEmpty(t, family)
	}

	for selector, details := range testSelectorsMap {
		family, err := GetTestSelectorFamily(selector)
		require.NoError(t, err,
			"Family not found for selector %d (chain id %d, name %s), please update selector_families.yml with the appropriate chain family for this chain",
			selector, details.ChainID, details.Name)
		require.NotEmpty(t, family)
	}
}

func Test_ChainSelectors(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainId       string
		expectErr     bool
	}{
		{
			name:          "bsc chain",
			chainSelector: 13264668187771770619,
			chainId:       "97",
		},
		{
			name:          "optimism chain",
			chainSelector: 2664363617261496610,
			chainId:       "420",
		},
		{
			name:          "not existing chain",
			chainSelector: 120398123,
			chainId:       "123454312",
			expectErr:     true,
		},
		{
			name:          "invalid selector and chain id",
			chainSelector: 0,
			chainId:       "0",
			expectErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chainId, err1 := ChainIdFromSelector(test.chainSelector)
			chainSelector, err2 := SelectorFromChainIdAndFamily(test.chainId, "")
			if test.expectErr {
				require.Error(t, err1)
				require.Error(t, err2)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.chainId, chainId)

			require.NoError(t, err2)
			assert.Equal(t, test.chainSelector, chainSelector)
		})
	}
}

func Test_TestChainIds(t *testing.T) {
	chainIds := TestChainIds()
	assert.Equal(t, len(chainIds), len(testSelectorsMap), "Should return correct number of test chain ids")

	tests := []struct {
		name          string
		chainSelector uint64
		chainId       string
		expectErr     bool
	}{
		{
			name:          "test chain",
			chainSelector: 17810359353458878177,
			chainId:       "90000020",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			details, exists := TestChainBySelector(test.chainSelector)
			assert.True(t, exists)
			assert.Equal(t, details.ChainID, test.chainId)
		})
	}

	for _, chainId := range chainIds {
		_, exist := testSelectorsMap[chainId]
		assert.True(t, exist)
	}
}

func Test_ChainNames(t *testing.T) {
	tests := []struct {
		name      string
		chainName string
		chainId   string
		expectErr bool
	}{
		{
			name:      "zkevm chain with a dedicated name",
			chainName: "ethereum-testnet-goerli-polygon-zkevm-1",
			chainId:   "1442",
		},
		{
			name:      "bsc chain with a dedicated name",
			chainName: "binance_smart_chain-testnet",
			chainId:   "97",
		},
		{
			name:      "chain without a name defined",
			chainName: "geth-testnet",
			chainId:   "1337",
		},
		{
			name:      "not existing chain",
			chainName: "avalanche-testnet-mumbai-1",
			chainId:   "120398123",
			expectErr: true,
		},
		{
			name:      "should return error if chain id passed as a name for chain with a full name",
			chainName: "1",
			expectErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chainId, err1 := ChainIdFromNameAndFamily(test.chainName, "")
			chainName, err2 := NameFromChainIdAndFamily(test.chainId, "")
			if test.expectErr {
				require.Error(t, err1)
				require.Error(t, err2)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.chainId, chainId)

			require.NoError(t, err2)
			assert.Equal(t, test.chainName, chainName)
		})
	}
}

func Test_ChainBySelector(t *testing.T) {
	testMap := loadTestChains()
	t.Run("exist", func(t *testing.T) {
		for selector, details := range testMap {
			v, exists := TestChainBySelector(selector)
			assert.True(t, exists)
			assert.Equal(t, details.ChainID, v.ChainID)
		}
	})

	t.Run("non existent", func(t *testing.T) {
		_, exists := TestChainBySelector(rand.Uint64())
		assert.False(t, exists)
	})
}

func Test_SelectorMap(t *testing.T) {
	selectorMap := loadChainDetailsBySelector()
	t.Run("exist", func(t *testing.T) {
		for selector, details := range selectorMap {
			v, err := SelectorFromChainIdAndFamily(details.ChainID, details.Family)
			assert.Nil(t, err)
			if selector != v {
				fmt.Printf("%v %v", selector, v)
			}
		}
	})

	t.Run("non existent", func(t *testing.T) {
		_, err := SelectorFromChainIdAndFamily(strconv.FormatUint(rand.Uint64(), 10), "")
		assert.Error(t, err)
	})
}

func Test_TestSelectorMap(t *testing.T) {
	testMap := loadTestChains()
	t.Run("exist", func(t *testing.T) {
		for selector := range testMap {
			exist, err := TestChainSelectorExist(selector)
			assert.NoError(t, err)
			assert.True(t, exist)
		}
	})

	t.Run("non existent", func(t *testing.T) {
		exist, err := TestChainSelectorExist(rand.Uint64())
		assert.Error(t, err)
		assert.False(t, exist)
	})
}
