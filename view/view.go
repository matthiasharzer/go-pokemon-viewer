package view

import "iter"

type PredicateFn[T any] = func(T) bool

type ReadOnlyView[T any] interface {
	Find(predicate PredicateFn[T]) (*T, error)
	All() iter.Seq2[T, error]
	Hash() string
}

type View[T any] interface {
	ReadOnlyView[T]
	Insert(data ...T) error
	Update(predicate PredicateFn[T], updateFn func(*T) error) error
	Delete(predicate PredicateFn[T]) error
	DeleteAll() error
	ReplaceAll(data ...T) error
}
