package inmemory

import (
	"hash/fnv"
	"iter"
	"strconv"
	"sync"

	"github.com/matthiasharzer/go-pokemon-viewer/view"
)

type Hashable interface {
	Hash() string
}

type View[T Hashable] struct {
	items []T
	hash  string

	mu sync.RWMutex
}

func NewView[T Hashable]() view.View[T] {
	return &View[T]{
		items: make([]T, 0),
		mu:    sync.RWMutex{},
	}
}

func (v *View[T]) recalculateHash() {
	h := fnv.New64a()

	for _, item := range v.items {
		_, _ = h.Write([]byte(item.Hash()))
	}

	v.hash = strconv.FormatUint(h.Sum64(), 36)
}

func (v *View[T]) Hash() string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.hash
}

func (v *View[T]) Insert(data ...T) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.items = append(v.items, data...)
	v.recalculateHash()

	return nil
}

func (v *View[T]) Update(predicate view.PredicateFn[T], updateFn func(*T) error) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	for i, item := range v.items {
		if predicate(item) {
			err := updateFn(&item)
			if err != nil {
				return err
			}
			v.items[i] = item
			v.recalculateHash()
			return nil
		}
	}
	return nil
}
func (v *View[T]) Find(predicate view.PredicateFn[T]) (*T, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	for _, item := range v.items {
		if predicate(item) {
			return &item, nil
		}
	}
	return nil, nil
}
func (v *View[T]) All() iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		v.mu.RLock()
		defer v.mu.RUnlock()

		for _, item := range v.items {
			if !yield(item, nil) {
				break
			}
		}
	}
}

func (v *View[T]) Delete(predicate view.PredicateFn[T]) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	for i, item := range v.items {
		if predicate(item) {
			v.items = append(v.items[:i], v.items[i+1:]...)
			v.recalculateHash()
			return nil
		}
	}
	return nil
}
func (v *View[T]) DeleteAll() error {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.items = make([]T, 0)
	v.recalculateHash()
	return nil
}
