package indexer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestIndexerIndexDoc(t *testing.T) {

	t.Run("Test regex split", func(t *testing.T) {
		test := "Some string   to @$$@test,for--+#idk212  whyy...''lol\" 	"
		log.Println(regexSubstrings(test, `[a-zA-Z0-9_]+`))
	})

	t.Run("Test create index", func(t *testing.T) {
		test := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
		indexes := makeInvertedIndexes(regexSubstrings(test, `[a-zA-Z0-9_]+`), 1)
		bts := new(bytes.Buffer)
		for w, m := range indexes {
			fmt.Fprintf(bts, "%s: {\n", w)
			for id, poss := range m {
				fmt.Fprintf(bts, "	%d: %+v,\n", id, poss)
			}
		}
		// log.Println(bts)
	})

	t.Run("Test adding some elements", func(t *testing.T) {
		indexer := NewIndexerBtree("../storage/btree-storage")
		document := []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`)
		indexer.IndexDocument(1, document)
		log.Println(indexer.btree)
	})

	t.Run("Test finding some elements", func(t *testing.T) {
		indexer := NewIndexerBtree("../storage/btree-storage")
		indexes, _ := indexer.GetDocsByKeyword("Lorem")
		log.Println(indexes)

	})

	t.Run("Test finding elements by 2 keywords", func(t *testing.T) {
		indexer := NewIndexerBtree("../storage/btree-storage")
		indexes, _ := indexer.GetDocsByKeywords("ipsum", "elit", 5)
		log.Println("by 2 kwrds", indexes)
	})

	t.Run("Test finding elements by prefix", func(t *testing.T) {
		indexer := NewIndexerBtree("../storage/btree-storage")
		indexes, _ := indexer.GetDocsByPrefix("pariat")
		log.Println(indexes)
		dir, _ := ioutil.ReadDir("../storage/btree-storage")
		for _, file := range dir {
			os.Remove(path.Join("../storage/btree-storage", file.Name()))
		}
	})
}
