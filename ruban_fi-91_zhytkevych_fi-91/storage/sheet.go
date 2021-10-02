package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"strings"
)

const t = 3

type DocId int
type PosIdx int

type SheetElement struct {
	Key  string
	Data map[DocId][]PosIdx
}

type Sheet struct {
	Keys     []*SheetElement
	Children []*Sheet
	Parent   *Sheet
}

func NewSheet() *Sheet {
	keys := make([]*SheetElement, 0, 2*t-1)
	return &Sheet{
		Keys:     keys,
		Children: nil,
		Parent:   nil,
	}
}

// func (s *Sheet) String() string {
// 	return fmt.Sprintf("%v", s.Keys)
// }

func (s *Sheet) Encode() bytes.Buffer {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(s)
	if err != nil {
		log.Println(err)
	}

	return buf
}

func (s *Sheet) AddChild(sheet *Sheet, pos int) error {
	s.Children = append(
		s.Children[:pos],
		sheet,
	)
	s.Children = append(s.Children, s.Children[pos:len(s.Children)-1]...)
	return nil
}

func (s *Sheet) AppendChildren(children []*Sheet) {
	s.Children = children
	for _, ch := range children {
		ch.Parent = s
	}
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

func (s *Sheet) Add(key string, data map[DocId][]PosIdx) error {
	if len(s.Keys) == cap(s.Keys) {
		return errors.New("Sheet full")
	}
	for i, v := range s.Keys {
		// if v == nil {
		// 	s.Keys[i] = &SheetElement{
		// 		Key:  key,
		// 		Data: data,
		// 	}
		// 	break
		// }
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
