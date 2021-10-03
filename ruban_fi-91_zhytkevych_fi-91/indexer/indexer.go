package indexer

import (
	"log"
	"regexp"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

type Indexer interface {
	IndexDocument(collId int, document []byte)
	GetDocsByKeyword(word string) map[int][]int
	GetDocsByPrefix(word string) map[int][]int
	GetDocsByKeywords(word1 string, word2 string, dist int) map[int][]int
}

type IndexerBtree struct {
	btree *storage.Btree
}

func (i *IndexerBtree) IndexDocument(docId int, document []byte) {
	words := regexp.Split(string(document), "[^a-zA-Z0-9_]")
	inverted_indexes := make(map[string]map[int][]int)
	for i, w := range words {
		if inverted_indexes[w] == nil {
			value := make(map[int][]int)
			value[docId] = []int{i}
			continue
		}
		inverted_indexes[w][docId] = append(inverted_indexes[w][docId], storage.PosIdx(i))
	}
	i.btree.AddIndexes(inverted_indexes)
}

func (i *IndexerBtree) GetDocsByKeyword(word string) map[int][]int {
	sheetEl, err := i.btree.Find(word)
	if err != nil {
		log.Println(err)
	}
	return sheetEl.Data
}
