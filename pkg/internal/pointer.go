package internal

// From allows to easily create pointers from
// stuff where it's not a trivial as it should
// be, e.g. from function calls. e.g. instead of
// f := Foo(); &f -> From(Foo())
func From[T any](v T) *T {
	return &v
}
