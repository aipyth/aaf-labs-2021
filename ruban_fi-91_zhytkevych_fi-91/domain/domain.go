package domain

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *SearchQuery) String() string {
	switch {
	case s.Keyword != "" && s.Prefix == "" && s.N == 0 && s.KeywordE == "":
		return fmt.Sprintf("{keywordsearch:%s}", s.Keyword)
	case s.Prefix != "" && s.Keyword == "" && s.N == 0 && s.KeywordE == "":
		return fmt.Sprintf("{prefixsearch:%s}", s.Prefix)
	case s.Keyword != "" && s.KeywordE != "" && s.Prefix == "":
		return fmt.Sprintf("{nsearch:%s %d %s}", s.Keyword, s.N, s.KeywordE)
	default:
		return fmt.Sprintf("{unknownsearch:%s; %s; %s; %d}", s.Keyword, s.Prefix, s.KeywordE, s.N)
	}
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

	err = d.Indexer.IndexDocument(doc.Id, collectionName, []byte(document))
	if err != nil {
		return err
	}

	return nil
}

func (d *Domain) Search(collectionName string, q SearchQuery) []*storage.Document {
	searchIds := make([]uint64, 0)
	documents := make([]*storage.Document, 0)

	switch {
	case q.KeywordE != "" && q.Keyword != "":
		ids, _ := d.Indexer.GetDocsByKeywords(
            collectionName,
			q.Keyword,
			q.KeywordE,
			q.N,
		)
		searchIds = append(searchIds, ids...)
	case q.Prefix != "":
		ids, _ := d.Indexer.GetDocsByPrefix(collectionName, q.Prefix)
		searchIds = append(searchIds, ids...)
	case q.Keyword != "":
		ids, _ := d.Indexer.GetDocsByKeyword(collectionName, q.Keyword)
		searchIds = append(searchIds, ids...)
	default:
		docs, _ := d.CollectionStorage.GetDocuments()
		documents = append(documents, docs...)
	}

	for _, v := range searchIds {
		doc, err := d.CollectionStorage.GetDocumentById(v)
		if err == nil {
			documents = append(documents, doc)
		}
	}

	return documents
}

func (d *Domain) IndexerRepresentationString() string {
    return d.Indexer.BtreesString()
}
