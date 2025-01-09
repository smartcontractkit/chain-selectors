package gotmpl

import (
	"bytes"
	"sort"
	"text/template"
)

// chain is a generic struct that can be used to represent chain families.
// C is the type of the chain ID.
//
// Supported types:
// EVM: uint64
// Solana: string
// Aptos: uint64
type chain[C uint64 | string] struct {
	ChainID  C
	Selector uint64
	Name     string
	VarName  string
	EnumName string
}

// Run runs the template with the given chains and returns the result.
// C is the type of the chain ID.
func Run[C uint64 | string](tmpl *template.Template, chainSelFunc func() map[C]uint64, nameFunc func(C) (string, error)) (string, error) {
	chains := make([]chain[C], 0)

	for chainID, chainSel := range chainSelFunc() {
		name, err := nameFunc(chainID)
		if err != nil {
			return "", err
		}

		chains = append(chains, chain[C]{
			ChainID:  chainID,
			Selector: chainSel,
			Name:     name,
			VarName:  encodeVarName(name, chainSel),
			EnumName: encodeEnumName(name, chainSel),
		})
	}

	sort.Slice(chains, func(i, j int) bool { return chains[i].VarName < chains[j].VarName })

	var wr = new(bytes.Buffer)

	if err := tmpl.Execute(wr, chains); err != nil {
		return "", err
	}

	return wr.String(), nil
}
