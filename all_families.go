package chain_selectors

import (
	"fmt"
	"slices"
)

// familySelectors maps each chain family to the selectors registered for it. It reads the same
// indexes getChainInfo uses, so it includes test selectors and extra selectors loaded from
// EXTRA_SELECTORS_FILE. When adding a family, add it here as well as in getChainInfo;
// TestFamilySelectorsMatchGetSelectorFamily checks that both agree.
var familySelectors = map[string]func() []uint64{
	FamilyEVM:      func() []uint64 { return mapKeys(evmChainsBySelector) },
	FamilySolana:   func() []uint64 { return mapKeys(solanaChainsBySelector) },
	FamilyAptos:    func() []uint64 { return mapKeys(aptosChainsBySelector) },
	FamilySui:      func() []uint64 { return mapKeys(suiChainsBySelector) },
	FamilyTron:     func() []uint64 { return mapKeys(tronChainIdBySelector) },
	FamilyTon:      func() []uint64 { return mapKeys(tonChainIdBySelector) },
	FamilyStarknet: func() []uint64 { return mapKeys(starknetChainsBySelector) },
	FamilyCanton:   func() []uint64 { return mapKeys(cantonChainsBySelector) },
	FamilyStellar:  func() []uint64 { return mapKeys(stellarChainsBySelector) },
}

// Families returns the chain families that have registered selectors, sorted.
func Families() []string {
	families := make([]string, 0, len(familySelectors))
	for family := range familySelectors {
		families = append(families, family)
	}
	slices.Sort(families)
	return families
}

// SelectorsByFamily returns every chain selector registered for the family, sorted. It returns an
// error for a family with no registered selectors.
func SelectorsByFamily(family string) ([]uint64, error) {
	selectors, ok := familySelectors[family]
	if !ok {
		return nil, fmt.Errorf("unknown chain family %q", family)
	}
	return sorted(selectors()), nil
}

// AllSelectors returns every registered chain selector across all families, sorted.
func AllSelectors() []uint64 {
	var all []uint64
	for _, selectors := range familySelectors {
		all = append(all, selectors()...)
	}
	return sorted(all)
}

func mapKeys[V any](m map[uint64]V) []uint64 {
	keys := make([]uint64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func sorted(selectors []uint64) []uint64 {
	slices.Sort(selectors)
	return selectors
}
