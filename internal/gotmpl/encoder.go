package gotmpl

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mr-tron/base58"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	reSeperators = regexp.MustCompile("[-_]+")
	caser        = cases.Title(language.English)
)

func encodeVarName(name string, chainSel uint64) string {
	const unnamed = "TEST"
	x := strings.ReplaceAll(name, "-", "_")
	x = strings.ToUpper(x)
	if len(x) > 0 && (unicode.IsDigit(rune(x[0])) || isSolTestChain(name)) {
		x = unnamed + "_" + x
	}
	if len(x) == 0 {
		x = unnamed + "_" + strconv.FormatUint(chainSel, 10)
	}
	return x
}

func encodeEnumName(name string, chainSel uint64) string {
	const unnamed = "Test"
	x := reSeperators.ReplaceAllString(name, " ")
	varName := strings.ReplaceAll(caser.String(x), " ", "")
	if len(varName) > 0 && (unicode.IsDigit(rune(varName[0])) || isSolTestChain(name)) {
		varName = unnamed + varName
	}
	if len(varName) == 0 {
		varName = unnamed + strconv.FormatUint(chainSel, 10)
	}
	return varName
}

// for evm, the above condition is used to detect if name == chainId == (some number) -> which means its a test chain
// for solana, as chainId is not a number but a base58 encoded hash, we cannot use the above condition
// we need to check if the name == chainId == a valid base58 encoded hash
func isSolTestChain(name string) bool {
	_, err := base58.Decode(name)
	return err == nil
}
