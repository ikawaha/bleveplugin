bleve plugin
===
[![GoDoc](https://godoc.org/github.com/ikawaha/bleveplugin/analysis/lang/ja?status.svg)](https://godoc.org/github.com/ikawaha/bleveplugin)

Japanese language analysis plugins for the [bleve v2](https://github.com/blevesearch/bleve) indexing/search library.


# Usage example

blog: comming soon

see. example/analyzer/main.go
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/ikawaha/bleveplugin/analysis/lang/ja"
)

func main() {
	if err := run(os.Args); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(_ []string) error {
	// document mapping
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name
	jaTextFieldMapping := bleve.NewTextFieldMapping()
	jaTextFieldMapping.Analyzer = "ja"
	dm := bleve.NewDocumentMapping()
	dm.AddFieldMappingsAt("type", keywordFieldMapping)
	dm.AddFieldMappingsAt("id", jaTextFieldMapping)
	dm.AddFieldMappingsAt("author", jaTextFieldMapping)
	dm.AddFieldMappingsAt("text", jaTextFieldMapping)

	// index mapping
	im := bleve.NewIndexMapping()
	im.TypeField = "type"
	im.AddDocumentMapping("book", dm)
	if err := im.AddCustomTokenizer("ja_tokenizer", map[string]any{
		"type":      ja.Name,
		"dict":      ja.DictIPA,
		"base_form": true,
		"stop_tags": true,
	}); err != nil {
		return fmt.Errorf("failed to create ja tokenizer: %w", err)
	}
	if err := im.AddCustomAnalyzer("ja", map[string]any{
		"type":      custom.Name,
		"tokenizer": "ja_tokenizer",
		"token_filters": []string{
			ja.StopWordsName,
			lowercase.Name,
		},
	}); err != nil {
		return fmt.Errorf("failed to create ja analyzer: %w", err)
	}

	// index
	index, err := bleve.NewMemOnly(im)
	defer index.Close()

	// documents
	docs := []string{
		`{"type": "book", "id": "1:赤い蝋燭と人魚", "author": "小川未明", "text": "人魚は南の方の海にばかり棲んでいるのではありません"}`,
		`{"type": "book", "id": "2:吾輩は猫である", "author": "夏目漱石", "text":   "吾輩は猫である。名前はまだない"}`,
		`{"type": "book", "id": "3:狐と踊れ", "author": "神林長平", "text": "踊っているのでなければ踊らされているのだろうさ"}`,
		`{"type": "book", "id": "4:ダンスダンスダンス", "author": "村上春樹", "text": "音楽の鳴っている間はとにかく踊り続けるんだ。おいらの言っていることはわかるかい？"}`,
	}
	// indexing
	for _, doc := range docs {
		var data map[string]any
		if err := json.Unmarshal([]byte(doc), &data); err != nil {
			log.Printf("SKIP: failed to unmarshal doc: %v, %s", err, doc)
			continue
		}
		if err := index.Index(data["id"].(string), data); err != nil {
			return fmt.Errorf("error indexing document: %w", err)
		}
		// printing ...
		fmt.Printf("indexed document with id:%s, %s\n", data["id"], doc)
	}
	dc, err := index.DocCount()
	if err != nil {
		return fmt.Errorf("failed to count documents: %w", err)
	}
	fmt.Printf("doc count: %d\n --------\n", dc)

	// query
	q := "踊る人形"
	query := bleve.NewMatchQuery(q)
	query.SetField("text")
	req := bleve.NewSearchRequest(query)
	fmt.Printf("query: search field %q, value %q\n", query.Field(), q)

	// search
	result, err := index.Search(req)
	if err != nil {
		log.Fatalf("error executing search: %v", err)
	}
	// search result
	for _, v := range result.Hits {
		fmt.Println(v.ID)
	}
	return nil
}
```

OUTPUT:
```
indexed document with id:1:赤い蝋燭と人魚, {"type": "book", "id": "1:赤い蝋燭と人魚", "author": "小川未明", "text": "人魚は南の方の海にばかり棲んでいるのではありません"}
indexed document with id:2:吾輩は猫である, {"type": "book", "id": "2:吾輩は猫である", "author": "夏目漱石", "text":   "吾輩は猫である。名前はまだない"}
indexed document with id:3:狐と踊れ, {"type": "book", "id": "3:狐と踊れ", "author": "神林長平", "text": "踊っているのでなければ踊らされているのだろうさ"}
indexed document with id:4:ダンスダンスダンス, {"type": "book", "id": "4:ダンスダンスダンス", "author": "村上春樹", "text": "音楽の鳴っている間はとにかく踊り続けるんだ。おいらの言っていることはわかるかい？"}
doc count: 4
 --------
query: search field "text", value "踊る人形"
3:狐と踊れ
4:ダンスダンスダンス
```
