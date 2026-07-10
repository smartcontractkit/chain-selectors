package chain_selectors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsMainnetChain(t *testing.T) {
	tests := []struct {
		name        string
		selector    uint64
		wantMainnet bool
		wantErr     bool
	}{
		// EVM Mainnets
		{
			name:        "Ethereum mainnet",
			selector:    ETHEREUM_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "BSC mainnet",
			selector:    BINANCE_SMART_CHAIN_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Polygon mainnet",
			selector:    POLYGON_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Arbitrum mainnet",
			selector:    ETHEREUM_MAINNET_ARBITRUM_1.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Optimism mainnet",
			selector:    ETHEREUM_MAINNET_OPTIMISM_1.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Base mainnet",
			selector:    ETHEREUM_MAINNET_BASE_1.Selector,
			wantMainnet: true,
			wantErr:     false,
		},

		// EVM Testnets
		{
			name:        "Ethereum Sepolia testnet",
			selector:    ETHEREUM_TESTNET_SEPOLIA.Selector,
			wantMainnet: false,
			wantErr:     false,
		},
		{
			name:        "BSC testnet",
			selector:    BINANCE_SMART_CHAIN_TESTNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},
		{
			name:        "Polygon Amoy testnet",
			selector:    POLYGON_TESTNET_AMOY.Selector,
			wantMainnet: false,
			wantErr:     false,
		},
		{
			name:        "Arbitrum Sepolia testnet",
			selector:    ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
			wantMainnet: false,
			wantErr:     false,
		},
		{
			name:        "Base Sepolia testnet",
			selector:    ETHEREUM_TESTNET_SEPOLIA_BASE_1.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Solana chains
		{
			name:        "Solana mainnet",
			selector:    SOLANA_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Solana testnet",
			selector:    SOLANA_TESTNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},
		{
			name:        "Solana devnet",
			selector:    SOLANA_DEVNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Aptos chains
		{
			name:        "Aptos mainnet",
			selector:    APTOS_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Aptos testnet",
			selector:    APTOS_TESTNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Sui chains
		{
			name:        "Sui mainnet",
			selector:    SUI_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Sui testnet",
			selector:    SUI_TESTNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Starknet chains
		{
			name:        "Starknet mainnet",
			selector:    ETHEREUM_MAINNET_STARKNET_1.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Starknet Sepolia testnet",
			selector:    ETHEREUM_TESTNET_SEPOLIA_STARKNET_1.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Canton chains
		{
			name:        "Canton mainnet",
			selector:    CANTON_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Canton testnet",
			selector:    CANTON_TESTNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// TON chains
		{
			name:        "TON mainnet",
			selector:    TON_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "TON testnet",
			selector:    TON_TESTNET.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Tron chains
		{
			name:        "Tron mainnet",
			selector:    TRON_MAINNET.Selector,
			wantMainnet: true,
			wantErr:     false,
		},
		{
			name:        "Tron testnet Nile",
			selector:    TRON_TESTNET_NILE.Selector,
			wantMainnet: false,
			wantErr:     false,
		},

		// Error cases
		{
			name:        "Unknown selector",
			selector:    9999999999999999999, // Invalid selector
			wantMainnet: false,
			wantErr:     true,
		},
		{
			name:        "Zero selector",
			selector:    0,
			wantMainnet: false,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsMainnetChain(tt.selector)

			if tt.wantErr {
				require.Error(t, err, "IsMainnetChain() should return error for selector %d", tt.selector)
				assert.Contains(t, err.Error(), "unknown chain selector", "Error message should indicate unknown selector")
			} else {
				require.NoError(t, err, "IsMainnetChain() should not return error for selector %d", tt.selector)
				assert.Equal(t, tt.wantMainnet, got, "IsMainnetChain() = %v, want %v for %s", got, tt.wantMainnet, tt.name)
			}
		})
	}
}

func TestGetChainDetailsByNetworkName(t *testing.T) {
	tests := []struct {
		name     string
		selector uint64
	}{
		{name: "EVM", selector: ETHEREUM_MAINNET.Selector},
		{name: "Solana", selector: SOLANA_MAINNET.Selector},
		{name: "Aptos", selector: APTOS_MAINNET.Selector},
		{name: "Sui", selector: SUI_MAINNET.Selector},
		{name: "Tron", selector: TRON_MAINNET.Selector},
		{name: "TON", selector: TON_MAINNET.Selector},
		{name: "Starknet", selector: ETHEREUM_MAINNET_STARKNET_1.Selector},
		{name: "Canton", selector: CANTON_MAINNET.Selector},
		{name: "Stellar", selector: STELLAR_MAINNET.Selector},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, err := GetChainDetails(tt.selector)
			require.NoError(t, err)

			got, err := GetChainDetailsByNetworkName(expected.ChainName)
			require.NoError(t, err)
			assert.Equal(t, expected, got)
		})
	}

	t.Run("unknown network name returns error", func(t *testing.T) {
		_, err := GetChainDetailsByNetworkName("unknown-network")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain details not found for network name unknown-network")
	})
}

func TestIsDeprecated(t *testing.T) {
	t.Run("known selector defaults to not deprecated", func(t *testing.T) {
		deprecated, err := IsDeprecated(ETHEREUM_MAINNET.Selector)
		require.NoError(t, err)
		assert.False(t, deprecated)
	})

	t.Run("unknown selector returns error", func(t *testing.T) {
		_, err := IsDeprecated(9999999999999999999)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown chain selector")
	})
}

// Parsing of both accepted formats plus empty/invalid.
func TestParseSunsetDate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantOK   bool
		wantTime time.Time
		wantErr  bool
	}{
		{"RFC 3339 datetime", "2020-01-02T03:04:05Z", true, time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC), false},
		{"date only", "2999-01-01", true, time.Date(2999, 1, 1, 0, 0, 0, 0, time.UTC), false},
		{"empty is not set", "", false, time.Time{}, false},
		{"invalid", "not-a-date", false, time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, ok, err := ParseSunsetDate(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid sunset date")
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantTime, ts)
		})
	}
}

