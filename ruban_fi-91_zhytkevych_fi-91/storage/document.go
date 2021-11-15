package storage

import (
	"fmt"
	"os"
    "path"
    "errors"
    "encoding/gob"
)

type Document struct {
    Id uint64
    Contents []byte
    Collection *Collection
}

func (d *Document) Save(p string) error {
    info, err := os.Stat(p)
    if err != nil {
        return err
    }
    if !info.IsDir() {
        return errors.New("the path is not a directory")
    }
    fileName := fmt.Sprintf("doc-%d.gob", d.Id)
    f, err := os.Create(path.Join(p, fileName))
    if err != nil {
        return err
    }
    encoder := gob.NewEncoder(f)
    return encoder.Encode(d)
}

func (d *Document) Load(p string) error {
    info, err := os.Stat(p)
    if err != nil {
        return err
    }
    var f *os.File
    if info.IsDir() {
        fileName := fmt.Sprintf("doc-%d.gob", d.Id)
        f, err = os.Open(path.Join(p, fileName))
        if err != nil {
            return err
        }
    } else {
        f, err = os.Open(p)
        if err != nil {
            return err
        }
    }
    decoder := gob.NewDecoder(f)
    return decoder.Decode(d)
}

func (d *Document) String() string {
    return fmt.Sprintf("[%d] %s: %s", d.Id, d.Collection.Name, string(d.Contents))
}
