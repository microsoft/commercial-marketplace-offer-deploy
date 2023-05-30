package component

import "sync"

type FactoryFunc[T any] func() (T, error)

// provider of singleton instance of T
type Provider[T any] interface {
	Get() T
	Error() error
}

type provider[T any] struct {
	instance T
	factory  FactoryFunc[T]
	once     sync.Once
	err      error
}

// Get implements Provider
func (p *provider[T]) Get() T {
	p.once.Do(func() {
		p.instance, p.err = p.factory()
	})
	return p.instance
}

func (p *provider[T]) Error() error {
	return p.err
}

func NewProvider[T any](factory FactoryFunc[T]) Provider[T] {
	return &provider[T]{
		factory: factory,
	}
}
