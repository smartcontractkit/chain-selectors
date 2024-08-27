package chain_selectors

import (
	"fmt"
	"math/rand"
	"regexp"
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

func TestAllChainSelectorsHaveFamilies(t *testing.T) {
	for _, ch := range ALL {
		family, err := GetSelectorFamily(ch.Selector)
		require.NoError(t, err,
			"Family not found for selector %d (chain id %d, name %s), please update selector_families.yml with the appropriate chain family for this chain",
			ch.Selector, ch.EvmChainID, ch.Name)
		require.NotEmpty(t, family)
	}
}

func Test_ChainSelectors(t *testing.T) {
	tests := []struct {
		name          string
		chainSelector uint64
		chainId       uint64
		expectErr     bool
	}{
		{
			name:          "bsc chain",
			chainSelector: 13264668187771770619,
			chainId:       97,
		},
		{
			name:          "optimism chain",
			chainSelector: 2664363617261496610,
			chainId:       420,
		},
		{
			name:          "test chain",
			chainSelector: 17810359353458878177,
			chainId:       90000020,
		},
		{
			name:          "not existing chain",
			chainSelector: 120398123,
			chainId:       123454312,
			expectErr:     true,
		},
		{
			name:          "invalid selector and chain id",
			chainSelector: 0,
			chainId:       0,
			expectErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chainId, err1 := ChainIdFromSelector(test.chainSelector)
			chainSelector, err2 := SelectorFromChainId(test.chainId)
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
		_, exist := testSelectorsMap[chainId]
		assert.True(t, exist)
	}
}

func Test_ChainNames(t *testing.T) {
	tests := []struct {
		name      string
		chainName string
		chainId   uint64
		expectErr bool
	}{
		{
			name:      "zkevm chain with a dedicated name",
			chainName: "ethereum-testnet-goerli-polygon-zkevm-1",
			chainId:   1442,
		},
		{
			name:      "bsc chain with a dedicated name",
			chainName: "binance_smart_chain-testnet",
			chainId:   97,
		},
		{
			name:      "chain without a name defined",
			chainName: "geth-testnet",
			chainId:   1337,
		},
		{
			name:      "test simulated chain without a dedicated name",
			chainName: "90000013",
			chainId:   90000013,
		},
		{
			name:      "not existing chain",
			chainName: "avalanche-testnet-mumbai-1",
			chainId:   120398123,
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
			chainId, err1 := ChainIdFromName(test.chainName)
			chainName, err2 := NameFromChainId(test.chainId)
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
	t.Run("exist", func(t *testing.T) {
		for _, ch := range ALL {
			v, exists := ChainBySelector(ch.Selector)
			assert.True(t, exists)
			assert.Equal(t, ch, v)
		}
	})

	t.Run("non existent", func(t *testing.T) {
		_, exists := ChainBySelector(rand.Uint64())
		assert.False(t, exists)
	})
}

func Test_ChainByEvmChainID(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		for _, ch := range ALL {
			v, exists := ChainByEvmChainID(ch.EvmChainID)
			assert.True(t, exists)
			assert.Equal(t, ch, v)
		}
	})

	t.Run("non existent", func(t *testing.T) {
		_, exists := ChainByEvmChainID(rand.Uint64())
		assert.False(t, exists)
	})
}

func Test_IsEvm(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		for _, ch := range ALL {
			isEvm, err := IsEvm(ch.Selector)
			assert.NoError(t, err)
			assert.True(t, isEvm)
		}
	})

	t.Run("non existent", func(t *testing.T) {
		isEvm, err := IsEvm(rand.Uint64())
		assert.Error(t, err)
		assert.False(t, isEvm)
	})
}

func Test_MainnetAndTestnetAreSynchronized(t *testing.T) {
	re := regexp.MustCompile(`^([a-zA-Z0-9-]+)-(mainnet|testnet)`)
	parseName := func(str string) (chain string, isMainnet bool, err error) {
		matches := re.FindStringSubmatch(str)

		if matches == nil {
			return "", false, fmt.Errorf("no matches found")
		}

		return matches[1], matches[2] == "mainnet", nil
	}

	type chainDetails struct {
		mainnet []string
		testnet []string
	}

	var chainMap = make(map[string]chainDetails)
	for _, chain := range ALL {
		name, isMainnet, err := parseName(chain.Name)
		if err == nil {
			details := chainMap[name]
			if isMainnet {
				details.mainnet = append(details.mainnet, chain.Name)
			} else {
				details.testnet = append(details.testnet, chain.Name)
			}
			chainMap[name] = details
		}
	}

	// Skip legacy testnet only chains
	delete(chainMap, "bitcoin")
	delete(chainMap, "geth")
	delete(chainMap, "hedera")
	delete(chainMap, "berachain")
	delete(chainMap, "cronos")
	delete(chainMap, "fantom")

	// analyze results
	for chain, details := range chainMap {
		if len(details.mainnet) == 0 && len(details.testnet) != 0 {
			assert.Fail(t, "Chain %s has testnet chains but no mainnet chains. If this is intentional add an exception to this test with a comment explaining why.", chain)
		}
	}
}
