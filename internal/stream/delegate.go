package stream

type Delegate struct {
	channel chan Item
	name    string
}

func NewDelegate(channel chan Item, name string) *Delegate {
	return &Delegate{
		name:    name,
		channel: channel,
	}
}
