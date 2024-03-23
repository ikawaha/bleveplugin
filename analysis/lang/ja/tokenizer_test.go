package ja

import (
	"reflect"
	"testing"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
)

func TestJapaneseTokenizer_Tokenize(t *testing.T) {
	cache := registry.NewCache()
	stopTags, err := cache.TokenMapNamed(StopTagsName)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name  string
		dict  *dict.Dict
		opts  []TokenizerOption
		input []byte
		want  analysis.TokenStream
	}{
		{
			name:  "文分割なし",
			input: []byte("私は鰻"),
			dict:  ipa.Dict(),
			want: analysis.TokenStream{
				{
					Start:    0,
					End:      3,
					Term:     []byte("私"),
					Position: 1,
					Type:     analysis.Ideographic,
				},
				{
					Start:    3,
					End:      6,
					Term:     []byte("は"),
					Position: 2,
					Type:     analysis.Ideographic,
				},
				{
					Start:    6,
					End:      9,
					Term:     []byte("鰻"),
					Position: 3,
					Type:     analysis.Ideographic,
				},
			},
		},
		{
			name:  "文分割あり",
			dict:  ipa.Dict(),
			input: []byte("私は鰻。ねこはいます。"),
			want: analysis.TokenStream{
				{
					Start:    0,
					End:      3,
					Term:     []byte("私"),
					Position: 1,
					Type:     analysis.Ideographic,
				},
				{
					Start:    3,
					End:      6,
					Term:     []byte("は"),
					Position: 2,
					Type:     analysis.Ideographic,
				},
				{
					Start:    6,
					End:      9,
					Term:     []byte("鰻"),
					Position: 3,
					Type:     analysis.Ideographic,
				},
				{
					Start:    9,
					End:      12,
					Term:     []byte("。"),
					Position: 4,
					Type:     analysis.Ideographic,
				},
				{
					Start:    12,
					End:      18,
					Term:     []byte("ねこ"),
					Position: 5,
					Type:     analysis.Ideographic,
				},
				{
					Start:    18,
					End:      21,
					Term:     []byte("は"),
					Position: 6,
					Type:     analysis.Ideographic,
				},
				{
					Start:    21,
					End:      24,
					Term:     []byte("い"),
					Position: 7,
					Type:     analysis.Ideographic,
				},
				{
					Start:    24,
					End:      30,
					Term:     []byte("ます"),
					Position: 8,
					Type:     analysis.Ideographic,
				},
				{
					Start:    30,
					End:      33,
					Term:     []byte("。"),
					Position: 9,
					Type:     analysis.Ideographic,
				},
			},
		},
		{
			name:  "文分割・フィルターあり",
			dict:  ipa.Dict(),
			input: []byte("私は鰻。ねこはいます。"),
			opts: []TokenizerOption{
				StopTagsFilter(stopTags),
				BaseFormFilter(DefaultInflected),
			},
			want: analysis.TokenStream{
				{
					Start:    0,
					End:      3,
					Term:     []byte("私"),
					Position: 1,
					Type:     analysis.Ideographic,
				},
				{
					Start:    6,
					End:      9,
					Term:     []byte("鰻"),
					Position: 3,
					Type:     analysis.Ideographic,
				},
				{
					Start:    12,
					End:      18,
					Term:     []byte("ねこ"),
					Position: 5,
					Type:     analysis.Ideographic,
				},
				{
					Start:    21,
					End:      24,
					Term:     []byte("いる"),
					Position: 7,
					Type:     analysis.Ideographic,
				},
			},
		},
		{
			// Start: 0  End: 3  Position: 1  Token: 私  Type: 1
			// Start: 3  End: 6  Position: 1  Token: は  Type: 1 <drop>
			// Start: 6  End: 9  Position: 1  Token: 鰻  Type: 1
			// Start: 9  End: 12  Position: 1  Token: 。  Type: 1 <drop>
			// --- 文区切り
			// Start: 12  End: 15  Position: 1  Token: は  Type: 1 <drop>
			// Start: 15  End: 18  Position: 1  Token: 。  Type: 1 <drop>
			// --- 文区切り
			// Start: 18  End: 24  Position: 1  Token: ねこ  Type: 1
			// Start: 24  End: 27  Position: 1  Token: は  Type: 1 <drop>
			// Start: 27  End: 30  Position: 1  Token: い  Type: 1
			// Start: 30  End: 36  Position: 1  Token: ます  Type: 1 <drop>
			// Start: 36  End: 39  Position: 1  Token: 。  Type: 1 <drop>
			name:  "文ごとDropされるとき",
			dict:  ipa.Dict(),
			input: []byte("私は鰻。は。ねこはいます。"),
			opts: []TokenizerOption{
				StopTagsFilter(stopTags),
				BaseFormFilter(DefaultInflected),
			},
			want: analysis.TokenStream{
				{
					Start:    0,
					End:      3,
					Term:     []byte("私"),
					Position: 1,
					Type:     analysis.Ideographic,
				},
				{
					Start:    6,
					End:      9,
					Term:     []byte("鰻"),
					Position: 3,
					Type:     analysis.Ideographic,
				},
				{
					Start:    18,
					End:      24,
					Term:     []byte("ねこ"),
					Position: 7,
					Type:     analysis.Ideographic,
				},
				{
					Start:    27,
					End:      30,
					Term:     []byte("いる"),
					Position: 9,
					Type:     analysis.Ideographic,
				},
			},
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			tz := NewJapaneseTokenizer(v.dict, v.opts...)
			if got := tz.Tokenize(v.input); !reflect.DeepEqual(got, v.want) {
				t.Errorf("got %+v, want %+v", got, v.want)
			}
		})
	}
}

func TestTokenizerDefaultStopTagFilterMatch(t *testing.T) {
	cache := registry.NewCache()
	stopTags, err := cache.TokenMapNamed(StopTagsName)
	if err != nil {
		t.Fatal(err)
	}
	var tnz JapaneseTokenizer
	opt := StopTagsFilter(stopTags)
	opt(&tnz)

	tests := []struct {
		name string
		pos  []string
		want bool
	}{
		{
			name: "match:接続詞",
			pos:  []string{"接続詞", "*", "*", "*"},
			want: true,
		},
		{
			name: "match:助詞-副助詞／並立助詞／終助詞",
			pos:  []string{"助詞", "副助詞／並立助詞／終助詞", "*", "*"},
			want: true,
		},
		{
			name: "match:助詞-格助詞-引用",
			pos:  []string{"助詞", "格助詞", "引用", "*"},
			want: true,
		},
		{
			name: "not match:名詞-接続詞的",
			pos:  []string{"名詞", "接続詞的", "*", "*"},
			want: false,
		},
		{
			name: "not match:名詞-代名詞-一般,*",
			pos:  []string{"名詞", "代名詞", "一般", "*"},
			want: false,
		},
		{
			name: "invalid pos hierarchy",
			pos:  []string{"接続詞", "*", "*"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tnz.stopTagFilter.Match(tt.pos); got != tt.want {
				t.Errorf("stopTagsFilter.Match(%+v) = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}
