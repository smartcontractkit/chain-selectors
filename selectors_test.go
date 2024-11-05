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

	for selector := range chainSelectorToDetails {
		_, exist := chainSelectors[selector]
		assert.False(t, exist, "Chain Selectors should be unique. Selector %d is duplicated for chain %d", selector)
		chainSelectors[selector] = struct{}{}
	}
}

func TestNoSameChainIDAndFamilyAreGenerated(t *testing.T) {
	chainIDAndFamily := map[string]struct{}{}

	for _, details := range chainSelectorToDetails {
		key := fmt.Sprintf("%s:%s", details.ChainID, details.Family)
		_, exist := chainIDAndFamily[key]
		assert.False(t, exist, "ChainID within single family should be unique. chainID %s is duplicated for family", details.ChainID, details.Family)
		chainIDAndFamily[key] = struct{}{}
	}
}

func TestNoOverlapBetweenRealAndTestChains(t *testing.T) {
	for k, _ := range selectorsMap {
		_, exist := testSelectorsMap[k]
		assert.False(t, exist, "Chain %d is duplicated between real and test chains", k)
	}
}

func TestBothSelectorsYmlAndTestSelectorsYmlAreValid(t *testing.T) {
	optimismGoerliSelector, err := SelectorFromChainIdAndFamily("420", "")
	require.NoError(t, err)
	assert.Equal(t, uint64(2664363617261496610), optimismGoerliSelector)

	evm, err := IsEvm(optimismGoerliSelector)
	assert.NoError(t, err)
	assert.True(t, evm)

	solanaMainnetSelector, err := SelectorFromChainIdAndFamily("5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d", "solana")
	require.NoError(t, err)
	assert.Equal(t, uint64(124615329519749607), solanaMainnetSelector)

	evm, err = IsEvm(solanaMainnetSelector)
	assert.NoError(t, err)
	assert.False(t, evm)

	testChainSelector, exist := ChainBySelector(17810359353458878177)
	require.True(t, exist)
	assert.Equal(t, "90000020", testChainSelector.ChainID)
}

func TestChainIdToChainSelectorReturningCopiedMap(t *testing.T) {
	selectors := ChainSelectorToChainDetails()
	tmp := selectors[5009297550715157269]
	tmp.ChainID = "2"
	selectors[5009297550715157269] = tmp

	chainID, err := GetChainIdFromSelector(5009297550715157269)
	assert.NoError(t, err)
	assert.NotEqual(t, chainID, tmp)
}

func TestAllChainSelectorsHaveFamilies(t *testing.T) {
	for _, ch := range ALL {
		family, err := GetSelectorFamily(ch.Selector)
		require.NoError(t, err,
			"Family not found for selector %d (chain id %s, name %s), please update test_selectors_restructured.yml with the appropriate chain family for this chain",
			ch.Selector, ch.ChainID, ch.Name)
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
			chainId, err1 := GetChainIdFromSelector(test.chainSelector)
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

	for _, chainId := range chainIds {
		selector, err := SelectorFromChainId(chainId)
		if err != nil {
			return
		}
		_, exist := testSelectorsMap[selector]
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
			chainID, err1 := ChainIdFromNameAndFamily(test.chainName, FamilyEVM)
			selector, err2 := SelectorFromChainIdAndFamily(chainID, "")
			if test.expectErr {
				require.Error(t, err1)
				require.Error(t, err2)
				return
			}
			require.NoError(t, err1)
			assert.Equal(t, test.chainId, chainID)

			require.NoError(t, err2)
			detail, _ := selectorsMap[selector]
			assert.Equal(t, test.chainName, detail.Name)
		})
	}
}

func Test_ChainBySelector(t *testing.T) {
	testMap := loadYML(testSelectorsYml)
	t.Run("exist", func(t *testing.T) {
		for selector, details := range testMap {
			v, exists := ChainBySelector(selector)
			assert.True(t, exists)
			assert.Equal(t, details.ChainID, v.ChainID)
		}
	})

	t.Run("non existent", func(t *testing.T) {
		_, exists := ChainBySelector(rand.Uint64())
		assert.False(t, exists)
	})
}

func Test_SelectorMap(t *testing.T) {
	selectorMap := loadYML(selectorYml)
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

func Test_IsEvm(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		for _, ch := range ALL {
			if ch.Family == FamilyEVM {
				exist, err := IsEvm(ch.Selector)
				assert.NoError(t, err)
				assert.True(t, exist)
			}
		}
	})

	t.Run("non existent", func(t *testing.T) {
		exist, err := IsEvm(rand.Uint64())
		assert.Error(t, err)
		assert.False(t, exist)
	})
}
