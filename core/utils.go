package core

func toIntSlice(original []interface{}) []int {
	res := make([]int, len(original))

	for i, d := range original {
		res[i] = int(d.(int64))
	}

	return res
}
