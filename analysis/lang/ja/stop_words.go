package ja

import (
	_ "embed"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/token/stop"
	"github.com/blevesearch/bleve/v2/registry"
)

func init() {
	registry.RegisterTokenMap(StopWordsName, StopWordsTokenMapConstructor)
	registry.RegisterTokenFilter(StopWordsName, StopWordsTokenFilterConstructor)
}

// StopWordsName is the name of the stop words filter.
const StopWordsName = "stop_words_ja"

// StopWordsBytes is a stop word list.
// see. https://github.com/apache/lucene-solr/blob/master/lucene/analysis/kuromoji/src/resources/org/apache/lucene/analysis/ja/stopwords.txt
//
//go:embed assets/stop_words.txt
var StopWordsBytes []byte

// StopWordsTokenMapConstructor returns a token map for stop words.
func StopWordsTokenMapConstructor(_ map[string]any, _ *registry.Cache) (analysis.TokenMap, error) {
	rv := analysis.NewTokenMap()
	err := rv.LoadBytes(StopWordsBytes)
	return rv, err
}

// StopWordsTokenFilterConstructor returns a token filter for stop words.
func StopWordsTokenFilterConstructor(_ map[string]any, cache *registry.Cache) (analysis.TokenFilter, error) { //nolint:ireturn
	tm, err := cache.TokenMapNamed(StopWordsName)
	if err != nil {
		return nil, err
	}
	return stop.NewStopTokensFilter(tm), nil
}
