package chain_selectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoSameChainSelectorsAreGenerated(t *testing.T) {
	chainSelectors := map[uint64]struct{}{}

	for k, v := range evmChainIdToChainSelector {
		selector := v.ChainSelector
		_, exist := chainSelectors[selector]
		assert.False(t, exist, "Chain Selectors should be unique. Selector %d is duplicated for chain %d", selector, k)
		chainSelectors[selector] = struct{}{}
	}
}

func TestNoOverlapBetweenRealAndTestChains(t *testing.T) {
	for k, _ := range selectorsMap {
		_, exist := testSelectorsMap[k]
		assert.False(t, exist, "Chain %d is duplicated between real and test chains", k)
	}
}

func TestBothSelectorsYmlAndTestSelectorsYmlAreValid(t *testing.T) {
	optimismGoerliSelector, err := SelectorFromChainId(420)
	require.NoError(t, err)
	assert.Equal(t, uint64(2664363617261496610), optimismGoerliSelector)

	testChainSelector, err := SelectorFromChainId(90000020)
	require.NoError(t, err)
	assert.Equal(t, uint64(17810359353458878177), testChainSelector)
}

func TestEvmChainIdToChainSelectorReturningCopiedMap(t *testing.T) {
	selectors := EvmChainIdToChainSelector()
	selectors[1] = 2

	_, err := ChainIdFromSelector(2)
	assert.Error(t, err)

	_, err = ChainIdFromSelector(1)
	assert.Error(t, err)
}

func TestChainIdFromSelector(t *testing.T) {
	_, err := ChainIdFromSelector(0)
	assert.Error(t, err, "Should return error if chain selector not found")

	_, err = ChainIdFromSelector(99999999)
	assert.Error(t, err, "Should return error if chain selector not found")

	chainId, err := ChainIdFromSelector(13264668187771770619)
	require.NoError(t, err)
	assert.Equal(t, uint64(97), chainId)
}

func TestSelectorFromChainId(t *testing.T) {
	_, err := SelectorFromChainId(0)
	require.Error(t, err)

	_, err = SelectorFromChainId(99999999)
	require.Error(t, err)

	chainSelectorId, err := SelectorFromChainId(97)
	require.NoError(t, err)
	assert.Equal(t, uint64(13264668187771770619), chainSelectorId)
}

func TestTestChainIds(t *testing.T) {
	chainIds := TestChainIds()
	assert.Equal(t, len(chainIds), len(testSelectorsMap), "Should return correct number of test chain ids")

	for _, chainId := range chainIds {
		_, exist := testSelectorsMap[chainId]
		assert.True(t, exist)
	}
}

func TestNameFromChainId(t *testing.T) {
	_, err := NameFromChainId(2)
	require.Error(t, err, "Should return error if chain not found")

	_, err = NameFromChainId(99999999)
	require.Error(t, err, "Should return error if chain not found")

	chainName, err := NameFromChainId(97)
	require.NoError(t, err)
	assert.Equal(t, "binance_smart_chain-testnet", chainName)

	chainName, err = NameFromChainId(1337)
	require.NoError(t, err)
	assert.Equal(t, "1337", chainName)
}
