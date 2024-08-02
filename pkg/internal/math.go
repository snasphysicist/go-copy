package internal

// Minimum returns the minimum value of all uint64s provided
func Minimum(n1 uint64, n2 ...uint64) uint64 {
	m := n1
	for _, n := range n2 {
		if n < m {
			m = n
		}
	}
	return m
}
