package storage

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"strings"
)

const t = 10

type DocId int
type PosIdx int

type SheetElement struct {
    Key string
    Data map[DocId][]PosIdx
}

type Sheet struct {
    Keys []*SheetElement
}

func NewSheet() *Sheet {
    return &Sheet{
        Keys: make([]*SheetElement, 0, 2*t-1),
    }
}

func (s *Sheet) String() string {
    return fmt.Sprintf("%v", s.Keys)
}

func (s *Sheet) Encode() bytes.Buffer {
    var buf bytes.Buffer
    encoder := gob.NewEncoder(&buf)
    err := encoder.Encode(s)
    if err != nil {
        log.Println(err)
    }

    return buf
}

func (s *Sheet) Add(key string, data map[DocId][]PosIdx) error {
    if len(s.Keys) == cap(s.Keys) {
        return errors.New("Sheet full")
    }

    for i, v := range s.Keys {
        if v == nil {
            s.Keys[i] = &SheetElement{
                Key: key,
                Data: data,
            }
            break
        }

        cmp := strings.Compare(v.Key, key)
        switch cmp {
        case 0:
            // update data -- merge two maps
            for k, v := range data {
                s.Keys[i].Data[k] = v
            }
            break
        case 1:
            s.Keys = append(
                s.Keys[:i],
                &SheetElement{
                    Key: key,
                    Data: data,
                },
            )
            s.Keys = append(s.Keys, s.Keys[i:len(s.Keys)-1]...)
            break
        case -1:
            continue
        }
    }
    return nil
}

