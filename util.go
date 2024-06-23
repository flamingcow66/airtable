package airtable

func chunk[T any](items []T, chunkSize int) [][]T {
	chunks := [][]T{}

	for chunkSize < len(items) {
		chunks = append(chunks, items[0:chunkSize:chunkSize])
		items = items[chunkSize:]
	}

	return append(chunks, items)
}
