package internal

// min returns the minimum value of all uint64s provided
func min(n1 uint64, n2 ...uint64) uint64 {
	m := n1
	for _, n := range n2 {
		if n < n1 {
			m = n
		}
	}
	return m
}
