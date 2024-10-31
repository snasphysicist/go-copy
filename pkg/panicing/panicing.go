package panicing

// OnError panics if the error is not nil
func OnError(err error) {
	if err != nil {
		panic(err)
	}
}

// OnWrite panics if the error is not nil,
// designed to work with io.Writer.Write
func OnWriteError(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
