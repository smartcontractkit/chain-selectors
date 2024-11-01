package gotmpl

import (
	"bytes"
	"sort"
	"text/template"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
)

type NameEncoder interface {
	VarName(name string, chainSel uint64) string
}

type chain struct {
	EvmChainID uint64
	Selector   uint64
	Name       string
	VarName    string
}

func Run(tmpl *template.Template, enc NameEncoder) (string, error) {
	var wr = new(bytes.Buffer)
	chains := make([]chain, 0)

	for evmChainID, chainSel := range chain_selectors.EvmChainIdToChainSelector() {
		name, err := chain_selectors.NameFromChainId(evmChainID)
		if err != nil {
			return "", err
		}

		chains = append(chains, chain{
			EvmChainID: evmChainID,
			Selector:   chainSel,
			Name:       name,
			VarName:    enc.VarName(name, chainSel),
		})
	}

	sort.Slice(chains, func(i, j int) bool { return chains[i].VarName < chains[j].VarName })

	if err := tmpl.Execute(wr, chains); err != nil {
		return "", err
	}

	return wr.String(), nil
}
