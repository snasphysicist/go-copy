package internal

// PanicOnError panics if the error is not nil
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// PanicOnWrite panics if the error is not nil,
// designed to work with io.Writer.Write
func PanicOnWriteError(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
