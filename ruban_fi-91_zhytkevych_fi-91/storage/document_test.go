package storage

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
)

func TestDocument(t *testing.T) {
    t.Run("document save and load with collection = nil", func(t *testing.T) {
        doc := &Document{
            Id: 1,
            Contents: []byte("denis go dota"),
            Collection: nil,
        }

        pathToSave := "./"
        err := doc.Save(pathToSave)
        if err != nil {
            t.Fatal(err)
        }
        
        loadedDoc := &Document{Id:1}
        err = loadedDoc.Load(pathToSave)
        if err != nil {
            t.Fatal(err)
        }

        os.Remove(path.Join(
            pathToSave,
            fmt.Sprintf("doc-%d.gob", doc.Id),
        ))

        if loadedDoc.Id != doc.Id {
            t.Error("doc and loaded doc ids does not match")
        }
        if bytes.Compare(loadedDoc.Contents, doc.Contents) != 0 {
            t.Error("doc and loaded doc contents does not match")
        }
        if loadedDoc.Collection != doc.Collection {
            t.Error("doc and loaded doc collection does not match")
        }
    })
    
}
