package storage

import (
    "testing"
)

func TestCollection(t *testing.T) {
    t.Run("can add documents id", func(t *testing.T) {
        collection := NewCollection("newcol")
        collection.AddDocument(4)
        if collection.Documents[0] != 4 {
            t.Error("document id is not added to collection")
        }
    })

    t.Run("contains document", func(t *testing.T) {
        collection := NewCollection("newcol")
        collection.AddDocument(4)
        if !collection.Contains(4) {
            t.Error("collection says that documents is not in it but it contains it")
        }
        if collection.Contains(0) {
            t.Error("collection says that it contains the document that does not exist")
        }       
    })
}
