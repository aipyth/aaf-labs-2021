package storage

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
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
	// 		Data: nil,
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
		btree := NewBtree("./btree-storage/")
		btree.AddWord("a", nil)
		btree.AddWord("b", nil)
		btree.AddWord("c", nil)
		btree.AddWord("d", nil)
		btree.AddWord("e", nil)
		btree.AddWord("f", nil)
		btree.AddWord("g", nil)
		btree.AddWord("h", nil)
		btree.AddWord("i", nil)
		btree.AddWord("j", nil)
		btree.AddWord("k", nil)
		btree.AddWord("l", nil)
		btree.AddWord("n", nil)
		btree.AddWord("m", nil)
		btree.AddWord("o", nil)
		btree.AddWord("p", nil)
		btree.AddWord("q", nil)
		btree.AddWord("r", nil)
		btree.AddWord("s", nil)
		btree.AddWord("t", nil)
		btree.AddWord("u", nil)
		btree.AddWord("v", nil)
		btree.AddWord("w", nil)
		btree.AddWord("x", nil)
		btree.AddWord("y", nil)
		btree.AddWord("z", nil)
		btree.AddWord("za", nil)
		btree.AddWord("zc", nil)
		btree.AddWord("zd", nil)
		btree.AddWord("ze", nil)
		btree.AddWord("bc", nil)
		log.Println(btree)
		os.RemoveAll("/btree-root/")
	})

	t.Run("Test element finding", func(t *testing.T) {
		btree := NewBtree("./btree-storage")
		m_a := make(map[DocId][]PosIdx)
		m_e := make(map[DocId][]PosIdx)
		m_g := make(map[DocId][]PosIdx)
		m_l := make(map[DocId][]PosIdx)
		m_a[1] = []PosIdx{0, 1}
		m_e[1] = []PosIdx{0, 1}
		m_g[1] = []PosIdx{0, 1}
		m_l[1] = []PosIdx{0, 1}
		btree.AddWord("a", m_a)
		btree.AddWord("b", nil)
		btree.AddWord("c", nil)
		btree.AddWord("d", nil)
		btree.AddWord("e", m_e)
		btree.AddWord("f", nil)
		btree.AddWord("g", m_g)
		btree.AddWord("h", nil)
		btree.AddWord("i", nil)
		btree.AddWord("j", nil)
		btree.AddWord("k", nil)
		btree.AddWord("l", m_l)
		log.Println(btree)
		log.Println(btree.Find("a"))
		log.Println(btree.Find("e"))
		log.Println(btree.Find("g"))
		log.Println(btree.Find("l"))
		log.Println(btree.Find("h"))
		os.RemoveAll("/btree-root/")
	})

}
