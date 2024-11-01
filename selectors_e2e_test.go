package chain_selectors_test

import (
	"strconv"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
)

func TestAccessGeneratedChain(t *testing.T) {
	assert.Equal(t, strconv.FormatUint(uint64(43114), 10), chain_selectors.AVALANCHE_MAINNET.ChainID)
}
