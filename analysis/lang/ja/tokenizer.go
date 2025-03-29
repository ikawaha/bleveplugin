package ja

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome-dict/uni"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

const (
	Name    = "ja_kagome"
	DictIPA = "ipa"
	DictUni = "uni"
)

func init() {
	registry.RegisterTokenizer(Name, TokenizerConstructor)
}

// TokenizerOption represents an option of the japanese tokenizer.
type TokenizerOption func(t *JapaneseTokenizer)

const (
	posHierarchy      = 4
	defaultPOSFeature = "*"
)

// StopTagsFilter returns a stop tags filter option.
func StopTagsFilter(m analysis.TokenMap) TokenizerOption {
	ps := make([]filter.POS, 0, len(m))
	for k := range m {
		pos := strings.Split(k, "-")
		for i := len(pos); i < posHierarchy; i++ {
			pos = append(pos, defaultPOSFeature)
		}
		ps = append(ps, pos)
	}
	ft := filter.NewPOSFilter(ps...)
	return func(t *JapaneseTokenizer) {
		t.stopTagFilter = ft
	}
}

// BaseFormFilter returns an base form filter option.
func BaseFormFilter(m analysis.TokenMap) TokenizerOption {
	ps := make([]filter.POS, 0, len(m))
	for k := range m {
		pos := strings.Split(k, "-")
		ps = append(ps, pos)
	}
	ft := filter.NewPOSFilter(ps...)
	return func(t *JapaneseTokenizer) {
		t.baseFormFilter = ft
	}
}

// JapaneseTokenizer represents a Japanese tokenizer with filters.
type JapaneseTokenizer struct {
	*tokenizer.Tokenizer
	stopTagFilter  *filter.POSFilter
	baseFormFilter *filter.POSFilter
}

var splitter = filter.SentenceSplitter{
	Delim:               []rune{'。', '．', '！', '!', '？', '?'},
	Follower:            []rune{'.', '｣', '」', '』', ')', '）', '｝', '}', '〉', '》'},
	SkipWhiteSpace:      false,
	DoubleLineFeedSplit: true,
	MaxRuneLen:          128,
}

// Tokenize tokenizes the input and filters them.
func (t *JapaneseTokenizer) Tokenize(input []byte) analysis.TokenStream {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Split(splitter.ScanSentences)
	base := 0
	position := 1
	var ret analysis.TokenStream
	for scanner.Scan() {
		inp := scanner.Text()
		tokens := t.Analyze(inp, tokenizer.Search)
		tokenLen := len(tokens)
		if t.stopTagFilter != nil {
			t.stopTagFilter.Drop(&tokens)
		}
		for _, v := range tokens {
			start := base + v.Position
			end := base + v.Position + len(v.Surface)
			term := input[start:end]
			if t.baseFormFilter != nil {
				if pos := v.POS(); t.baseFormFilter.Match(pos) {
					if base, ok := v.BaseForm(); ok {
						term = []byte(base)
					}
				}
			}
			ret = append(ret, &analysis.Token{
				Start:    start,
				End:      end,
				Term:     term,
				Position: position + v.Index,
				Type:     analysis.Ideographic,
				KeyWord:  false,
			})
		}
		base += len(inp)
		position += tokenLen
	}
	return ret
}

// NewJapaneseTokenizer returns a Japanese tokenizer.
func NewJapaneseTokenizer(dict *dict.Dict, opts ...TokenizerOption) *JapaneseTokenizer {
	t, err := tokenizer.New(dict, tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}
	ret := &JapaneseTokenizer{
		Tokenizer: t,
	}
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

func TokenizerConstructor(config map[string]any, cache *registry.Cache) (analysis.Tokenizer, error) { //nolint:ireturn
	var d *dict.Dict
	kind, ok := config["dict"]
	if !ok {
		return nil, errors.New(`config requires dict, e.g. "ipa" or "uni"`)
	}
	switch kind {
	case DictIPA:
		d = ipa.Dict()
	case DictUni:
		d = uni.Dict()
	default:
		return nil, fmt.Errorf("unsupported dictionary: %s", kind)
	}
	var opts []TokenizerOption
	if ok, _ := config["stop_tags"].(bool); ok {
		stopTags, err := cache.TokenMapNamed(StopTagsName)
		if err != nil {
			return nil, err
		}
		opts = append(opts, StopTagsFilter(stopTags))
	}
	if ok, _ := config["base_form"].(bool); ok {
		opts = append(opts, BaseFormFilter(DefaultInflected))
	}
	return NewJapaneseTokenizer(d, opts...), nil
}
