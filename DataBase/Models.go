package DataBase

/* IndexedType */

type Indexable interface {
	GetId() uint32
}

type IndexedType struct {
	Id uint32
}

func NewIndexedType(id uint32) IndexedType {
	return IndexedType{id}
}

func (item IndexedType) GetId() uint32 {
	return item.Id
}

/* AppData */

type AppData struct {
	IndexedType
	Data string
}

func NewAppData(id uint32, data string) AppData {
	return AppData{NewIndexedType(id), data}
}
