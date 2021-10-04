package domain

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/indexer"
	"github.com/aipyth/aaf-labs-2021/ruban_fi-91_zhytkevych_fi-91/storage"
)

const defaultDomainConfigPath = "./"
const domainConfigName = "dddb-conf.json"

var defaultDomainConf = &DomainConf{
	StoragePath: "./",
	StorageType: "fs",
    IndexerPath: "./",
    IndexerType: "fs",
}

type Domain struct {
	CollectionStorage storage.CollectionStorage
	Indexer           indexer.Indexer
}

type DomainConf struct {
	StoragePath string `json:"storage_path"`
	StorageType string `json:"storage_type"`
	IndexerPath string `json:"indexer_path"`
	IndexerType string `json:"indexer_type"`
}

type SearchQuery struct {
	Keyword  string
	Prefix   string
	N        uint
	KeywordE string
}

func NewDomain() *Domain {
	domain := &Domain{}
	conf := initDomainConfiguration()

	var collectionsStorage storage.CollectionStorage
	var indexr indexer.Indexer
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

	switch conf.IndexerType {
	case "mem":
		panic("in memory indexer is not implemented")
	case "fs":
		indexr = indexer.NewIndexerBtree(conf.IndexerPath)
	}

	domain.CollectionStorage = collectionsStorage
	domain.Indexer = indexr

	return domain
}

func initDomainConfiguration() *DomainConf {
	f, err := os.Open(path.Join(defaultDomainConfigPath, domainConfigName))
	if err != nil {
		if err.Error() != "open "+domainConfigName+": no such file or directory" {
			panic(err)
		} else {
			return defaultDomainConf
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
	doc, err := d.CollectionStorage.AddDocument(collectionName, []byte(document))
	if err != nil {
		return err
	}

	err = d.Indexer.IndexDocument(doc.Id, []byte(document))
	if err != nil {
		return err
	}

	return nil
}

func (d *Domain) Search(q SearchQuery) []*storage.Document {
	resp := make([]*storage.Document, 0)

	switch {
	case q.KeywordE != "" && q.Keyword != "":
		docs, _ := d.Indexer.GetDocsByKeywords(q.Keyword, q.KeywordE, q.N)
		for _, v := range docs {
			doc, err := d.CollectionStorage.GetDocumentById(v)
			if err == nil {
				resp = append(resp, doc)
			}
		}

	case q.Prefix != "":
		docs, _ := d.Indexer.GetDocsByPrefix(q.Prefix)
		for _, v := range docs {
			doc, err := d.CollectionStorage.GetDocumentById(v)
			if err == nil {
				resp = append(resp, doc)
			}
		}

	case q.Keyword != "":
		docs, _ := d.Indexer.GetDocsByKeyword(q.Keyword)
		for _, v := range docs {
			doc, err := d.CollectionStorage.GetDocumentById(v)
			if err == nil {
				resp = append(resp, doc)
			}
		}
	}

	return resp
}
