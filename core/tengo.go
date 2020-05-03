package core

import "github.com/d5/tengo/v2"

func awdawdo(array []tengo.Object) []interface{} {
	res := make([]interface{}, len(array))

	for i, v := range array {
		if v.TypeName() == "string" {
			res[i] = v.(*tengo.String).Value
		}

		if v.TypeName() == "int" {
			res[i] = int(v.(*tengo.Int).Value)
		}
	}

	return res
}
