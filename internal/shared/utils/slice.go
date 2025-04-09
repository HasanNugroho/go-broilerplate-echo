package utils

func Intersection[T comparable](a, b []T) []T {
	set := make(map[T]bool)
	result := []T{}
	seen := make(map[T]bool)

	for _, val := range a {
		set[val] = true
	}
	for _, val := range b {
		if set[val] && !seen[val] {
			result = append(result, val)
			seen[val] = true
		}
	}

	return result
}
