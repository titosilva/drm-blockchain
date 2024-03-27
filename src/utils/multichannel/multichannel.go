package multichannel

type MultiChannel[T any] struct {
	subscribers []chan T
	closed      bool
}

func New[T any]() *MultiChannel[T] {
	mc := new(MultiChannel[T])
	mc.subscribers = make([]chan T, 0)
	return mc
}

func (mc *MultiChannel[T]) Subscribe() <-chan T {
	if mc.closed {
		panic("channel is closed!")
	}

	c := make(chan T)
	mc.subscribers = append(mc.subscribers, c)
	return c
}

func (mc *MultiChannel[T]) Notify(t T) {
	if mc.closed {
		panic("channel is closed!")
	}

	for _, c := range mc.subscribers {
		go func(c0 chan T) {
			c0 <- t
		}(c)
	}
}

func (mc *MultiChannel[T]) Close() {
	if mc.closed {
		return
	}

	for _, c := range mc.subscribers {
		close(c)
	}
}
