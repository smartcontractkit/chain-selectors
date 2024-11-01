//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"gopkg.in/yaml.v3"
)

func main() {
	rsDir := os.Getenv("PWD")
	if !strings.HasSuffix(rsDir, "/rs") {
		rsDir = path.Join(rsDir, "rs")
	}
	tmplRaw, err := os.ReadFile(path.Join(rsDir, "generated_chains.rs.tmpl"))
	if err != nil {
		panic(err)
	}
	chains, err := readChainsFromSelectors(
		path.Join(rsDir, "..", "selectors.yml"),
		path.Join(rsDir, "..", "test_selectors.yml"),
	)
	if err != nil {
		panic(err)
	}

	generatedFileName := "generated_chains.rs"
	tmpl, err := template.New(generatedFileName).Parse(string(tmplRaw))
	if err != nil {
		panic(err)
	}

	generatedFilePath := path.Join(rsDir, "chainselectors", "src", generatedFileName)
	existingContent, err := os.ReadFile(generatedFilePath)
	if err != nil {
		panic(err)
	}
	var wr = new(bytes.Buffer)
	if err := tmpl.Execute(wr, chains); err != nil {
		panic(err)
	}
	tmpFile := path.Join(os.TempDir(), generatedFileName)
	if err := os.WriteFile(tmpFile, wr.Bytes(), 0644); err != nil {
		panic(err)
	}
	// execute rustfmt on the temporary generated file
	cmd := exec.Command("rustfmt", tmpFile)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	formatted, err := os.ReadFile(tmpFile)
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpFile)

	if string(existingContent) == string(formatted) {
		fmt.Println("rust: no changes detected")
		return
	}
	if err := os.WriteFile(generatedFilePath, formatted, 0644); err != nil {
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

func readChainsFromSelectors(selectorsYml, testSelectorsYml string) ([]Chain, error) {
	selectors, err := readSelectorsYaml(selectorsYml)
	if err != nil {
		return nil, err
	}
	testSelectors, err := readSelectorsYaml(testSelectorsYml)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("[-_]+")
	caser := cases.Title(language.English)
	chains := make([]Chain, 0, len(selectors.Selectors)+len(testSelectors.Selectors))
	for chainID, selector := range selectors.Selectors {
		chains = append(chains, Chain{
			EvmChainID: chainID,
			Selector:   selector.Selector,
			Name:       selector.Name,
			VarName:    toVarName(selector.Name, selector.Selector, caser, re),
		})
	}

	sort.Slice(chains, func(i, j int) bool { return chains[i].VarName < chains[j].VarName })

	return chains, nil
}

func toVarName(name string, chainSel uint64, caser cases.Caser, reSep *regexp.Regexp) string {
	x := reSep.ReplaceAllString(name, " ")
	varName := strings.ReplaceAll(caser.String(x), " ", "")
	if len(varName) > 0 && unicode.IsDigit(rune(varName[0])) {
		varName = "Test" + varName
	}
	if len(varName) == 0 {
		varName = "Test" + strconv.FormatUint(chainSel, 10)
	}
	return varName
}
