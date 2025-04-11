//go:build ignore

package main

import (
	"fmt"
	"go/format"
	"os"
	"text/template"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chain-selectors/internal/gotmpl"
)

const filename = "generated_chains_aptos.go"

type chain struct {
	ChainID  uint64
	Selector uint64
	Name     string
	VarName  string
}

var chainTemplate, _ = template.New("").Parse(`// Code generated by go generate please DO NOT EDIT
package chain_selectors

type AptosChain struct {
	ChainID    uint64
	Selector   uint64
	Name       string
	VarName    string
}

var (
{{ range . }}
	{{.VarName}} = AptosChain{ChainID: {{ .ChainID }}, Selector: {{ .Selector }}, Name: "{{ .Name }}"}{{ end }}
)

var AptosALL = []AptosChain{
{{ range . }}{{ .VarName }},
{{ end }}
}

`)

func main() {
	src, err := gotmpl.Run(chainTemplate, chain_selectors.AptosChainIdToChainSelector, chain_selectors.AptosNameFromChainId)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source([]byte(src))
	if err != nil {
		panic(err)
	}

	existingContent, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if string(existingContent) == string(formatted) {
		fmt.Println("aptos: no changes detected")
		return
	}
	fmt.Println("aptos: updating generations")

	err = os.WriteFile(filename, formatted, 0644)
	if err != nil {
		panic(err)
	}
}
