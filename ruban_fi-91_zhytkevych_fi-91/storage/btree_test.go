package storage

import (
	"bytes"
	"encoding/gob"
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

func TestSheetEncode(t *testing.T) {
    t.Run("encode empty sheet", func(t *testing.T) {
        sh := NewSheet()
        bts := sh.Encode()
        encoded := encode(t, &Sheet{
            Keys: make([]*SheetElement, 0),
        })
        if bytes.Compare(bts.Bytes(), encoded.Bytes()) != 0 {
            t.Errorf("Encoded data does not match\n")
        }
    })
}
