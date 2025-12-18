package chain_selectors

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
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

func setSelectorEnv(t *testing.T, filePath string) {
	t.Setenv("EXTRA_SELECTORS_FILE", filePath)
}

// setRemoteDatasourceURL overrides the remoteDatasourceURL for testing
// Typically used with httptest.NewServer to create a mock HTTP server
// Usage: setRemoteDatasourceURL(t, server.URL)
func setRemoteDatasourceURL(t *testing.T, url string) {
	original := remoteDatasourceURL
	remoteDatasourceURL = url
	t.Cleanup(func() {
		remoteDatasourceURL = original
	})
}

func runTestWithYaml(t *testing.T, testName string, yamlContent string, validate func(*testing.T, extraSelectorsData)) {
	t.Run(testName, func(t *testing.T) {
		filePath := createTempYamlFile(t, yamlContent)
		defer os.Remove(filePath)

		setSelectorEnv(t, filePath)

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
  "ASwXBTzJM5evpfrWSHSjZaxPErZRuiGJnFixGUHi4NQT":  #Random solana chainID
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
starknet:
  "TEST_SN":
    selector: 1111111111111111111
    name: "test-starknet-chain"
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
		assert.Empty(t, result.Starknet)

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
		assert.Len(t, result.Starknet, 1)

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
		assert.Empty(t, result.Starknet)
	})

}

func TestExtraSelectorsInvalidFormat(t *testing.T) {
	t.Run("Invalid YAML syntax should panic", func(t *testing.T) {
		invalidYaml := `
evm:
  999:
    selector: 1234567890123456789
    name: "test-evm-chain"
invalid: yaml: syntax: [unclosed
`
		filePath := createTempYamlFile(t, invalidYaml)
		defer os.Remove(filePath)

		setSelectorEnv(t, filePath)

		assert.Panics(t, func() {
			loadAndParseExtraSelectors()
		}, "Expected panic for invalid YAML syntax")
	})

	t.Run("Invalid chain ID type should panic", func(t *testing.T) {
		invalidEvmChainIdYaml := `
evm:
  "abc":  # string
    selector: 1234567890123456789
    name: "test-evm-chain"
`
		filePath := createTempYamlFile(t, invalidEvmChainIdYaml)
		defer os.Remove(filePath)

		setSelectorEnv(t, filePath)

		assert.Panics(t, func() {
			loadAndParseExtraSelectors()
		}, "Expected panic for invalid EVM chain ID type")
	})

	t.Run("Non-existent file should panic", func(t *testing.T) {
		setSelectorEnv(t, "/non/existent/file.yaml")

		assert.Panics(t, func() {
			loadAndParseExtraSelectors()
		}, "Expected panic for non-existent file")
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
		t.Skipf("Skipping test because EXTRA_SELECTORS_FILE is not set to %s", testFilePath)
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
	chain, ok := ChainBySelector(2442541497099098535)
	assert.Equal(t, true, ok)
	assert.Equal(t, "hyperliquid-mainnet", chain.Name)
	assert.Equal(t, uint64(2442541497099098535), selector)
	name, err = NameFromChainId(999)
	assert.NoError(t, err)
	assert.Equal(t, "hyperliquid-mainnet", name)
}

// Validates a custom provide file for formating errors. This can be used in external CI checks to ensure the file is valid.
func TestExtraSelectorsValidateCustomFile(t *testing.T) {
	extraSelectorsFile := os.Getenv("EXTRA_SELECTORS_FILE")
	if extraSelectorsFile == "" {
		t.Skip("EXTRA_SELECTORS_FILE environment variable is not set")
		return
	}

	// loadAndParseExtraSelectors will panic if the file has errors
	assert.NotPanics(t, func() {
		loadAndParseExtraSelectors()
	}, "Loading extra selectors file should not panic if file is valid")
}

// ==================== Remote Datasource Tests ====================

func TestIsRemoteDatasourceEnabled(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		os.Unsetenv("ENABLE_REMOTE_DATASOURCE")
		assert.False(t, isRemoteDatasourceEnabled())
	})

	t.Run("enabled when set to true", func(t *testing.T) {
		t.Setenv("ENABLE_REMOTE_DATASOURCE", "true")
		assert.True(t, isRemoteDatasourceEnabled())
	})

	t.Run("enabled when set to 1", func(t *testing.T) {
		t.Setenv("ENABLE_REMOTE_DATASOURCE", "1")
		assert.True(t, isRemoteDatasourceEnabled())
	})

	t.Run("disabled when set to false", func(t *testing.T) {
		t.Setenv("ENABLE_REMOTE_DATASOURCE", "false")
		assert.False(t, isRemoteDatasourceEnabled())
	})

	t.Run("disabled when set to invalid value", func(t *testing.T) {
		t.Setenv("ENABLE_REMOTE_DATASOURCE", "invalid")
		assert.False(t, isRemoteDatasourceEnabled())
	})
}

