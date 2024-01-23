// Code generated by go generate please DO NOT EDIT
package chain_selectors

type Chain struct {
	EvmChainID uint64
	Selector   uint64
	Name       string
	VarName    string
}

var (
	AVALANCHE_MAINNET                       = Chain{EvmChainID: 43114, Selector: 6433500567565415381, Name: "avalanche-mainnet"}
	AVALANCHE_TESTNET_FUJI                  = Chain{EvmChainID: 43113, Selector: 14767482510784806043, Name: "avalanche-testnet-fuji"}
	BINANCE_SMART_CHAIN_MAINNET             = Chain{EvmChainID: 56, Selector: 11344663589394136015, Name: "binance_smart_chain-mainnet"}
	BINANCE_SMART_CHAIN_TESTNET             = Chain{EvmChainID: 97, Selector: 13264668187771770619, Name: "binance_smart_chain-testnet"}
	BITTORRENT_CHAIN_MAINNET                = Chain{EvmChainID: 199, Selector: 3776006016387883143, Name: "bittorrent_chain-mainnet"}
	BITTORRENT_CHAIN_TESTNET                = Chain{EvmChainID: 1029, Selector: 4459371029167934217, Name: "bittorrent_chain-testnet"}
	ETHEREUM_MAINNET                        = Chain{EvmChainID: 1, Selector: 5009297550715157269, Name: "ethereum-mainnet"}
	ETHEREUM_MAINNET_ARBITRUM_1             = Chain{EvmChainID: 42161, Selector: 4949039107694359620, Name: "ethereum-mainnet-arbitrum-1"}
	ETHEREUM_MAINNET_BASE_1                 = Chain{EvmChainID: 8453, Selector: 15971525489660198786, Name: "ethereum-mainnet-base-1"}
	ETHEREUM_MAINNET_KROMA_1                = Chain{EvmChainID: 255, Selector: 3719320017875267166, Name: "ethereum-mainnet-kroma-1"}
	ETHEREUM_MAINNET_MANTLE_1               = Chain{EvmChainID: 5000, Selector: 1556008542357238666, Name: "ethereum-mainnet-mantle-1"}
	ETHEREUM_MAINNET_OPTIMISM_1             = Chain{EvmChainID: 10, Selector: 3734403246176062136, Name: "ethereum-mainnet-optimism-1"}
	ETHEREUM_MAINNET_POLYGON_ZKEVM_1        = Chain{EvmChainID: 1101, Selector: 4348158687435793198, Name: "ethereum-mainnet-polygon-zkevm-1"}
	ETHEREUM_MAINNET_SCROLL_1               = Chain{EvmChainID: 534352, Selector: 13204309965629103672, Name: "ethereum-mainnet-scroll-1"}
	ETHEREUM_TESTNET_GOERLI_ARBITRUM_1      = Chain{EvmChainID: 421613, Selector: 6101244977088475029, Name: "ethereum-testnet-goerli-arbitrum-1"}
	ETHEREUM_TESTNET_GOERLI_BASE_1          = Chain{EvmChainID: 84531, Selector: 5790810961207155433, Name: "ethereum-testnet-goerli-base-1"}
	ETHEREUM_TESTNET_GOERLI_MANTLE_1        = Chain{EvmChainID: 5001, Selector: 1226473277236831298, Name: "ethereum-testnet-goerli-mantle-1"}
	ETHEREUM_TESTNET_GOERLI_OPTIMISM_1      = Chain{EvmChainID: 420, Selector: 2664363617261496610, Name: "ethereum-testnet-goerli-optimism-1"}
	ETHEREUM_TESTNET_GOERLI_POLYGON_ZKEVM_1 = Chain{EvmChainID: 1442, Selector: 11059667695644972511, Name: "ethereum-testnet-goerli-polygon-zkevm-1"}
	ETHEREUM_TESTNET_GOERLI_ZKSYNC_1        = Chain{EvmChainID: 280, Selector: 6802309497652714138, Name: "ethereum-testnet-goerli-zksync-1"}
	ETHEREUM_TESTNET_SEPOLIA                = Chain{EvmChainID: 11155111, Selector: 16015286601757825753, Name: "ethereum-testnet-sepolia"}
	ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1     = Chain{EvmChainID: 421614, Selector: 3478487238524512106, Name: "ethereum-testnet-sepolia-arbitrum-1"}
	ETHEREUM_TESTNET_SEPOLIA_BASE_1         = Chain{EvmChainID: 84532, Selector: 10344971235874465080, Name: "ethereum-testnet-sepolia-base-1"}
	ETHEREUM_TESTNET_SEPOLIA_KROMA_1        = Chain{EvmChainID: 2358, Selector: 5990477251245693094, Name: "ethereum-testnet-sepolia-kroma-1"}
	ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1     = Chain{EvmChainID: 11155420, Selector: 5224473277236331295, Name: "ethereum-testnet-sepolia-optimism-1"}
	ETHEREUM_TESTNET_SEPOLIA_SCROLL_1       = Chain{EvmChainID: 534351, Selector: 2279865765895943307, Name: "ethereum-testnet-sepolia-scroll-1"}
	KAVA_MAINNET                            = Chain{EvmChainID: 2222, Selector: 7550000543357438061, Name: "kava-mainnet"}
	KAVA_TESTNET                            = Chain{EvmChainID: 2221, Selector: 2110537777356199208, Name: "kava-testnet"}
	POLYGON_MAINNET                         = Chain{EvmChainID: 137, Selector: 4051577828743386545, Name: "polygon-mainnet"}
	POLYGON_TESTNET_MUMBAI                  = Chain{EvmChainID: 80001, Selector: 12532609583862916517, Name: "polygon-testnet-mumbai"}
	TEST_1000                               = Chain{EvmChainID: 1000, Selector: 11787463284727550157, Name: "1000"}
	TEST_1337                               = Chain{EvmChainID: 1337, Selector: 3379446385462418246, Name: "1337"}
	TEST_2337                               = Chain{EvmChainID: 2337, Selector: 12922642891491394802, Name: "2337"}
	TEST_76578                              = Chain{EvmChainID: 76578, Selector: 781901677223027175, Name: "76578"}
	TEST_90000001                           = Chain{EvmChainID: 90000001, Selector: 909606746561742123, Name: "90000001"}
	TEST_90000002                           = Chain{EvmChainID: 90000002, Selector: 5548718428018410741, Name: "90000002"}
	TEST_90000003                           = Chain{EvmChainID: 90000003, Selector: 789068866484373046, Name: "90000003"}
	TEST_90000004                           = Chain{EvmChainID: 90000004, Selector: 5721565186521185178, Name: "90000004"}
	TEST_90000005                           = Chain{EvmChainID: 90000005, Selector: 964127714438319834, Name: "90000005"}
	TEST_90000006                           = Chain{EvmChainID: 90000006, Selector: 8966794841936584464, Name: "90000006"}
	TEST_90000007                           = Chain{EvmChainID: 90000007, Selector: 8412806778050735057, Name: "90000007"}
	TEST_90000008                           = Chain{EvmChainID: 90000008, Selector: 4066443121807923198, Name: "90000008"}
	TEST_90000009                           = Chain{EvmChainID: 90000009, Selector: 6747736380229414777, Name: "90000009"}
	TEST_90000010                           = Chain{EvmChainID: 90000010, Selector: 8694984074292254623, Name: "90000010"}
	TEST_90000011                           = Chain{EvmChainID: 90000011, Selector: 328334718812072308, Name: "90000011"}
	TEST_90000012                           = Chain{EvmChainID: 90000012, Selector: 7715160997071429212, Name: "90000012"}
	TEST_90000013                           = Chain{EvmChainID: 90000013, Selector: 3574539439524578558, Name: "90000013"}
	TEST_90000014                           = Chain{EvmChainID: 90000014, Selector: 4543928599863227519, Name: "90000014"}
	TEST_90000015                           = Chain{EvmChainID: 90000015, Selector: 6443235356619661032, Name: "90000015"}
	TEST_90000016                           = Chain{EvmChainID: 90000016, Selector: 13087962012083037329, Name: "90000016"}
	TEST_90000017                           = Chain{EvmChainID: 90000017, Selector: 11985232338641871056, Name: "90000017"}
	TEST_90000018                           = Chain{EvmChainID: 90000018, Selector: 7777066535355430289, Name: "90000018"}
	TEST_90000019                           = Chain{EvmChainID: 90000019, Selector: 1273605685587320666, Name: "90000019"}
	TEST_90000020                           = Chain{EvmChainID: 90000020, Selector: 17810359353458878177, Name: "90000020"}
	WEMIX_MAINNET                           = Chain{EvmChainID: 1111, Selector: 5142893604156789321, Name: "wemix-mainnet"}
	WEMIX_TESTNET                           = Chain{EvmChainID: 1112, Selector: 9284632837123596123, Name: "wemix-testnet"}
)

