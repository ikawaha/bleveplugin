package ja

import (
	_ "embed"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

func init() {
	registry.RegisterTokenMap(StopTagsName, StopTagsTokenMapConstructor)
}

const StopTagsName = "stop_tags_ja"

// StopTagsBytes is a stop tag list.
// see. https://github.com/apache/lucene-solr/blob/master/lucene/analysis/kuromoji/src/resources/org/apache/lucene/analysis/ja/stoptags.txt
//
//go:embed assets/stop_tags.txt
var StopTagsBytes []byte

// StopTagsTokenMapConstructor returns a token map for stop tags (for IPA dict).
func StopTagsTokenMapConstructor(_ map[string]any, _ *registry.Cache) (analysis.TokenMap, error) {
	rv := analysis.NewTokenMap()
	err := rv.LoadBytes(StopTagsBytes)
	return rv, err
}
