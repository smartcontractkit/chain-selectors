package chain_selectors

import "testing"

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
