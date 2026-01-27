package chain_selectors

import (
	"testing"

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
