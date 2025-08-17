package chain_selectors

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempYamlFile(t *testing.T, yamlContent string) string {
	tmpFile, err := os.CreateTemp("", "test_*.yaml")
	require.NoError(t, err)

	_, err = tmpFile.WriteString(yamlContent)
	require.NoError(t, err)
	tmpFile.Close()

	return tmpFile.Name()
}

func setSelectorEnv(_ *testing.T, filePath string) func() {
	os.Setenv("EXTRA_SELECTORS_FILE", filePath)
	return func() {
		os.Unsetenv("EXTRA_SELECTORS_FILE")
	}
}

func runTestWithYaml(t *testing.T, testName string, yamlContent string, validate func(*testing.T, extraSelectorsData)) {
	t.Run(testName, func(t *testing.T) {
		filePath := createTempYamlFile(t, yamlContent)
		defer os.Remove(filePath)

		cleanup := setSelectorEnv(t, filePath)
		defer cleanup()

		result := loadAndParseExtraSelectors()
		validate(t, result)
	})
}

const (
	yamlSingleFamily = `
evm:
  999:
    selector: 1234567890123456789
    name: "test-evm-chain"
`

	yamlMultipleFamilies = `
evm:
  999:
    selector: 1234567890123456789
    name: "test-evm-chain"
solana:
  "test-solana-genesis": 
    selector: 1111111111111111111
    name: "test-solana-chain"
ton:
  -666:
    selector: 3333333333333333333
    name: "test-ton-chain"
aptos:
  888:
    selector: 9876543210987654321
    name: "test-aptos-chain"
sui:
  777:
    selector: 2222222222222222222
    name: "test-sui-chain"
tron:
  123:
    selector: 123557890123456789
    name: "test-tron-chain"
`
)

func TestExtraSelectors(t *testing.T) {
	runTestWithYaml(t, "Extra selectors: single family", yamlSingleFamily, func(t *testing.T, result extraSelectorsData) {
		assert.Len(t, result.Evm, 1)
		assert.Empty(t, result.Aptos)
		assert.Empty(t, result.Solana)
		assert.Empty(t, result.Sui)
		assert.Empty(t, result.Ton)
		assert.Empty(t, result.Tron)

		evmChain, exists := result.Evm[999]
		assert.True(t, exists)
		assert.Equal(t, uint64(1234567890123456789), evmChain.ChainSelector)
		assert.Equal(t, "test-evm-chain", evmChain.ChainName)
	})

	runTestWithYaml(t, "Extra selectors: multiple families", yamlMultipleFamilies, func(t *testing.T, result extraSelectorsData) {
		assert.Len(t, result.Evm, 1)
		assert.Len(t, result.Aptos, 1)
		assert.Len(t, result.Sui, 1)
		assert.Len(t, result.Ton, 1)
		assert.Len(t, result.Tron, 1)

		aptosChain, exists := result.Aptos[888]
		assert.True(t, exists)
		assert.Equal(t, uint64(9876543210987654321), aptosChain.ChainSelector)
		assert.Equal(t, "test-aptos-chain", aptosChain.ChainName)

		suiChain, exists := result.Sui[777]
		assert.True(t, exists)
		assert.Equal(t, uint64(2222222222222222222), suiChain.ChainSelector)
		assert.Equal(t, "test-sui-chain", suiChain.ChainName)
	})

	runTestWithYaml(t, "Extra selectors: empty YAML file", ``, func(t *testing.T, result extraSelectorsData) {
		assert.Empty(t, result.Evm)
		assert.Empty(t, result.Aptos)
		assert.Empty(t, result.Solana)
		assert.Empty(t, result.Sui)
		assert.Empty(t, result.Ton)
		assert.Empty(t, result.Tron)
	})

}

func TestExtraSelectorsE2E(t *testing.T) {
	var err error
	var cwd string

	cwd, err = os.Getwd()
	require.NoError(t, err)
	testFilePath := filepath.Join(cwd, "test_extra_selectors.yml")
	envFilePath := filepath.Join(cwd, os.Getenv("EXTRA_SELECTORS_FILE"))

	if envFilePath != testFilePath {
		log.Printf("Skipping test because EXTRA_SELECTORS_FILE is not set to %s", testFilePath)
		return
	}

	var selector uint64
	var chainID uint64
	//Should return data for valid addition
	chainID, err = ChainIdFromSelector(1234567890123456789)
	assert.NoError(t, err)
	assert.Equal(t, uint64(90909090111), chainID)

	selector, err = SelectorFromChainId(90909090111)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1234567890123456789), selector)

	name, err := NameFromChainId(90909090111)
	assert.NoError(t, err)
	assert.Equal(t, "test-evm-chain", name)

	//Should not override chain id that already exists in selectors.yml
	selector, err = SelectorFromChainId(999)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2442541497099098535), selector)
	name, err = NameFromChainId(999)
	assert.NoError(t, err)
	assert.Equal(t, "hyperliquid-mainnet", name)
}
