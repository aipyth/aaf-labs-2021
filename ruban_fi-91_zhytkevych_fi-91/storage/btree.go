package storage

import (
	"log"
	"strings"
)

type Btree struct {
	root *Sheet
	path string
}

func NewBtree(storagePath string) *Btree {
	root, err := ReadSheet("root", storagePath)
	if err != nil {
		root = NewSheet(storagePath)
		root.Name = FilePath("root")
		WriteSheet(root, storagePath)
	}
	return &Btree{
		root: root,
		path: storagePath,
	}
}

func (t *Btree) splitSheetsByHalf(sheet *Sheet) (*Sheet, *Sheet, *SheetElement) {
	mid := len(sheet.Keys) / 2
	firstPart := NewSheet(t.path)
	firstPart.Keys = append(firstPart.Keys, sheet.Keys[:mid]...)
	if sheet.Children != nil {
		firstPart.AppendChildren(sheet.Children[:mid+1])
	}
	if sheet.Parent != nil {
		firstPart.Parent = sheet.Parent
	}
	WriteSheet(firstPart, t.path)

	secPart := NewSheet(t.path)
	secPart.Keys = append(secPart.Keys, sheet.Keys[mid+1:]...)
	if sheet.Children != nil {
		secPart.AppendChildren(sheet.Children[mid+1:])
	}
	if sheet.Parent != nil {
		secPart.Parent = sheet.Parent
	}
	WriteSheet(secPart, t.path)
	return firstPart, secPart, sheet.Keys[mid]
}

func indexOf(slice []FilePath, item FilePath) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func (t *Btree) appendToSheet(sheet *Sheet, elem *SheetElement) error {
	err := sheet.Add(elem.Key, elem.Data)
	if err == nil {
		WriteSheet(sheet, t.path)
		return nil
	}
	if err.Error() == "Sheet full" {
		firstPart, secPart, popped := t.splitSheetsByHalf(sheet)
		cmp := strings.Compare(elem.Key, popped.Key)
		switch cmp {
		case -1:
			firstPart.Add(elem.Key, elem.Data)
			WriteSheet(firstPart, t.path)
		case 1:
			secPart.Add(elem.Key, elem.Data)
			WriteSheet(secPart, t.path)
		}
		if sheet.Parent != nil {
			newChildren := make([]FilePath, 0)
			i := indexOf(sheet.Parent.Children, sheet.Name)
			newChildren = append(newChildren, sheet.Parent.Children[:i]...)
			newChildren = append(newChildren, firstPart.Name)
			newChildren = append(newChildren, secPart.Name)
			newChildren = append(newChildren, sheet.Parent.Children[i+1:]...)
			sheet.Parent.Children = newChildren
			err = t.appendToSheet(sheet.Parent, popped)
			WriteSheet(sheet, t.path)
			return nil
		} else {
			keys := make([]*SheetElement, 0, cap(sheet.Keys))
			sheet.Keys = append(keys, sheet.Keys[len(sheet.Keys)/2:len(sheet.Keys)/2+1]...)
			firstPart.Parent = sheet
			secPart.Parent = sheet
			sheet.Children = []FilePath{firstPart.Name, secPart.Name}
			WriteSheet(sheet, t.path)
		}
	} else {
		return err
	}
	return nil
}

func (t *Btree) Find(word string) (*SheetElement, error) {
	currSheet := t.root
	for {
		el, childIndex, err := currSheet.Find(word)
		if err != nil && currSheet.Children != nil {
			currSheet, err = ReadSheet(currSheet.Children[childIndex], t.path)
			if err != nil {
				log.Println(err)
			}
			continue
		}
		return el, err
	}
}

func (t *Btree) AddIndex(word string, data map[uint64][]int) error {
	currSheet := t.root
	index := 0
	for currSheet.Children != nil {
		cmp := strings.Compare(currSheet.Keys[index].Key, word)
		switch cmp {
		case 0:
			for k, v := range data {
				currSheet.Keys[index].Data[k] = v
			}
			return nil
		case 1:
			nextSheet, err := ReadSheet(currSheet.Children[index], t.path)
			if err != nil {
				log.Println(err)
			}
			nextSheet.Parent = currSheet
			currSheet = nextSheet
			index = 0
			continue
		case -1:
			if index == len(currSheet.Keys)-1 {
				nextSheet, err := ReadSheet(currSheet.Children[index+1], t.path)
				if err != nil {
					log.Println(err)
				}
				nextSheet.Parent = currSheet
				currSheet = nextSheet
				index = 0
				continue
			}
			index++
			continue
		}
	}
	t.appendToSheet(currSheet, &SheetElement{
		Key:  word,
		Data: data,
	})
	return nil
}

func (t *Btree) AddIndexes(indexedDoc map[string]map[uint64][]int) error {
	for word, data := range indexedDoc {
		err := t.AddIndex(word, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func childrenToString(sheet *Sheet, path string) string {
	str := sheet.String()
	if len(sheet.Children) != 0 {
		str += ": {\n"
		for _, v := range sheet.Children {
			nextSheet, err := ReadSheet(v, path)
			if err != nil {
				log.Println(err)
			}
			str += "	" + childrenToString(nextSheet, path) + "\n"
		}
		str += "},\n"
	}
	return str
}

func (t *Btree) String() string {
	return childrenToString(t.root, t.path)
}
