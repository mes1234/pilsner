package stream

type Delegate struct {
	Channel chan Item
	Name    string
}

func NewDelegate(name string) *Delegate {

	subscriber := make(chan Item, BufferSize)

	return &Delegate{
		Name:    name,
		Channel: subscriber,
	}
}
