//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chain-selectors/internal/gotmpl"
)

func main() {
	gen("generated_chains.rs")
}

func wd() string {
	rsDir := os.Getenv("PWD")
	if !strings.HasSuffix(rsDir, "/rs") {
		rsDir = path.Join(rsDir, "rs")
	}
	return rsDir
}

func gen(generatedFileName string) {
	rsDir := wd()
	tmplFileName := strings.ReplaceAll(generatedFileName, ".rs", ".rs.tmpl")
	tmplRaw, err := os.ReadFile(path.Join(rsDir, tmplFileName))
	if err != nil {
		panic(fmt.Errorf("failed to read template file: %w", err))
	}
	tmpl, err := template.New(generatedFileName).Parse(string(tmplRaw))
	if err != nil {
		panic(fmt.Errorf("failed to parse template: %w", err))
	}

	generatedFilePath := path.Join(rsDir, "chainselectors", "src", generatedFileName)
	existingContent, err := os.ReadFile(generatedFilePath)
	if err != nil {
		panic(fmt.Errorf("failed to read existing file: %w", err))
	}

	raw, err := gotmpl.Run(tmpl, chain_selectors.EvmChainIdToChainSelector, chain_selectors.NameFromChainId)
	if err != nil {
		panic(fmt.Errorf("failed to run go-template: %w", err))
	}

	formatted, err := rustfmt([]byte(raw), generatedFileName)
	if err != nil {
		panic(fmt.Errorf("failed to format rust code: %w", err))
	}

	if string(existingContent) == string(formatted) {
		fmt.Println("evm (rust): no changes detected")
		return
	}

	if err := os.WriteFile(generatedFilePath, formatted, 0644); err != nil {
		panic(fmt.Errorf("failed to update %s: %w", generatedFileName, err))
	}
}

func rustfmt(src []byte, generatedFileName string) ([]byte, error) {
	tmpFile := path.Join(os.TempDir(), generatedFileName)
	if err := os.WriteFile(tmpFile, src, 0644); err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	cmd := exec.Command("rustfmt", tmpFile)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return nil, err
		}
		// if rustfmt is not installed, try to use docker
		vol := fmt.Sprintf("%s:/usr/src/app/%s", tmpFile, generatedFileName)
		rustfmtCmd := fmt.Sprintf("rustup component add rustfmt &>/dev/null && rustfmt %s", generatedFileName)
		if dockerErr := exec.Command(
			"docker", "run", "--rm",
			"-v", vol, "-w", "/usr/src/app",
			"rust:1.82-alpine",
			"/bin/sh", "-c", rustfmtCmd,
		).Run(); dockerErr != nil {
			return nil, err
		}
	}
	formatted, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, err
	}

	return formatted, nil
}
