//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/smartcontractkit/chain-selectors/internal/gotmpl"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	generatedFileName = "generated_chains.rs"
	tmplFileName      = "generated_chains.rs.tmpl"
)

func wd() string {
	rsDir := os.Getenv("PWD")
	if !strings.HasSuffix(rsDir, "/rs") {
		rsDir = path.Join(rsDir, "rs")
	}
	return rsDir
}

func main() {
	rsDir := wd()
	tmplRaw, err := os.ReadFile(path.Join(rsDir, tmplFileName))
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New(generatedFileName).Parse(string(tmplRaw))
	if err != nil {
		panic(err)
	}

	generatedFilePath := path.Join(rsDir, "chainselectors", "src", generatedFileName)
	existingContent, err := os.ReadFile(generatedFilePath)
	if err != nil {
		panic(err)
	}

	raw, err := gotmpl.Run(tmpl, newRustNameEncoder())
	if err != nil {
		panic(err)
	}

	formatted, err := rustfmt([]byte(raw))
	if err != nil {
		panic(err)
	}

	if string(existingContent) == string(formatted) {
		fmt.Println("rust: no changes detected")
		return
	}

	if err := os.WriteFile(generatedFilePath, formatted, 0644); err != nil {
		panic(err)
	}
}

func rustfmt(src []byte) ([]byte, error) {
	tmpFile := path.Join(os.TempDir(), generatedFileName)
	if err := os.WriteFile(tmpFile, src, 0644); err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	if err := exec.Command("rustfmt", tmpFile).Run(); err != nil {
		// if rustfmt is not installed, try to use docker
		cmd := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/usr/src/app/generated_chains.rs", tmpFile), "-w", "/usr/src/app", "rust:1.82-alpine", "/bin/sh", "-c", "rustup component add rustfmt &>/dev/null && rustfmt generated_chains.rs")
		if dockerErr := cmd.Run(); dockerErr != nil {
			return nil, err
		}
	}
	formatted, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

type rustNameEncoder struct {
	re    *regexp.Regexp
	caser cases.Caser
}

func newRustNameEncoder() *rustNameEncoder {
	return &rustNameEncoder{
		re:    regexp.MustCompile("[-_]+"),
		caser: cases.Title(language.English),
	}
}

func (enc *rustNameEncoder) VarName(name string, chainSel uint64) string {
	x := enc.re.ReplaceAllString(name, " ")
	varName := strings.ReplaceAll(enc.caser.String(x), " ", "")
	if len(varName) > 0 && unicode.IsDigit(rune(varName[0])) {
		varName = "Test" + varName
	}
	if len(varName) == 0 {
		varName = "Test" + strconv.FormatUint(chainSel, 10)
	}
	return varName
}
