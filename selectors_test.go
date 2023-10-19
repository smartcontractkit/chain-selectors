package chain_selectors

import (
	"testing"
)

func TestNoSameChainSelectorsAreGenerated(t *testing.T) {
	chainSelectors := map[uint64]bool{}

	for k, v := range evmChainIdToChainSelector {
		selector := v.ChainSelector
		if _, exist := chainSelectors[selector]; exist {
			t.Errorf("Chain Selectors should be unique. Selector %d is duplicated for chain %d", selector, k)
		}
		chainSelectors[selector] = true
	}
}

func TestNoOverlapBetweenRealAndTestChains(t *testing.T) {
	for k, _ := range selectorsMap {
		if _, exist := testSelectorsMap[k]; exist {
			t.Errorf("Chain %d is duplicated between real and test chains", k)
		}
	}
}

func TestBothSelectorsYmlAndTestSelectorsYmlAreValid(t *testing.T) {
	optimismGoerliSelector, err := SelectorFromChainId(420)
	if err != nil {
		t.Error("Chain Selectors not found for chain 420")
	}
	if optimismGoerliSelector != 2664363617261496610 {
		t.Error("Invalid Chain Selector for chain 420")
	}

	testChainSelector, err := SelectorFromChainId(90000020)
	if err != nil {
		t.Error("Chain Selectors not found for test chain 90000020")
	}
	if testChainSelector != 17810359353458878177 {
		t.Error("Invalid Chain Selector for chain 90000020")
	}
}

func TestEvmChainIdToChainSelectorReturningCopiedMap(t *testing.T) {
	selectors := EvmChainIdToChainSelector()
	selectors[1] = 2

	_, err := ChainIdFromSelector(2)
	if err == nil {
		t.Error("Changes to map should not affect the original map")
	}

	_, err = ChainIdFromSelector(1)
	if err == nil {
		t.Error("Changes to map should not affect the original map")
	}
}

func TestChainIdFromSelector(t *testing.T) {
	_, err := ChainIdFromSelector(0)
	if err == nil {
		t.Error("Should return error if chain selector not found")
	}

	_, err = ChainIdFromSelector(99999999)
	if err == nil {
		t.Error("Should return error if chain selector not found")
	}

	chainId, err := ChainIdFromSelector(13264668187771770619)
	if err != nil {
		t.Error("Should return chain id if chain selector found")
	}
	if chainId != 97 {
		t.Error("Should return correct chain id")
	}
}

func TestSelectorFromChainId(t *testing.T) {
	_, err := SelectorFromChainId(0)
	if err == nil {
		t.Error("Should return error if chain not found")
	}

	_, err = SelectorFromChainId(99999999)
	if err == nil {
		t.Error("Should return error if chain not found")
	}

	chainSelectorId, err := SelectorFromChainId(97)
	if err != nil {
		t.Error("Should return chain selector id if chain found")
	}
	if chainSelectorId != 13264668187771770619 {
		t.Error("Should return correct chain selector id")
	}
}

func TestTestChainIds(t *testing.T) {
	chainIds := TestChainIds()
	if len(chainIds) != len(testSelectorsMap) {
		t.Error("Should return correct number of test chain ids")
	}

	for _, chainId := range chainIds {
		if _, exist := testSelectorsMap[chainId]; !exist {
			t.Error("Should return correct test chain ids")
		}
	}
}

func TestNameFromChainId(t *testing.T) {
	_, err := NameFromChainId(0)
	if err == nil {
		t.Error("Should return error if chain not found")
	}

	_, err = NameFromChainId(99999999)
	if err == nil {
		t.Error("Should return error if chain not found")
	}

	chainName, err := NameFromChainId(97)
	if err != nil {
		t.Error("Should return chain name if chain found")
	}
	if chainName != "bsc-testnet" {
		t.Error("Should return correct chain name")
	}

	chainName, err = NameFromChainId(1337)
	if err != nil {
		t.Error("Should return chain name if chain found")
	}
	if chainName != "1337" {
		t.Error("Should return chain id if name is not defined")
	}
}
