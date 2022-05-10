package stream

type Iterator interface {
	Next() (next bool, item Item)
}

type iterator struct {
	position int
	pointer  int
	data     *[]Item
}

func (i *iterator) Next() (next bool, item Item) {

	if i.pointer <= i.position {
		i.pointer++
		return true, (*i.data)[i.pointer-1]
	} else {
		return false, Item{}
	}

}

func NewIterator(data *[]Item) Iterator {
	return &iterator{
		data:     data,
		position: len(*data) - 1,
	}
}
