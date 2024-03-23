package ja

import (
	"reflect"
	"testing"

	"github.com/blevesearch/bleve/v2/analysis"
	"github.com/blevesearch/bleve/v2/registry"
)

func TestDefaultStopWords(t *testing.T) {
	cache := registry.NewCache()
	words, err := cache.TokenMapNamed(StopWordsName)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(words), 109; got != want {
		t.Errorf("stop words size: got %d, want %d", got, want)
	}
	want := analysis.TokenMap{
		"い": true, "う": true, "お": true, "か": true,
		"が": true, "き": true, "さ": true, "し": true,
		"ず": true, "せ": true, "た": true, "だ": true,
		"つ": true, "て": true, "で": true, "と": true,
		"な": true, "に": true, "の": true, "は": true,
		"ば": true, "へ": true, "も": true, "や": true,
		"ら": true, "れ": true, "を": true, "ん": true,
		"あっ": true, "あり": true, "ある": true, "いう": true,
		"いる": true, "うち": true, "おり": true, "から": true,
		"ここ": true, "こと": true, "この": true, "これ": true,
		"する": true, "せる": true, "その": true, "それ": true,
		"たち": true, "ため": true, "たり": true, "だっ": true,
		"でき": true, "です": true, "では": true, "でも": true,
		"とき": true, "とも": true, "ない": true, "なお": true,
		"なく": true, "なっ": true, "など": true, "なら": true,
		"なり": true, "なる": true, "にて": true, "ので": true,
		"のみ": true, "ほか": true, "ほど": true, "ます": true,
		"また": true, "まで": true, "もの": true, "よう": true,
		"より": true, "られ": true, "れる": true, "及び": true,
		"特に": true, "および": true, "かつて": true, "これら": true,
		"さらに": true, "しかし": true, "そして": true, "その他": true,
		"その後": true, "ただし": true, "できる": true, "という": true,
		"ところ": true, "として": true, "と共に": true, "なかっ": true,
		"ながら": true, "により": true, "による": true, "または": true,
		"ものの": true, "られる": true, "それぞれ": true, "といった": true,
		"とともに": true, "において": true, "における": true, "について": true,
		"によって": true, "に対して": true, "に対する": true, "に関する": true,
		"ほとんど": true,
	}
	if !reflect.DeepEqual(words, want) {
		t.Errorf("got %+v, want %+v", words, want)
	}
}
