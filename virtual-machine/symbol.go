package virtualmachine

import "bytes"

type Symbol struct {
	name       []byte
	value_type byte
	value_ptr  int16
}

type SymbolTable []Symbol

func NewSymbolTable() SymbolTable {
	return make([]Symbol, 0)
}

func (st *SymbolTable) Search(name []byte) *Symbol {
	for i := range *st {
		if bytes.Equal((*st)[i].name, name) {
			return &(*st)[i]
		}
	}
	return nil
}

func (st *SymbolTable) Add(name []byte, value_type byte, value_ptr int16) {
	*st = append(*st, Symbol{name, value_type, value_ptr})
}
