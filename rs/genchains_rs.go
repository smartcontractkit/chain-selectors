//go:build ignore

package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

func main() {
	tmplRaw, err := os.ReadFile("rs/generated_chains.rs.tmpl")
	if err != nil {
		panic(err)
	}
	chains, err := readChainsFromSelectors()
	if err != nil {
		panic(err)
	}

	tmpl, err := template.New("generated_chains.rs").Parse(string(tmplRaw))
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("rs/chainselectors/src/generated_chains.rs", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := tmpl.Execute(f, chains); err != nil {
		panic(err)
	}
}

type SelectorsYamlEntry struct {
	Name     string `yaml:"name"`
	Selector uint64 `yaml:"selector"`
}

type SelectorsYaml struct {
	Selectors map[uint64]SelectorsYamlEntry `yaml:"selectors"`
}

type Chain struct {
	EvmChainID uint64
	Selector   uint64
	Name       string
	VarName    string
}

func readSelectorsYaml(filePath string) (*SelectorsYaml, error) {
	selectorsRaw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read selectors yml: %w", err)
	}
	var selectors SelectorsYaml
	err = yaml.Unmarshal(selectorsRaw, &selectors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse selectors yml: %w", err)
	}
	return &selectors, nil
}

func readChainsFromSelectors() ([]Chain, error) {
	selectors, err := readSelectorsYaml(path.Join( /*"..", */ "selectors.yml"))
	if err != nil {
		return nil, err
	}
	testSelectors, err := readSelectorsYaml(path.Join( /*"..", */ "test_selectors.yml"))
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("[-_]+")
	chains := make([]Chain, 0, len(selectors.Selectors)+len(testSelectors.Selectors))
	for chainID, selector := range selectors.Selectors {
		chains = append(chains, Chain{
			EvmChainID: chainID,
			Selector:   selector.Selector,
			Name:       selector.Name,
			VarName:    toVarName(selector.Name, selector.Selector, re),
		})
	}
	return chains, nil
}

func toVarName(name string, chainSel uint64, reSep *regexp.Regexp) string {
	x := reSep.ReplaceAllString(name, " ")
	varName := strings.ReplaceAll(strings.Title(x), " ", "")
	if len(varName) > 0 && unicode.IsDigit(rune(varName[0])) {
		varName = "Test" + varName
	}
	if len(varName) == 0 {
		varName = "Test" + strconv.FormatUint(chainSel, 10)
	}
	return varName
}
