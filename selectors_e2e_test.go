package chain_selectors_test

import (
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
)

func TestAccessGeneratedChain(t *testing.T) {
	assert.Equal(t, uint64(43114), chain_selectors.AVALANCHE_MAINNET.EvmChainID)
}
