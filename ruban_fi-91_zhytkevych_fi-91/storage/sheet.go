package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
)

const t = 10

type FilePath string

type SheetElement struct {
	Key  string
	Data map[uint64][]int
}

type Sheet struct {
	Name     FilePath
	Keys     []*SheetElement
	Children []FilePath
	Parent   *Sheet
}

func NewSheet(folder string) *Sheet {
	keys := make([]*SheetElement, 0, 2*t-1)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Println(err)
	}
	name := fmt.Sprintf("%d", len(files))
	return &Sheet{
		Name:     FilePath(name),
		Keys:     keys,
		Children: nil,
		Parent:   nil,
	}
}

func serialize(s *Sheet) (bytes.Buffer, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(s)
	if err != nil {
		return buf, err
	}
	return buf, nil
}

func deserialize(file *os.File, sheet *Sheet) (*Sheet, error) {
	decoder := gob.NewDecoder(file)
	err := decoder.Decode(sheet)
	if err != nil {
		return nil, err
	}
	return sheet, nil
}

func ReadSheet(filePath FilePath, folder string) (*Sheet, error) {
	file, err := os.Open(path.Join(folder, "btr-"+string(filePath)+".gob"))
	if err != nil {
		return nil, err
	}
	sheet := NewSheet(folder)
	sheet, err = deserialize(file, sheet)
	if err != nil {
		return nil, err
	}
	return sheet, nil
}

func WriteSheet(sheet *Sheet, folder string) {
	sheet.Parent = nil
	buffer, err := serialize(sheet)
	err = os.WriteFile(path.Join(folder, "btr-"+string(sheet.Name)+".gob"), buffer.Bytes(), 0700)
	if err != nil {
		log.Println(err)
	}
}

func (s *Sheet) AddChild(sheet *Sheet, pos int) error {
	s.Children = append(
		s.Children[:pos],
		sheet.Name,
	)
	s.Children = append(s.Children, s.Children[pos:len(s.Children)-1]...)
	return nil
}

func (s *Sheet) AppendChildren(children []FilePath) {
	s.Children = children
}

func (s *Sheet) Find(key string) (*SheetElement, int, error) {
	length := len(s.Keys)
	cmp := strings.Compare(key, s.Keys[length-1].Key)
	switch cmp {
	case 0:
		return s.Keys[length-1], length - 1, nil
	case 1:
		return nil, length, errors.New("Not found")
	}
	i := sort.Search(length, func(i int) bool { return s.Keys[i].Key >= key })
	if i < length && s.Keys[i].Key == key {
		return s.Keys[i], i, nil
	} else {
		return nil, i, errors.New("Not found")
	}
}

func (s *Sheet) SearchMatches(prefix string) ([]*SheetElement, []int, error) {
	elements := make([]*SheetElement, 0)
	childIndexes := make([]int, 0)
	length := len(s.Keys)
	cmp := strings.Compare(prefix, s.Keys[length-1].Key)
	if cmp == 1 {
		return elements, []int{length}, nil
	}
	i := sort.Search(length, func(i int) bool { return s.Keys[i].Key >= prefix })
	if i < length && strings.HasPrefix(s.Keys[i].Key, prefix) {
		for {
			elements = append(elements, s.Keys[i])
			childIndexes = append(childIndexes, i)
			i++
			if !strings.HasPrefix(s.Keys[i].Key, prefix) || i == length {
				return elements, append(childIndexes, i), nil
			}
		}
	} else {
		return elements, []int{i}, nil
	}
}

func (s *Sheet) Add(key string, data map[uint64][]int) error {
	if len(s.Keys) == cap(s.Keys) {
		return errors.New("Sheet full")
	}
	for i, v := range s.Keys {
		cmp := strings.Compare(v.Key, key)
		switch cmp {
		// update data -- merge two maps
		case 0:
			for k, v := range data {
				s.Keys[i].Data[k] = v
			}
			return nil
		case 1:
			newKeys := make([]*SheetElement, 0, cap(s.Keys))
			toAdd := &SheetElement{
				Key:  key,
				Data: data,
			}
			newKeys = append(newKeys, s.Keys[:i]...)
			newKeys = append(newKeys, toAdd)
			newKeys = append(newKeys, s.Keys[i:]...)
			s.Keys = newKeys
			return nil
		case -1:
			continue
		}
	}

	s.Keys = append(s.Keys, &SheetElement{
		Key:  key,
		Data: data,
	})
	return nil
}

func (s *Sheet) String() string {
	key_data_string := "["
	for i, el := range s.Keys {
		key_data_string += "\"" + el.Key + "\": {"
		j := 0
		for docid, posids := range el.Data {
			key_data_string += fmt.Sprint(docid) + ": " + fmt.Sprint(posids)
			if j != len(el.Data)-1 {
				key_data_string += ", "
			}
			j++
		}
		key_data_string += "}"
		if i != len(s.Keys)-1 {
			key_data_string += "; "
		}
	}
	return key_data_string + "]"
}
