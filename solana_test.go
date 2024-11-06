package chain_selectors

import "testing"

func TestSolana(t *testing.T) {
	for k, v := range solanaSelectorsMap {
		t.Logf("k: %s, v: %v", k, v)
	}
}
