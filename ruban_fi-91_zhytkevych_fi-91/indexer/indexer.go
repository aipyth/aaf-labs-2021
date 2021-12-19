package indexer

import (
	"regexp"
	"strings"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

type Indexer interface {
	IndexDocument(docId uint64, document []byte) error
	IndexString() string
	GetDocsByKeyword(word string) ([]uint64, error)
	GetDocsByPrefix(word string) ([]uint64, error)
	GetDocsByKeywords(word1 string, word2 string, dist uint) ([]uint64, error)
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

func regexSubstrings(str string, reg string) []string {
	return regexp.MustCompile(reg).FindAllString(str, -1)
}

func makeInvertedIndexes(words []string, docId uint64) map[string]map[uint64][]int {
	m := make(map[string]map[uint64][]int)
	for i, w := range words {
		w = strings.ToLower(w)
		if len(m[w]) == 0 {
			m[w] = map[uint64][]int{docId: []int{i}}
			continue
		}
		m[w][docId] = append(m[w][docId], int(i))
	}
	return m
}

func (i *IndexerBtree) IndexDocument(docId uint64, document []byte) error {
	words := regexSubstrings(string(document), `[a-zA-Z0-9_]+`)
	return i.btree.AddIndexes(makeInvertedIndexes(words, docId))
}

func mapKeys(m map[uint64][]int) []uint64 {
	keys := make([]uint64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func Includes(slice []uint64, el uint64) bool {
	for _, x := range slice {
		if x == el {
			return true
		}
	}
	return false
}

func GetDocIds(shEls []*storage.SheetElement) []uint64 {
	docIds := make([]uint64, 0)
	for _, shEl := range shEls {
		for docId := range shEl.Data {
			if !Includes(docIds, docId) {
				docIds = append(docIds, docId)
			}
		}
	}
	return docIds
}

func (i *IndexerBtree) GetDocsByKeyword(word string) ([]uint64, error) {
	word = strings.ToLower(word)
	sheetEl, err := i.btree.Find(word)
	if err != nil {
		return []uint64{}, err
	}
	return mapKeys(sheetEl.Data), nil
}

func makePositionsHash(data map[uint64][]int) map[int]uint64 {
	hashmap := make(map[int]uint64, 0)
	for id, positions := range data {
		for _, pos := range positions {
			hashmap[pos] = id
		}
	}
	return hashmap
}

func (i *IndexerBtree) GetDocsByKeywords(word1 string, word2 string, dist uint) ([]uint64, error) {
	e1, err := i.btree.Find(word1)
	e2, err := i.btree.Find(word2)
	var docIds []uint64
	if err != nil {
		return docIds, err
	}
	e1_hash := makePositionsHash(e1.Data)
	for e2_id, positions := range e2.Data {
		for _, pos := range positions {
			if e1_id, ok := e1_hash[pos+int(dist)]; ok && e1_id == e2_id {
				docIds = append(docIds, e1_id)
			}
			if e1_id, ok := e1_hash[pos-int(dist)]; ok && e1_id == e2_id {
				docIds = append(docIds, e1_id)
			}
		}
	}
	return docIds, nil
}

func (i *IndexerBtree) GetDocsByPrefix(prefix string) ([]uint64, error) {
	prefix = strings.ToLower(prefix)
	shEls, err := i.btree.FindByPrefix(prefix)
	if err != nil {
		return []uint64{}, err
	}
	docIds := GetDocIds(shEls)
	return docIds, nil
}

func (i *IndexerBtree) IndexString() string {
	return i.btree.String()
}
