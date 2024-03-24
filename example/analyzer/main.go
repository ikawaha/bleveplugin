package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
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

func run(args []string) error {
	im := bleve.NewIndexMapping()
	if err := im.AddCustomTokenizer("ja", map[string]any{
		"type":      ja.Name,
		"dict":      ja.DictIPA,
		"base_form": true,
		"stop_tags": true,
	}); err != nil {
		return fmt.Errorf("failed to create ja tokenizer: %w", err)
	}
	if err := im.AddCustomAnalyzer("ja", map[string]any{
		"type":      custom.Name,
		"tokenizer": "ja",
		"token_filters": []string{
			ja.StopWordsName,
			lowercase.Name,
		},
	}); err != nil {
		return fmt.Errorf("failed to create ja analyzer: %w", err)
	}

	// document mapping
	dm := bleve.NewDocumentMapping()
	im.AddDocumentMapping("book", dm)
	author := bleve.NewTextFieldMapping()
	author.Name = "author"
	author.Analyzer = "ja"
	dm.AddFieldMapping(author)
	body := bleve.NewTextFieldMapping()
	body.Name = "text"
	body.Analyzer = "ja"
	dm.AddFieldMapping(body)

	// indexer
	index, err := bleve.NewMemOnly(im)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// 対象ドキュメント
	docs := []string{
		`{"id": "1:赤い蝋燭と人魚", "author": "小川未明", "text": "人魚は南の方の海にばかり棲んでいるのではありません"}`,
		`{"id": "2:吾輩は猫である", "author": "夏目漱石", "text":   "吾輩は猫である。名前はまだない"}`,
		`{"id": "3:狐と踊れ", "author": "神林長平", "text": "踊っているのでなければ踊らされているのだろうさ"}`,
		`{"id": "4:ダンスダンスダンス", "author": "村上春樹", "text": "音楽の鳴っている間はとにかく踊り続けるんだ。おいらの言っていることはわかるかい？"}`,
	}
	// indexing
	for _, doc := range docs {
		var data map[string]string
		if err := json.Unmarshal([]byte(doc), &data); err != nil {
			log.Printf("SKIP: failed to unmarshal doc: %v, %s", err, doc)
			continue
		}
		id := data["id"]
		if err := index.Index(id, doc); err != nil {
			return fmt.Errorf("error indexing document: %w", err)
		}
		// 表示
		fmt.Printf("indexed document with id:%s\n", id)
		//for _, v := range doc.Fields {
		//	fmt.Printf("\t%s: %s\n", v.Name(), v.Value())
		//}
	}

	// クエリ
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
	// 検索結果
	for _, v := range result.Hits {
		fmt.Println(v.ID, v.Expl, v.Fields)
	}
	return nil
}