func TestGetRemoteChainByID(t *testing.T) {
	// Save original state
	originalRemoteSelectors := remoteSelectors
	originalRemoteSelectorsFetched := remoteSelectorsFetched
	originalOnce := remoteSelectorsOnce

	// Cleanup runs after all subtests complete
	t.Cleanup(func() {
		remoteSelectors = originalRemoteSelectors
		remoteSelectorsFetched = originalRemoteSelectorsFetched
		remoteSelectorsOnce = originalOnce
	})

	// Reset sync.Once and mark as already fetched to prevent network calls
	remoteSelectorsOnce = sync.Once{}
	remoteSelectorsOnce.Do(func() {}) // Mark as done

	// Setup test data
	remoteSelectors = extraSelectorsData{
		Evm: map[uint64]ChainDetails{
			12345: {ChainSelector: 9999999999, ChainName: "test-remote-evm"},
		},
		Solana: map[string]ChainDetails{
			"TestSolanaChainID": {ChainSelector: 8888888888, ChainName: "test-remote-solana"},
		},
		Aptos: map[uint64]ChainDetails{
			111: {ChainSelector: 7777777777, ChainName: "test-remote-aptos"},
		},
	}
	remoteSelectorsFetched = true

	// Enable remote datasource
	t.Setenv("ENABLE_REMOTE_DATASOURCE", "true")

	t.Run("EVM chain found", func(t *testing.T) {
		details, ok := getRemoteChainByID(FamilyEVM, "12345")
		assert.True(t, ok)
		assert.Equal(t, uint64(9999999999), details.ChainSelector)
		assert.Equal(t, "test-remote-evm", details.ChainName)
	})

	t.Run("EVM chain not found", func(t *testing.T) {
		_, ok := getRemoteChainByID(FamilyEVM, "99999")
		assert.False(t, ok)
	})

	t.Run("Solana chain found", func(t *testing.T) {
		details, ok := getRemoteChainByID(FamilySolana, "TestSolanaChainID")
		assert.True(t, ok)
		assert.Equal(t, uint64(8888888888), details.ChainSelector)
		assert.Equal(t, "test-remote-solana", details.ChainName)
	})

	t.Run("Aptos chain found", func(t *testing.T) {
		details, ok := getRemoteChainByID(FamilyAptos, "111")
		assert.True(t, ok)
		assert.Equal(t, uint64(7777777777), details.ChainSelector)
		assert.Equal(t, "test-remote-aptos", details.ChainName)
	})

	t.Run("Invalid chain ID format", func(t *testing.T) {
		_, ok := getRemoteChainByID(FamilyEVM, "not-a-number")
		assert.False(t, ok)
	})

	t.Run("Unknown family", func(t *testing.T) {
		_, ok := getRemoteChainByID("unknown-family", "12345")
		assert.False(t, ok)
	})
}

