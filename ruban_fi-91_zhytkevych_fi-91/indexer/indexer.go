package indexer

import (
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

type Indexer interface {
	IndexDocument(docId int, document []byte)
	GetDocsByKeyword(word string) map[int][]int
	GetDocsByPrefix(word string) map[int][]int
	GetDocsByKeywords(word1 string, word2 string, dist int) map[int][]int
}

type IndexerBtree struct {
	btree *storage.Btree
}

func NewIndexer(storagePath string) *IndexerBtree {
	btree := storage.NewBtree(storagePath)
	return &IndexerBtree{
		btree: btree,
	}
}

func regexSubstrings(str string, reg string) []string {
	return regexp.MustCompile(reg).FindAllString(str, -1)
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

func (i *IndexerBtree) IndexDocument(docId int, document []byte) {
	words := regexSubstrings(string(document), `[a-zA-Z0-9_]+`)
	i.btree.AddIndexes(makeInvertedIndexes(words, docId))
}

func mapKeys(m map[int][]int) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func (i *IndexerBtree) GetDocsByKeyword(word string) ([]int, error) {
	word = strings.ToLower(word)
	sheetEl, err := i.btree.Find(word)
	log.Println(sheetEl.Data)
	if err != nil {
		return []int{}, err
	}
	return mapKeys(sheetEl.Data), nil
}

func (i *IndexerBtree) GetDocsByKeywords(word1 string, word2 string, dist uint) ([]int, error) {
	e1, err := i.btree.Find(word1)
	e2, err := i.btree.Find(word2)
	var docIds []int
	if err != nil {
		return docIds, err
	}
	for id, positions := range e1.Data {
		for _, pos := range positions {
			k := sort.Search(len(e2.Data[id]), func(j int) bool { return e2.Data[id][j] == pos+1+int(dist) })
			log.Println("k", k, pos+1+int(dist), e2.Data[id])
			if k < len(e2.Data[id]) && e2.Data[id][k] == pos+1+int(dist) {
				docIds = append(docIds, id)
				break
			}
		}
	}
	return docIds, nil
}
