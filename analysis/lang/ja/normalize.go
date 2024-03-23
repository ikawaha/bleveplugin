package ja

import (
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"golang.org/x/text/unicode/norm"
)

const NormalizeCharFilterName = "ja_normalize_unicode"

func init() {
	registry.RegisterCharFilter(NormalizeCharFilterName, UnicodeNormalizeCharFilterConstructor)
}

var forms = map[string]norm.Form{
	"nfc":  norm.NFC,
	"nfd":  norm.NFD,
	"nfkc": norm.NFKC,
	"nfkd": norm.NFKD,
}

// UnicodeNormalizeCharFilter represents unicode char filter.
type UnicodeNormalizeCharFilter struct {
	form norm.Form
}

// Filter applies per-char normalization.
func (f UnicodeNormalizeCharFilter) Filter(input []byte) []byte {
	return f.form.Bytes(input)
}

// NewUnicodeNormalizeCharFilter returns a normalize char filter.
func NewUnicodeNormalizeCharFilter(form norm.Form) analysis.CharFilter {
	return UnicodeNormalizeCharFilter{
		form: form,
	}
}

func UnicodeNormalizeCharFilterConstructor(config map[string]any, _ *registry.Cache) (analysis.CharFilter, error) {
	formVal, ok := config["form"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify form")
	}
	form, ok := forms[strings.ToLower(formVal)]
	if !ok {
		return nil, fmt.Errorf("no form named %s", formVal)
	}
	return NewUnicodeNormalizeCharFilter(form), nil
}
