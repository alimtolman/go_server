package DataBase

import (
	"Server/Helpers/File"
	"encoding/json"
)

type JsonTable[T Indexable] struct {
	directory string
	path      string
	lastId    uint32
	data      []T
}

/* public */

func NewJsonTable[T Indexable](name string) *JsonTable[T] {
	result := &JsonTable[T]{}

	result.directory = File.CurrentDirectory() + "/db"
	result.path = result.directory + "/" + name + ".json"
	result.lastId = 0
	result.data = make([]T, 0)

	if File.Exists(result.path) {
		bytes := File.ReadBytes(result.path)
		json.Unmarshal(bytes, &result.data)
	} else {
		File.CreateDirectory(result.directory)
		File.WriteText(result.path, "[]")
	}

	result.updateLastId()

	return result
}

func (table *JsonTable[T]) Add(item T) {
	item.GetId()
	table.data = append(table.data, item)

	table.updateLastId()
	table.save()
}

func (table *JsonTable[T]) AddOrSet(item T) {
	index := table.findIndex(item)

	if index >= 0 {
		table.data[index] = item
	} else {
		table.data = append(table.data, item)
	}

	table.updateLastId()
	table.save()
}

func (table *JsonTable[T]) Get(predicate func(T) bool) T {
	for _, item := range table.data {
		if predicate(item) {
			return item
		}
	}

	return table.data[0]
}

func (table *JsonTable[T]) GetAt(index uint32) T {
	return table.data[index]
}

func (table *JsonTable[T]) GetAll() []T {
	return table.data
}

func (table *JsonTable[T]) NextId(predicate func(T) bool) uint32 {
	if predicate == nil {
		return table.lastId + 1
	}

	for _, item := range table.data {
		if predicate(item) {
			return item.GetId()
		}
	}

	return table.lastId + 1
}

func (table *JsonTable[T]) Remove(predicate func(T) bool) {
	data := make([]T, 0)

	if predicate != nil {
		for _, item := range table.data {
			if !predicate(item) {
				data = append(data, item)
			}
		}
	}

	table.data = data

	table.updateLastId()
	table.save()
}

func (table *JsonTable[T]) RemoveAll() {
	table.Remove(nil)
}

/* private */

func (table *JsonTable[T]) findIndex(item T) int {
	index := -1

	for i, existing := range table.data {
		if existing.GetId() == item.GetId() {
			index = i
			break
		}
	}

	return index
}

func (table *JsonTable[T]) save() {
	data, _ := json.Marshal(table.data)

	File.WriteBytes(table.path, data)
}

func (table *JsonTable[T]) updateLastId() {
	if len(table.data) != 0 {
		table.lastId = table.data[len(table.data)-1].GetId()
	}
}
