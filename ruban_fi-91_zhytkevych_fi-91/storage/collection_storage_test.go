package storage

import (
	"os"
	"testing"
)

func TestCollectionStorageFS(t *testing.T) {
    t.Run("create collectionStorageFs", func(t *testing.T) {
        _, err := NewCollectionStorageFS("./")
        if err != nil {
            t.Fatal(err)
        }
    })

    t.Run("test load documets and collections", func(t *testing.T) {
        cs, err := NewCollectionStorageFS("./")
        if err != nil {
            t.Fatal(err)
        }

        if err := cs.CreateCollection("denis"); err != nil {
            t.Error(err)
        }

        if _, err := cs.AddDocument("denis", []byte("download APEX pls")); err != nil {
            t.Error(err)
        }

        if _, err := cs.AddDocument("denis", []byte("or go and play dota")); err != nil {
            t.Error(err)
        }

        if err := cs.CreateCollection("ivan"); err != nil {
            t.Error(err)
        }

        if _, err := cs.AddDocument("ivan", []byte("likes to play dumb games at the end of day")); err != nil {
            t.Error(err)
        }

        csT, err := NewCollectionStorageFS("./")
        if err != nil {
            t.Fatal(err)
        }

        if len(csT.Documents) != 3 {
            t.Errorf("CollectionStorageFS loaded %d documents instead of %d created\n", len(csT.Documents), 3)
        }
        if len(cs.Collections) != 2 {
            t.Errorf("CollectionStorageFS loaded %d collections instead of %d created\n", len(csT.Collections), 2)
        }

        os.Remove("collections.gob")
        os.Remove("doc-1.gob")
        os.Remove("doc-2.gob")
        os.Remove("doc-3.gob")
    })

    t.Run("document saves collection info", func(t *testing.T) {
        cs, err := NewCollectionStorageFS("./")
        if err != nil {
            t.Fatal(err)
        }

        if err := cs.CreateCollection("denis"); err != nil {
            t.Fatal(err)
        }

        if _, err := cs.AddDocument("denis", []byte("DENIS")); err != nil {
            t.Fatal(err)
        }
        defer func () {
            os.Remove(collectionsFileName)
            os.Remove("doc-1.gob")
        }()

        csN, err := NewCollectionStorageFS("./")
        if err != nil {
            t.Fatal(err)
        }

        doc, err := csN.GetDocumentById(1)
        if err := doc.Load(cs.path); err != nil {
            t.Error(err)
        }
        if doc.Collection.Name != "denis" {
            t.Error("document collection name does not match")
        }
    })
}
