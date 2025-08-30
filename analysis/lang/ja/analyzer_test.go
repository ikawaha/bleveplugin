package ja

import (
	"reflect"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
)

func TestCustomAnalyzer(t *testing.T) {
	tests := []struct {
		title  string
		input  []byte
		output analysis.TokenStream
	}{
		{
			title: "tokenize",
			input: []byte("関西国際空港"),
			output: analysis.TokenStream{
				{Term: []byte("関西"), Position: 1, Start: 0, End: 6, Type: analysis.Ideographic},
				{Term: []byte("国際"), Position: 2, Start: 6, End: 12, Type: analysis.Ideographic},
				{Term: []byte("空港"), Position: 3, Start: 12, End: 18, Type: analysis.Ideographic},
			},
		},
		{
			title: "filtered results: stop tags filter & stop word filter",
			input: []byte("これらは私の猫"),
			output: analysis.TokenStream{
				{Term: []byte("私"), Position: 3, Start: 12, End: 15, Type: analysis.Ideographic},
				{Term: []byte("猫"), Position: 5, Start: 18, End: 21, Type: analysis.Ideographic},
			},
		},
		{
			title: "filtered results:赤い蝋燭と人魚",
			input: []byte("人魚は、南の方の海にばかり棲んでいるのではありません。"),
			// 人魚は、南の方の海にばかり棲んでいるのではありません。
			// 人魚	名詞,一般,*,*,*,*,人魚,ニンギョ,ニンギョ
			// は	助詞,係助詞,*,*,*,*,は,ハ,ワ
			// 、	記号,読点,*,*,*,*,、,、,、
			// 南	名詞,一般,*,*,*,*,南,ミナミ,ミナミ
			// の	助詞,連体化,*,*,*,*,の,ノ,ノ
			// 方	名詞,非自立,一般,*,*,*,方,ホウ,ホー
			// の	助詞,連体化,*,*,*,*,の,ノ,ノ
			// 海	名詞,一般,*,*,*,*,海,ウミ,ウミ
			// に	助詞,格助詞,一般,*,*,*,に,ニ,ニ
			// ばかり	助詞,副助詞,*,*,*,*,ばかり,バカリ,バカリ
			// 棲ん	動詞,自立,*,*,五段・マ行,連用タ接続,棲む,スン,スン
			// で	助詞,接続助詞,*,*,*,*,で,デ,デ
			// いる	動詞,非自立,*,*,一段,基本形,いる,イル,イル
			// の	名詞,非自立,一般,*,*,*,の,ノ,ノ
			// で	助動詞,*,*,*,特殊・ダ,連用形,だ,デ,デ
			// は	助詞,係助詞,*,*,*,*,は,ハ,ワ
			// あり	動詞,自立,*,*,五段・ラ行,連用形,ある,アリ,アリ
			// ませ	助動詞,*,*,*,特殊・マス,未然形,ます,マセ,マセ
			// ん	助動詞,*,*,*,不変化型,基本形,ん,ン,ン
			// 。	記号,句点,*,*,*,*,。,。,。
			output: analysis.TokenStream{
				{Term: []byte("人魚"), Position: 1, Start: 0, End: 6, Type: analysis.Ideographic},
				{Term: []byte("南"), Position: 4, Start: 12, End: 15, Type: analysis.Ideographic},
				{Term: []byte("方"), Position: 6, Start: 18, End: 21, Type: analysis.Ideographic},
				{Term: []byte("海"), Position: 8, Start: 24, End: 27, Type: analysis.Ideographic},
				{Term: []byte("棲む"), Position: 11, Start: 39, End: 45, Type: analysis.Ideographic}, // Note: base form
			},
		},
	}
	im := bleve.NewIndexMapping()
	if err := im.AddCustomTokenizer("ja", map[string]any{
		"type":      Name,
		"dict":      "ipa",
		"base_form": true,
		"stop_tags": true,
	}); err != nil {
		t.Fatal(err)
	}
	if err := im.AddCustomAnalyzer("ja", map[string]any{
		"type":      custom.Name,
		"tokenizer": "ja",
		"token_filters": []string{
			StopWordsName,
			lowercase.Name,
		},
	}); err != nil {
		t.Fatal(err)
	}
	analyzer := im.AnalyzerNamed("ja")
	if analyzer == nil {
		t.Fatal("analyzer is nil")
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			actual := analyzer.Analyze(test.input)
			if !reflect.DeepEqual(actual, test.output) {
				t.Errorf("want %+v, got %+v", test.output, actual)
			}
		})
	}
}

func BenchmarkJapaneseAnalyzer(b *testing.B) {
	im := bleve.NewIndexMapping()
	if err := im.AddCustomTokenizer("ja", map[string]any{
		"type":      Name,
		"dict":      "uni",
		"base_form": true,
		"stop_tags": true,
	}); err != nil {
		b.Fatal(err)
	}
	if err := im.AddCustomAnalyzer("ja", map[string]any{
		"type":      custom.Name,
		"tokenizer": "ja",
		"token_filters": []string{
			StopWordsName,
			lowercase.Name,
		},
	}); err != nil {
		b.Fatal(err)
	}
	analyzer := im.AnalyzerNamed("ja")
	if analyzer == nil {
		b.Fatal("analyzer is nil")
	}
	sen := []byte("人魚は、南の方の海にばかり棲んでいるのではありません。")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.Analyze(sen)
	}
}