var ALL = []Chain{
	AVALANCHE_MAINNET,
	AVALANCHE_TESTNET_FUJI,
	BINANCE_SMART_CHAIN_MAINNET,
	BINANCE_SMART_CHAIN_TESTNET,
	BITTORRENT_CHAIN_MAINNET,
	BITTORRENT_CHAIN_TESTNET,
	ETHEREUM_MAINNET,
	ETHEREUM_MAINNET_ARBITRUM_1,
	ETHEREUM_MAINNET_BASE_1,
	ETHEREUM_MAINNET_KROMA_1,
	ETHEREUM_MAINNET_MANTLE_1,
	ETHEREUM_MAINNET_OPTIMISM_1,
	ETHEREUM_MAINNET_POLYGON_ZKEVM_1,
	ETHEREUM_MAINNET_SCROLL_1,
	ETHEREUM_TESTNET_GOERLI_ARBITRUM_1,
	ETHEREUM_TESTNET_GOERLI_BASE_1,
	ETHEREUM_TESTNET_GOERLI_MANTLE_1,
	ETHEREUM_TESTNET_GOERLI_OPTIMISM_1,
	ETHEREUM_TESTNET_GOERLI_POLYGON_ZKEVM_1,
	ETHEREUM_TESTNET_GOERLI_ZKSYNC_1,
	ETHEREUM_TESTNET_SEPOLIA,
	ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1,
	ETHEREUM_TESTNET_SEPOLIA_BASE_1,
	ETHEREUM_TESTNET_SEPOLIA_KROMA_1,
	ETHEREUM_TESTNET_SEPOLIA_OPTIMISM_1,
	ETHEREUM_TESTNET_SEPOLIA_SCROLL_1,
	KAVA_MAINNET,
	KAVA_TESTNET,
	POLYGON_MAINNET,
	POLYGON_TESTNET_MUMBAI,
	TEST_1000,
	TEST_1337,
	TEST_2337,
	TEST_76578,
	TEST_90000001,
	TEST_90000002,
	TEST_90000003,
	TEST_90000004,
	TEST_90000005,
	TEST_90000006,
	TEST_90000007,
	TEST_90000008,
	TEST_90000009,
	TEST_90000010,
	TEST_90000011,
	TEST_90000012,
	TEST_90000013,
	TEST_90000014,
	TEST_90000015,
	TEST_90000016,
	TEST_90000017,
	TEST_90000018,
	TEST_90000019,
	TEST_90000020,
	WEMIX_MAINNET,
	WEMIX_TESTNET,
}
