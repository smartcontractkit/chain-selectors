package chain_selectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFamilySelectorsMatchGetSelectorFamily(t *testing.T) {
	for _, family := range Families() {
		selectors, err := SelectorsByFamily(family)
		require.NoError(t, err)
		require.NotEmpty(t, selectors, "family %s has no selectors", family)
		for _, selector := range selectors {
			got, err := GetSelectorFamily(selector)
			require.NoError(t, err)
			assert.Equal(t, family, got, "selector %d", selector)
		}
	}
}

func TestAllSelectors(t *testing.T) {
	all := AllSelectors()

	total := 0
	for _, family := range Families() {
		selectors, err := SelectorsByFamily(family)
		require.NoError(t, err)
		total += len(selectors)
	}
	assert.Len(t, all, total, "selectors must be unique across families")
	assert.IsIncreasing(t, all)

	assert.Contains(t, all, ETHEREUM_MAINNET.Selector)
	assert.Contains(t, all, SOLANA_MAINNET.Selector)
	assert.Contains(t, all, APTOS_MAINNET.Selector)
	for _, ch := range ALL {
		assert.Contains(t, all, ch.Selector)
	}
}

func TestSelectorsByFamily_UnknownFamily(t *testing.T) {
	_, err := SelectorsByFamily(FamilyCosmos)
	require.Error(t, err)
}
