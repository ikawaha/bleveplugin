package ja

import (
	"github.com/blevesearch/bleve/v2/analysis"
)

const (
	Inf動詞   = "動詞"   //nolint:gosmopolitan,asciicheck
	Inf形容詞  = "形容詞"  //nolint:gosmopolitan,asciicheck
	Inf形容動詞 = "形容動詞" //nolint:gosmopolitan,asciicheck
)

// DefaultInflected represents POSs which has inflected form.
var DefaultInflected = analysis.TokenMap{
	Inf動詞:   true,
	Inf形容詞:  true,
	Inf形容動詞: true,
}
