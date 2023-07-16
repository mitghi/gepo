package storage

type Builder[T any] func() Storage[T]
type store[T any] map[string]T

type Storage[T any] interface {
	Add(string, T)
	Get(string) (T, bool)
	Delete(string)
	Exists(string) bool
	Keys() []string
	Values() []T
	Size() int
	Clone() Storage[T]
	Drain()
}

func New[T any]() Storage[T] {
	output := make(store[T])
	return output
}

func (s store[T]) Add(id string, value T) { s[id] = value }
func (s store[T]) Delete(id string)       { delete(s, id) }
func (s store[T]) Size() int              { return len(s) }

func (s store[T]) Get(id string) (value T, ok bool) {
	value, ok = s[id]
	return value, ok
}

func (s store[T]) Clone() Storage[T] {
	output := make(store[T])
	for k, v := range s {
		output[k] = v
	}
	return output
}

func (s store[T]) Exists(id string) (ok bool) {
	_, ok = s[id]
	return ok
}

func (s store[T]) Values() (values []T) {
	len := len(s)
	if len == 0 {
		return nil
	}
	values = make([]T, 0, len)
	for _, v := range s {
		values = append(values, v)
	}
	return values
}

func (s store[T]) Keys() (keys []string) {
	if len(s) == 0 {
		return nil
	}
	keys = make([]string, 0, len(s))
	for k, _ := range s {
		keys = append(keys, k)
	}
	return keys
}

func (s store[T]) Drain() {
	if len(s) == 0 {
		return
	}
	keys := s.Keys()
	for _, v := range keys {
		delete(s, v)
	}
	return
}