// Sunset liveness: future = live, past = dead, unset = never dead.
func TestSunsetPassed(t *testing.T) {
	now := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		sunsetAt   string
		wantPassed bool
		wantErr    bool
	}{
		{"past date is dead", "2020-01-01T00:00:00Z", true, false},
		{"exact sunset time is dead", "2025-06-01T00:00:00Z", true, false},
		{"future date is live", "2999-01-01", false, false},
		{"empty is not sunset", "", false, false},
		{"invalid", "not-a-date", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passed, err := SunsetPassed(tt.sunsetAt, now)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid sunset date")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantPassed, passed)
		})
	}
}

// Registry wiring and error propagation for GetSunsetDate / IsSunset.
// (Past-date sunset is verified end-to-end in TestExtraSelectorsE2E.)
func TestSunsetAccessors(t *testing.T) {
	tests := []struct {
		name        string
		selector    uint64
		wantHasDate bool
		wantSunset  bool
		wantErr     bool
	}{
		{"chain without sunset date", ETHEREUM_MAINNET.Selector, false, false, false},
		{"unknown selector", 9999999999999999999, false, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok, err := GetSunsetDate(tt.selector)
			sunset, serr := IsSunset(tt.selector)
			if tt.wantErr {
				require.Error(t, err)
				require.Error(t, serr)
				return
			}
			require.NoError(t, err)
			require.NoError(t, serr)
			assert.Equal(t, tt.wantHasDate, ok)
			assert.Equal(t, tt.wantSunset, sunset)
		})
	}
}

// Every configured sunset_at must parse, so a malformed date fails CI and can't be merged.
func TestAllSunsetDatesValid(t *testing.T) {
	check := func(name, sunsetAt string) {
		if sunsetAt == "" {
			return
		}
		_, ok, err := ParseSunsetDate(sunsetAt)
		require.NoErrorf(t, err, "chain %q has invalid sunset_at %q", name, sunsetAt)
		assert.Truef(t, ok, "chain %q sunset_at %q parsed but reported not set", name, sunsetAt)
	}
	checkAll := func(selectors []ChainDetails) {
		for _, d := range selectors {
			check(d.ChainName, d.SunsetAt)
		}
	}
	values := func(m map[uint64]ChainDetails) []ChainDetails {
		out := make([]ChainDetails, 0, len(m))
		for _, d := range m {
			out = append(out, d)
		}
		return out
	}
	stringKeyed := func(m map[string]ChainDetails) []ChainDetails {
		out := make([]ChainDetails, 0, len(m))
		for _, d := range m {
			out = append(out, d)
		}
		return out
	}
	checkAll(values(evmChainIdToChainSelector))
	checkAll(stringKeyed(solanaChainIdToChainSelector))
	checkAll(values(aptosSelectorsMap))
	checkAll(values(suiSelectorsMap))
	checkAll(values(tronSelectorsMap))
	checkAll(stringKeyed(starknetSelectorsMap))
	checkAll(stringKeyed(cantonChainsByChainId))
	checkAll(stringKeyed(stellarChainsByChainId))
	for _, d := range tonSelectorsMap {
		check(d.ChainName, d.SunsetAt)
	}
}
