package indexer

type SetUint struct {
	data map[uint64]bool
}

func NewSetUint() *SetUint {
	data := make(map[uint64]bool, 0)
	return &SetUint{
		data,
	}
}

func (s *SetUint) Add(el uint64) {
	if !s.data[el] {
		s.data[el] = true
	}
}

func (s *SetUint) ToArray() []uint64 {
	arr := make([]uint64, 0)
	for id := range s.data {
		arr = append(arr, id)
	}
	return arr
}
