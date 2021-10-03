package domain

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

const defaultDomainConfigPath = "./"
const domainConfigName = "dddb-conf.json"

type Domain struct {
    CollectionStorage storage.CollectionStorage
}

type DomainConf struct {
    StoragePath string `json:"storage_path"`
    StorageType string `json:"storage_type"`
}

func NewDomain() *Domain {
    domain := &Domain{}
    conf := initDomainConfiguration()

    var collectionsStorage storage.CollectionStorage
    var err error

    switch conf.StorageType {
    case "mem":
        panic("in memory storage is not implemented")
    case "fs":
        collectionsStorage, err = storage.NewCollectionStorageFS(conf.StoragePath)
        if err != nil {
            panic(err)
        }
    }
    
    domain.CollectionStorage = collectionsStorage

    return domain
}

func initDomainConfiguration() *DomainConf {
    f, err := os.Open(path.Join(defaultDomainConfigPath, domainConfigName))
    if err != nil {
        if err.Error() != "open "  + domainConfigName + ": no such file or directory" {
           panic(err)
        } else {
            return &DomainConf{
                StoragePath: "./",
                StorageType: "fs",
            }
        }
    }

    conf := &DomainConf{}
    err = json.NewDecoder(f).Decode(conf)
    if err != nil {
        panic(err)
    }

    return conf
}

// CreateCollection creates new collection is storage with non empty name
func (d *Domain) CreateCollection(name string) error {
    if name == "" {
        return errors.New("collection name is empty")
    }
    return d.CollectionStorage.CreateCollection(name)
}

// InsertDocument adds non empty document to storage and indexes it's words
func (d *Domain) InsertDocument(collectionName string, document string) error {
    if document == "" {
        return errors.New("document cannot be empty")
    }
    _, err := d.CollectionStorage.AddDocument(collectionName, []byte(document))
    if err != nil {
        return err
    }
    // TODO: index words
    return nil
}
