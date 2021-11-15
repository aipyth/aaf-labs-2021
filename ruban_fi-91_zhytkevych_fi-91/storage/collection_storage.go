package storage

import (
	"encoding/gob"
	"errors"
	"os"
	"path"
	"strconv"
)

const collectionsFileName = "collections.gob"

// CollectionStorage inteface describes needed methods of different collections storages
// implementations. This is essential for Domain.
type CollectionStorage interface {
    AddDocument(collectionName string, content []byte) (*Document, error)
    CreateCollection(name string) error
    GetDocumentById(id uint64) (*Document, error)
    ContainsCollection(name string) bool
    FindCollection(name string) *Collection
    GetDocuments() ([]*Document, error)    
}

// CollectionStorageFS implements storage based of saving documents and collections to file system.
type CollectionStorageFS struct {
    Collections []*Collection
    Documents map[uint64]*Document
    path string
}

// NewCollectionStorageFS creates new CollectionStorageFS that resides in specified path.
func NewCollectionStorageFS(p string) (*CollectionStorageFS, error) {
    cs := &CollectionStorageFS{
        Collections: make([]*Collection, 0),
        Documents: make(map[uint64]*Document),
        path: p,
    }
    if err := cs.loadCollections(); err != nil &&
        err.Error() != "open " + path.Join(cs.path, collectionsFileName) + ": no such file or directory" {
        return nil, err
    }
    if err := cs.loadDocuments(); err != nil {
        return nil, err
    }
    return cs, nil
}

// GetDocuments returns all documents in storage
func(cs *CollectionStorageFS) GetDocuments() ([]*Document, error) {
    docs := make([]*Document, 0,)
    for _, v := range cs.Documents {
        err := v.Load(cs.path)
        if err != nil {
            return docs, err
        }
        docs = append(docs, v)
    }
    return docs, nil
}

// saveCollections just saves all collections info to the file collections.gob
func (cs *CollectionStorageFS) saveCollections() error {
    f, err := os.Create(path.Join(cs.path, collectionsFileName))
    if err != nil {
        return err
    }
    encoder := gob.NewEncoder(f)
    return encoder.Encode(cs.Collections)
}
// loadCollections loads structs of collections saved to collections.gob directly into
// CollectionStorageFS
func (cs *CollectionStorageFS) loadCollections() error {
    f, err := os.Open(path.Join(cs.path, collectionsFileName))
    if err != nil {
        return err
    }
    decoder := gob.NewDecoder(f)
    return decoder.Decode(&cs.Collections)
}

// loadDocuments reads filenames in the directory and saves information that saved documents in this
// dir exists. Any document may be loaded later by Document.Load()
func (cs *CollectionStorageFS) loadDocuments() error {
    dirEntries, err := os.ReadDir(cs.path)
    if err != nil {
        return err
    }

    for _, v := range dirEntries {
        name := v.Name()
        if len(name) > 3 && name[len(name)-4:] == ".gob" &&
            name != collectionsFileName {
            docId, err := strconv.ParseUint(name[4:len(name)-4], 10, 64)
            if err != nil {
                // log.Println("[CollectionStorageFS:loadDocuments]", err)
                continue
            }

            cs.Documents[docId] = &Document{
                Id: docId,
            }
        }
    }
    return nil
}

func (cs *CollectionStorageFS) ContainsCollection(name string) bool {
    for _, v := range cs.Collections {
        if v.Name == name {
            return true
        }
    }
    return false
}

func (cs *CollectionStorageFS) FindCollection(name string) *Collection {
    for _, v := range cs.Collections {
        if v.Name == name {
            return v
        }
    }
    return nil
}

// AddDocument creates documents, adds it's entry to specified collection and saves data about
// changed collection and document to filesystem.
func (cs *CollectionStorageFS) AddDocument(collectionName string, content []byte) (*Document, error) {
    collection := cs.FindCollection(collectionName)
    if collection == nil {
        return nil, errors.New("cannot add document to unexisting collection")
    }

    document := &Document{
        Id: uint64(len(cs.Documents) + 1),
        Contents: content,
        Collection: collection,
    }

    err := document.Save(cs.path)
    if err != nil {
        return nil, err
    }
    collection.AddDocument(document.Id)
    err = cs.saveCollections()
    if err != nil {
        // TODO: delete created document file
        return nil, err
    }
    cs.Documents[document.Id] = document
    return document, nil
}

// CreateCollection creates new collection
func (cs *CollectionStorageFS) CreateCollection(name string) error {
    if cs.FindCollection(name) != nil {
        return errors.New("this collection already exist")
    }
    collection := NewCollection(name)
    cs.Collections = append(cs.Collections, collection)
    return cs.saveCollections()
} 

// GetDocumentById return document by id
func (cs *CollectionStorageFS) GetDocumentById(id uint64) (*Document, error) {
    doc, ok := cs.Documents[id]
    if !ok {
        return nil, errors.New("this document does not exist")
    }
    err := doc.Load(cs.path)
    if err != nil {
        return doc, err
    }
    return doc, nil
}


