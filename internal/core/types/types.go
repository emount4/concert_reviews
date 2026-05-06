package core_types

type Nullable[T any] struct {
	Value *T
	Set   bool
}
