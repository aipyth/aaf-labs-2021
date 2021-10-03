package indexer

import (
	"log"
	"regexp"
	"strings"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

type Indexer interface {
	IndexDocument(collId int, document []byte) error
	GetDocsByKeyword(word string) map[int][]int
	GetDocsByPrefix(word string) map[int][]int
	GetDocsByKeywords(word1 string, word2 string, dist int) map[int][]int
}

type IndexerBtree struct {
	btree *storage.Btree
}

func NewIndexerBtree(storagePath string) *IndexerBtree {
	btree := storage.NewBtree(storagePath)
	return &IndexerBtree{
		btree: btree,
	}
}

func regexSplit(str string, reg string) []string {
	return regexp.MustCompile(reg).Split(str, -1)
}

func makeInvertedIndexes(words []string, docId int) map[string]map[int][]int {
	m := make(map[string]map[int][]int)
	for i, w := range words {
		w = strings.ToLower(w)
		if len(m[w]) == 0 {
			m[w] = map[int][]int{docId: []int{i}}
			continue
		}
		m[w][docId] = append(m[w][docId], int(i))
	}
	return m
}

func (i *IndexerBtree) IndexDocument(docId int, document []byte) error {
	words := regexSplit(string(document), `[^a-zA-Z0-9_+]`)
	return i.btree.AddIndexes(makeInvertedIndexes(words, docId))
}

func (i *IndexerBtree) GetDocsByKeyword(word string) map[int][]int {
	word = strings.ToLower(word)
	sheetEl, err := i.btree.Find(word)
	if err != nil {
		log.Println(err)
	}
	return sheetEl.Data
}
