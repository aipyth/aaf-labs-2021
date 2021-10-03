package storage


type Collection struct {
    Name string
    Documents []uint64
}

func NewCollection(name string) *Collection {
    return &Collection{
        Name: name,
        Documents: []uint64{},
    }
}

func (c *Collection) AddDocument(id uint64) {
    c.Documents = append(c.Documents, id)
}

func (c *Collection) Contains(doc uint64) bool {
    for _, v := range c.Documents {
        if v == doc {
            return true
        }
    }
    return false
}

