package stream

type Terminator interface {
	Terminate()
}

func (s *stream) Terminate() {
	s.terminator <- struct{}{}
}
