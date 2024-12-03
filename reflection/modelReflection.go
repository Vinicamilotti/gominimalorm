package reflection

import "reflect"

func StructFieldPtr(x interface{}) []any {
	xv := reflect.ValueOf(x).Elem()
	ret := []any{}

	for i := 0; i < xv.NumField(); i++ {
		if !xv.Field(i).Addr().CanInterface() {
			continue
		}

		v := xv.Field(i).Addr().Interface()
		ret = append(ret, v)
	}
	return ret
}