func TestGetRemoteChainBySelector(t *testing.T) {
	// Save original state
	originalRemoteSelectors := remoteSelectors
	originalRemoteSelectorsFetched := remoteSelectorsFetched
	originalOnce := remoteSelectorsOnce

	// Cleanup runs after all subtests complete
	t.Cleanup(func() {
		remoteSelectors = originalRemoteSelectors
		remoteSelectorsFetched = originalRemoteSelectorsFetched
		remoteSelectorsOnce = originalOnce
	})

	// Reset sync.Once and mark as already fetched to prevent network calls
	remoteSelectorsOnce = sync.Once{}
	remoteSelectorsOnce.Do(func() {}) // Mark as done

	// Setup test data
	remoteSelectors = extraSelectorsData{
		Evm: map[uint64]ChainDetails{
			12345: {ChainSelector: 9999999999, ChainName: "test-remote-evm"},
		},
		Solana: map[string]ChainDetails{
			"TestSolanaChainID": {ChainSelector: 8888888888, ChainName: "test-remote-solana"},
		},
		Aptos: map[uint64]ChainDetails{
			111: {ChainSelector: 7777777777, ChainName: "test-remote-aptos"},
		},
	}
	remoteSelectorsFetched = true

	// Enable remote datasource
	t.Setenv("ENABLE_REMOTE_DATASOURCE", "true")

	t.Run("EVM selector found", func(t *testing.T) {
		family, chainID, details, ok := getRemoteChainBySelector(9999999999)
		assert.True(t, ok)
		assert.Equal(t, FamilyEVM, family)
		assert.Equal(t, "12345", chainID)
		assert.Equal(t, "test-remote-evm", details.ChainName)
	})

	t.Run("Solana selector found", func(t *testing.T) {
		family, chainID, details, ok := getRemoteChainBySelector(8888888888)
		assert.True(t, ok)
		assert.Equal(t, FamilySolana, family)
		assert.Equal(t, "TestSolanaChainID", chainID)
		assert.Equal(t, "test-remote-solana", details.ChainName)
	})

	t.Run("Aptos selector found", func(t *testing.T) {
		family, chainID, details, ok := getRemoteChainBySelector(7777777777)
		assert.True(t, ok)
		assert.Equal(t, FamilyAptos, family)
		assert.Equal(t, "111", chainID)
		assert.Equal(t, "test-remote-aptos", details.ChainName)
	})

	t.Run("Selector not found", func(t *testing.T) {
		_, _, _, ok := getRemoteChainBySelector(1111111111)
		assert.False(t, ok)
	})
}

func TestTryLazyFetchRemoteSelectorsDisabled(t *testing.T) {
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Should return false when disabled
	result := tryLazyFetchRemoteSelectors()
	assert.False(t, result)
}

