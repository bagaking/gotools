package reflectool

import "reflect"

type FieldHandler func(field *reflect.Value, fieldType reflect.StructField) error

func ForEachField(target interface{}, fn FieldHandler) error {
	r := reflect.ValueOf(target)
	elem := r
	if r.Kind() == reflect.Ptr {
		elem = r.Elem()
	}
	rType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field, fieldType := elem.Field(i), rType.Field(i)
		if err := fn(&field, fieldType); err != nil {
			return err
		}
	}
	return nil
}
