package domain

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDomain(t *testing.T) {
    // removes all .gob files
    t.Cleanup(func() {
        dirEntries, _ := os.ReadDir("./")
        for _, v := range dirEntries {
            if filepath.Ext(v.Name()) == ".gob" {
                os.Remove(v.Name())
            }
        }
    })

    t.Run("create domain", func(t *testing.T) {
        domain := NewDomain()
        if domain == nil {
            t.Fatal("domain is nil")
        }
        if domain.CollectionStorage == nil {
            t.Fatal("domain's collection storage is nil")
        }
    })

    t.Run("creates collection", func (t *testing.T) {
        domain := NewDomain()
        const newCollectionName = "denis na pidzhake"
        err := domain.CreateCollection(newCollectionName)
        if err != nil {
            t.Fatal(err)
        }
        if !domain.CollectionStorage.ContainsCollection(newCollectionName) {
            t.Fatal("collection storage does not contain new collection")
        }
    })

    t.Run("inserts document", func (t *testing.T) {
        domain := NewDomain()
        const newCollectionName = "denis na pidzhake"
        err := domain.CreateCollection(newCollectionName)
        if err != nil {
            t.Fatal(err)
        }

        err = domain.InsertDocument(newCollectionName, "pidzhak na denise")
        if err != nil {
            t.Fatal(err)
        }

        doc, err := domain.CollectionStorage.GetDocumentById(1)
        if err != nil {
            t.Fatal(err)
        }
        if doc.Collection.Name != newCollectionName {
            t.Error("document has wrong collection")
        }

        col := domain.CollectionStorage.FindCollection(newCollectionName)
        if col == nil {
            t.Fatal("created collection is not found")
        }
        if !col.Contains(doc.Id) {
            t.Error("the collection does not contain the document")
        }
    })

}
