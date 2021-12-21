package indexer

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

type Indexer interface {
	IndexDocument(docId uint64, collectionName string, document []byte) error
	BtreesString() string
	GetDocsByKeyword(collectionName string, word string) ([]uint64, error)
	GetDocsByPrefix(collectionName string, word string) ([]uint64, error)
	GetDocsByKeywords(collectionName string, word1 string, word2 string, dist uint) ([]uint64, error)
}

type BtreeEntry struct {
	Name  string
	Btree *storage.Btree
}

type IndexerBtree struct {
	path   string
	btrees []BtreeEntry
}

func NewIndexerBtree(storagePath string) *IndexerBtree {
	return &IndexerBtree{
		path:   storagePath,
		btrees: readPathBtrees(storagePath),
	}
}

func readPathBtrees(path string) (btrees []BtreeEntry) {
	entries, _ := os.ReadDir(path)
	btrees = make([]BtreeEntry, 0)

	var entryPath string
	for _, entry := range entries {
		if entry.IsDir() {
			entryPath = filepath.Join(path, entry.Name())
			btrees = append(btrees, BtreeEntry{
				Name:  entry.Name(),
				Btree: storage.NewBtree(entryPath),
			})
		}
	}
	return
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

func (i *IndexerBtree) getBtreeByCollection(name string) *storage.Btree {
	for _, bt := range i.btrees {
		if bt.Name == name {
			return bt.Btree
		}
	}
	return nil
}

func (i *IndexerBtree) IndexDocument(docId uint64, collectionName string, document []byte) error {
	words := regexSubstrings(string(document), `[a-zA-Z0-9_]+`)
	btree := i.getBtreeByCollection(collectionName)
	if btree == nil {
		// create new btree
		btreePath := filepath.Join(i.path, collectionName)
		if err := os.Mkdir(btreePath, os.ModeDir|os.ModePerm); err != nil {
			panic(err)
		}
		btree = storage.NewBtree(btreePath)
		i.btrees = append(i.btrees, BtreeEntry{
			Name:  collectionName,
			Btree: btree,
		})
	}
	return btree.AddIndexes(makeInvertedIndexes(words, docId))
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

func (i *IndexerBtree) GetDocsByKeyword(collectionName string, word string) ([]uint64, error) {
	btree := i.getBtreeByCollection(collectionName)
	if btree == nil {
		return []uint64{}, errors.New("No such collection")
	}
	word = strings.ToLower(word)
	sheetEl, err := btree.Find(word)
	if err != nil {
		return []uint64{}, err
	}
	return mapKeys(sheetEl.Data), nil
}

func makePositionsHash(data map[uint64][]int) map[uint64]map[int]bool {
	hashmap := make(map[uint64]map[int]bool, 0)
	for id, positions := range data {
		hashmap[id] = make(map[int]bool, 0)
		for _, pos := range positions {
			hashmap[id][pos] = true
		}
	}
	return hashmap
}

func (i *IndexerBtree) GetDocsByKeywords(collectionName string, word1 string, word2 string, dist uint) ([]uint64, error) {
	btree := i.getBtreeByCollection(collectionName)
	if btree == nil {
		return []uint64{}, errors.New("No such collection")
	}
	e1, err := btree.Find(word1)
	e2, err := btree.Find(word2)
	docIds := NewSetUint()
	if err != nil {
		return docIds.ToArray(), err
	}
	e1_hash := makePositionsHash(e1.Data)
	for e2_id, positions := range e2.Data {
		for _, pos := range positions {
			if e1_hash[e2_id][pos+int(dist)] {
				docIds.Add(e2_id)
			}
			if e1_hash[e2_id][pos-int(dist)] {
				docIds.Add(e2_id)
			}
		}
	}
	return docIds.ToArray(), nil
}

func (i *IndexerBtree) GetDocsByPrefix(collectionName string, prefix string) ([]uint64, error) {
	btree := i.getBtreeByCollection(collectionName)
	if btree == nil {
		return []uint64{}, errors.New("No such collection")
	}
	prefix = strings.ToLower(prefix)
	shEls, err := btree.FindByPrefix(prefix)
	if err != nil {
		return []uint64{}, err
	}
	docIds := GetDocIds(shEls)
	return docIds, nil
}

func (i *IndexerBtree) BtreesString() (out string) {
	for _, entry := range i.btrees {
		out += entry.Name + ": "
		out += entry.Btree.String()
		out += "\n"
	}
	return
}
