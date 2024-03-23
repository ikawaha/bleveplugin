package ja

import (
	"reflect"
	"testing"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

func TestStopTags(t *testing.T) {
	cache := registry.NewCache()
	tags, err := cache.TokenMapNamed(StopTagsName)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(tags), 27; got != want {
		t.Errorf("stop tags size: got %d, want %d", got, want)
	}
	want := analysis.TokenMap{
		"接続詞": true,
		"助詞-副助詞／並立助詞／終助詞": true,
		"助詞-特殊":     true,
		"助詞-間投助詞":   true,
		"助詞-並立助詞":   true,
		"助詞-連体化":    true,
		"助詞-副助詞":    true,
		"助詞-接続助詞":   true,
		"助詞-係助詞":    true,
		"助詞-終助詞":    true,
		"助詞-副詞化":    true,
		"記号-括弧開":    true,
		"その他-間投":    true,
		"助詞-格助詞":    true,
		"助詞-格助詞-連語": true,
		"非言語音":      true,
		"記号-句点":     true,
		"記号-括弧閉":    true,
		"フィラー":      true,
		"助詞-格助詞-引用": true,
		"記号":        true,
		"記号-一般":     true,
		"記号-空白":     true,
		"助詞-格助詞-一般": true,
		"助動詞":       true,
		"記号-読点":     true,
		"助詞":        true,
	}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("got %+v, want %+v", tags, want)
	}
}
