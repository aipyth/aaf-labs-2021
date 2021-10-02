package storage

import (
	"strings"
)

type Btree struct {
	root *Sheet
}

func NewBtree() *Btree {
	return &Btree{
		root: NewSheet(),
	}
}

func (t *Btree) Find(word string) (*SheetElement, error) {
	currSheet := t.root
	for {
		el, childIndex, err := currSheet.Find(word)
		if err != nil && currSheet.Children != nil {
			currSheet = currSheet.Children[childIndex]
			continue
		}
		return el, err
	}
}

func splitSheetsByHalf(sheet *Sheet) (*Sheet, *Sheet, *SheetElement) {
	mid := len(sheet.Keys) / 2
	firstPart := NewSheet()
	secPart := NewSheet()
	firstPart.Keys = append(firstPart.Keys, sheet.Keys[:mid]...)
	secPart.Keys = append(secPart.Keys, sheet.Keys[mid+1:]...)
	if sheet.Children != nil {
		firstPart.AppendChildren(sheet.Children[:mid+1])
		secPart.AppendChildren(sheet.Children[mid+1:])
	}
	if sheet.Parent != nil {
		firstPart.Parent = sheet.Parent
		secPart.Parent = sheet.Parent
	}
	return firstPart, secPart, sheet.Keys[mid]
}

func indexOf(slice []*Sheet, item *Sheet) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}
	return -1
}

func appendToSheet(sheet *Sheet, elem *SheetElement) error {
	err := sheet.Add(elem.Key, elem.Data)
	if err == nil {
		return nil
	}
	if err.Error() == "Sheet full" {
		firstPart, secPart, popped := splitSheetsByHalf(sheet)
		cmp := strings.Compare(elem.Key, popped.Key)
		switch cmp {
		case -1:
			firstPart.Add(elem.Key, elem.Data)
		case 1:
			secPart.Add(elem.Key, elem.Data)
		}
		if sheet.Parent != nil {
			newChildren := make([]*Sheet, 0)
			i := indexOf(sheet.Parent.Children, sheet)
			newChildren = append(newChildren, sheet.Parent.Children[:i]...)
			newChildren = append(newChildren, firstPart)
			newChildren = append(newChildren, secPart)
			newChildren = append(newChildren, sheet.Parent.Children[i+1:]...)
			sheet.Parent.Children = newChildren
			err = appendToSheet(sheet.Parent, popped)
			return nil
		} else {
			keys := make([]*SheetElement, 0, cap(sheet.Keys))
			sheet.Keys = append(keys, sheet.Keys[len(sheet.Keys)/2:len(sheet.Keys)/2+1]...)
			firstPart.Parent = sheet
			secPart.Parent = sheet
			sheet.Children = []*Sheet{firstPart, secPart}
		}
	} else {
		return nil
	}
	return nil
}

func (t *Btree) AddWord(word string, data map[DocId][]PosIdx) error {
	currSheet := t.root
	index := 0
	for currSheet.Children != nil {
		cmp := strings.Compare(currSheet.Keys[index].Key, word)
		switch cmp {
		case 0:
			for k, v := range data {
				currSheet.Keys[index].Data[k] = v
			}
		case 1:
			currSheet = currSheet.Children[index]
			index = 0
			continue
		case -1:
			if index == len(currSheet.Keys)-1 {
				currSheet = currSheet.Children[index+1]
				index = 0
				continue
			}
			index++
			continue
		}
	}
	appendToSheet(currSheet, &SheetElement{
		Key:  word,
		Data: data,
	})
	return nil
}

func childrenToString(sheet *Sheet) string {
	str := sheet.String()
	if len(sheet.Children) != 0 {
		str += ": {\n"
		for _, v := range sheet.Children {
			str += "	" + childrenToString(v) + "\n"
		}
		str += "},\n"
	}
	return str
}

func (t *Btree) String() string {
	return childrenToString(t.root)
}
