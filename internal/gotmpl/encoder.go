package gotmpl

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func newNameEncoder() *nameEncoder {
	return &nameEncoder{
		re:    regexp.MustCompile("[-_]+"),
		caser: cases.Title(language.English),
	}
}

type nameEncoder struct {
	re    *regexp.Regexp
	caser cases.Caser
}

func (*nameEncoder) varName(name string, chainSel uint64) string {
	const unnamed = "TEST"
	x := strings.ReplaceAll(name, "-", "_")
	x = strings.ToUpper(x)
	if len(x) > 0 && unicode.IsDigit(rune(x[0])) {
		x = unnamed + "_" + x
	}
	if len(x) == 0 {
		x = unnamed + "_" + strconv.FormatUint(chainSel, 10)
	}
	return x
}

func (enc *nameEncoder) enumName(name string, chainSel uint64) string {
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
