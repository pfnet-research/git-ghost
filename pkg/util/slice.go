package util

type empty struct{}

func UniqueStringSlice(slice []string) []string {
	m := make(map[string]empty)

	for _, ele := range slice {
		m[ele] = empty{}
	}

	uniq := []string{}
	for i := range m {
		uniq = append(uniq, i)
	}

	return uniq
}
