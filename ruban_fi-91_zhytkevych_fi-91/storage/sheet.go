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
	"strings"
)

const t = 20

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
	file, err := os.Open(path.Join(folder, string(filePath)+".gob"))
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
	err = os.WriteFile(path.Join(folder, string(sheet.Name)+".gob"), buffer.Bytes(), 0700)
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
	for i, v := range s.Keys {
		cmp := strings.Compare(v.Key, key)
		switch cmp {
		case 0:
			return v, i, nil
		case 1:
			return nil, i, errors.New("Not found")
		case -1:
			continue
		}
	}
	return nil, len(s.Keys), errors.New("Not found")
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
	str := "["
	for i, el := range s.Keys {
		if i != len(s.Keys)-1 {
			str += el.Key + ", "
		} else {
			str += el.Key
		}
	}
	return str + "]"
}
