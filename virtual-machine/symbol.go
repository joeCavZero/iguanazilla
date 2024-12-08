package virtualmachine

type Symbol struct {
	name      string
	value_ptr int16
}

type SymbolTable []Symbol

func NewSymbolTable() SymbolTable {
	return make([]Symbol, 0)
}

func (st *SymbolTable) Search(name string) *Symbol {
	for i := range *st {
		if name == (*st)[i].name {
			return &(*st)[i]
		}
	}
	return nil
}

func (st *SymbolTable) Add(name string, value_ptr int16) {
	*st = append(*st, Symbol{name, value_ptr})
}
