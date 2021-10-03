package indexer

type Indexer interface {
}

type IndexerBtree struct {
}

func Indexation(document []bytes) map[int][]int
