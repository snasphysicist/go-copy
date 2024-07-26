package main

// assert panics with message if predicate returns false
func assert(predicate func() bool, message string) {
	if !predicate() {
		panic(message)
	}
}
