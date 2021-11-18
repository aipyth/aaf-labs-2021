package storage

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func encode(t *testing.T, i interface{}) bytes.Buffer {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(i)
	if err != nil {
		t.Error(err.Error())
	}
	return buf
}

func TestSheetAddElement(t *testing.T) {
	// t.Run("empty sheet encode", func(t *testing.T) {
	// 	sh := NewSheet()
	// 	bts := sh.Encode()
	// 	encoded := encode(t, &Sheet{
	// 		Keys: make([]*SheetElement, 0),
	// 	})
	// 	if bytes.Compare(bts.Bytes(), encoded.Bytes()) != 0 {
	// 		t.Errorf("Encoded data does not match\n")
	// 	}
	// })

	// t.Run("sheet add elements", func(t *testing.T) {
	// 	sheet := NewSheet()
	// 	sheet.Add("vanya", nil)
	// 	sheet.Add("go", nil)
	// 	sheet.Add("dota", nil)
	// 	appendToSheet(sheet, &SheetElement{
	// 		Key:  "chert",
	// 		Data: nil,
	// 	})
	// 	appendToSheet(sheet, &SheetElement{
	// 		Key:  "katat",
	// 		Data: nil,
	// 	})

	// 	err := appendToSheet(sheet, &SheetElement{
	// 		Key:  "posle",
	// 		Data: nil,
	// 	})

	// 	err = appendToSheet(sheet.Children[1], &SheetElement{
	// 		Key:  "luche",
	// 		Data: nil,
	// 	})

	// 	err = appendToSheet(sheet.Children[1], &SheetElement{
	// 		Key:  "spat",
	// 		Data: nil,
	// 	})

	// 	err = appendToSheet(sheet.Children[1], &SheetElement{
	// 		Key:  "poyti",
	// 		Data: nil
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// })

	// t.Run("split sheets by two on half", func(t *testing.T) {
	// 	sheet := NewSheet()
	// 	sheet.Add("ya", nil)
	// 	sheet.Add("hochu", nil)
	// 	sheet.Add("pivo", nil)
	// 	sheet.Add("vodku", nil)
	// 	sheet.Add("mb", nil)
	// 	left, right, elem := splitSheetsByHalf(sheet)
	// 	log.Println(left, right, elem.Key)
	// })

	t.Run("Test tree adding elements and find", func(t *testing.T) {
		log.Println("starting")
		btree := NewBtree("./btree-storage/")
		btree.AddIndex("a", nil)
		btree.AddIndex("b", nil)
		btree.AddIndex("c", nil)
		btree.AddIndex("d", nil)
		btree.AddIndex("e", nil)
		btree.AddIndex("f", nil)
		btree.AddIndex("g", nil)
		btree.AddIndex("h", nil)
		btree.AddIndex("i", nil)
		btree.AddIndex("j", nil)
		btree.AddIndex("k", nil)
		btree.AddIndex("l", nil)
		btree.AddIndex("n", nil)
		btree.AddIndex("m", nil)
		btree.AddIndex("o", nil)
		btree.AddIndex("p", nil)
		btree.AddIndex("q", nil)
		btree.AddIndex("r", nil)
		btree.AddIndex("s", nil)
		btree.AddIndex("t", nil)
		btree.AddIndex("u", nil)
		btree.AddIndex("v", nil)
		btree.AddIndex("w", nil)
		btree.AddIndex("x", nil)
		btree.AddIndex("y", nil)
		btree.AddIndex("z", nil)
		btree.AddIndex("za", nil)
		btree.AddIndex("zc", nil)
		btree.AddIndex("zd", nil)
		btree.AddIndex("ze", nil)
		btree.AddIndex("bc", nil)
		btree.AddIndex("obladi", nil)
		log.Println(btree)
		dir, _ := ioutil.ReadDir("./btree-storage")
		for _, file := range dir {
			os.Remove(path.Join("./btree-storage", file.Name()))
		}
	})

	t.Run("Test element finding", func(t *testing.T) {
		btree := NewBtree("./btree-storage")
		m_a := make(map[uint64][]int)
		m_e := make(map[uint64][]int)
		m_g := make(map[uint64][]int)
		m_l := make(map[uint64][]int)
		m_a[1] = []int{0, 1}
		m_e[1] = []int{0, 1}
		m_g[1] = []int{0, 1}
		m_l[1] = []int{0, 1}
		btree.AddIndex("a", m_a)
		btree.AddIndex("b", nil)
		btree.AddIndex("c", nil)
		btree.AddIndex("d", nil)
		btree.AddIndex("e", m_e)
		btree.AddIndex("f", nil)
		btree.AddIndex("g", m_g)
		btree.AddIndex("h", nil)
		btree.AddIndex("i", nil)
		btree.AddIndex("j", nil)
		btree.AddIndex("k", nil)
		btree.AddIndex("l", m_l)
		btree.AddIndex("obladi", m_l)
		btree.AddIndex("oblada", m_l)
		btree.AddIndex("obitel", m_l)
		btree.AddIndex("obivatel", m_l)
		btree.AddIndex("obratimost", m_l)
		btree.AddIndex("obuv", m_l)
		btree.AddIndex("obrubok", m_l)
		btree.AddIndex("obichnost", m_l)
		log.Println(btree)
		log.Println(btree.Find("a"))
		log.Println(btree.Find("e"))
		log.Println(btree.Find("g"))
		log.Println(btree.Find("l"))
		log.Println(btree.Find("h"))
		found, err := btree.FindByPrefix("ob")
		log.Println(len(found), err)
		dir, _ := ioutil.ReadDir("./btree-storage")
		for _, file := range dir {
			os.Remove(path.Join("./btree-storage", file.Name()))
		}
	})

}