func TestLoadRemoteDatasource(t *testing.T) {
	// Test with a mock HTTP server
	yamlContent := `
evm:
  99999:
    selector: 1234567890123456789
    name: "mock-evm-chain"
solana:
  "MockSolanaChainID":
    selector: 9876543210987654321
    name: "mock-solana-chain"
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(yamlContent))
	}))
	t.Cleanup(func() { server.Close() })

	// We can't easily test loadRemoteDatasource directly because it uses a hardcoded URL
	// But we can test fetchFromURL which it uses internally
	t.Run("fetchFromURL success", func(t *testing.T) {
		content, err := fetchFromURL(server.URL)
		require.NoError(t, err)
		assert.Contains(t, string(content), "mock-evm-chain")
	})

	t.Run("fetchFromURL 404", func(t *testing.T) {
		server404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		t.Cleanup(func() { server404.Close() })

		_, err := fetchFromURL(server404.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("fetchFromURL 500", func(t *testing.T) {
		server500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(func() { server500.Close() })

		_, err := fetchFromURL(server500.URL)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})
}

func TestRemoteFallbackInLookupFunctions(t *testing.T) {
	// Save original state
	originalRemoteSelectors := remoteSelectors
	originalRemoteSelectorsFetched := remoteSelectorsFetched
	originalOnce := remoteSelectorsOnce

	// Cleanup runs after all subtests complete
	t.Cleanup(func() {
		remoteSelectors = originalRemoteSelectors
		remoteSelectorsFetched = originalRemoteSelectorsFetched
		remoteSelectorsOnce = originalOnce
	})

	// Reset sync.Once and mark as already fetched to prevent network calls
	remoteSelectorsOnce = sync.Once{}
	remoteSelectorsOnce.Do(func() {}) // Mark as done

	// Use chain IDs that definitely don't exist in embedded chains
	// (very high numbers that won't conflict with any real chain)
	testChainID := uint64(88888888888)
	testSelector := uint64(7777777777777777777)

	// Setup test data
	remoteSelectors = extraSelectorsData{
		Evm: map[uint64]ChainDetails{
			testChainID: {ChainSelector: testSelector, ChainName: "remote-only-chain"},
		},
	}
	remoteSelectorsFetched = true

	// Enable remote datasource
	t.Setenv("ENABLE_REMOTE_DATASOURCE", "true")

	t.Run("SelectorFromChainId falls back to remote", func(t *testing.T) {
		selector, err := SelectorFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, testSelector, selector)
	})

	t.Run("ChainIdFromSelector falls back to remote", func(t *testing.T) {
		chainID, err := ChainIdFromSelector(testSelector)
		require.NoError(t, err)
		assert.Equal(t, testChainID, chainID)
	})

	t.Run("NameFromChainId falls back to remote", func(t *testing.T) {
		name, err := NameFromChainId(testChainID)
		require.NoError(t, err)
		assert.Equal(t, "remote-only-chain", name)
	})

	t.Run("ChainBySelector falls back to remote", func(t *testing.T) {
		chain, ok := ChainBySelector(testSelector)
		assert.True(t, ok)
		assert.Equal(t, "remote-only-chain", chain.Name)
		assert.Equal(t, testChainID, chain.EvmChainID)
	})

	t.Run("ChainByEvmChainID falls back to remote", func(t *testing.T) {
		chain, ok := ChainByEvmChainID(testChainID)
		assert.True(t, ok)
		assert.Equal(t, "remote-only-chain", chain.Name)
	})
}

func TestRemoteDisabledDoesNotFallback(t *testing.T) {
	// Make sure remote datasource is disabled
	os.Unsetenv("ENABLE_REMOTE_DATASOURCE")

	// Use chain IDs that definitely don't exist
	nonExistentChainID := uint64(77777777777)
	nonExistentSelector := uint64(6666666666666666666)

	t.Run("SelectorFromChainId returns error when not found", func(t *testing.T) {
		_, err := SelectorFromChainId(nonExistentChainID)
		assert.Error(t, err)
	})

	t.Run("ChainBySelector returns false when not found", func(t *testing.T) {
		_, ok := ChainBySelector(nonExistentSelector)
		assert.False(t, ok)
	})
}

// Example test demonstrating how to mock the remote datasource using httptest.NewServer
func TestLoadRemoteDatasourceWithMockServer(t *testing.T) {
	// Create a mock HTTP server
	mockYAML := `
evm:
  999888:
    selector: 9876543210
    name: "test-mock-chain"
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockYAML))
	}))
	t.Cleanup(server.Close)

	// Override the remote datasource URL to point to our mock server
	setRemoteDatasourceURL(t, server.URL)

	// Reset remote state for this test
	originalRemoteSelectors := remoteSelectors
	originalRemoteSelectorsFetched := remoteSelectorsFetched
	originalOnce := remoteSelectorsOnce
	t.Cleanup(func() {
		remoteSelectors = originalRemoteSelectors
		remoteSelectorsFetched = originalRemoteSelectorsFetched
		remoteSelectorsOnce = originalOnce
	})

	remoteSelectorsOnce = sync.Once{}
	remoteSelectorsFetched = false

	// Enable remote datasource
	t.Setenv("ENABLE_REMOTE_DATASOURCE", "true")

	// Trigger the remote datasource load
	result := loadRemoteDatasource()

	// Verify the mock data was loaded
	require.NotNil(t, result.Evm)
	details, exists := result.Evm[999888]
	require.True(t, exists, "Expected chain ID 999888 to exist in result")
	assert.Equal(t, uint64(9876543210), details.ChainSelector)
	assert.Equal(t, "test-mock-chain", details.ChainName)
}
