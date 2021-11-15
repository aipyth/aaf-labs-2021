package storage

type Storage interface {
    CreateCollection(name string) error
    AddDocument(collection string) error
}

