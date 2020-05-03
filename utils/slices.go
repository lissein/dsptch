package utils

func ToIntSlice(original []interface{}) []int {
	res := make([]int, len(original))

	for i, d := range original {
		res[i] = d.(int)
	}

	return res
}
