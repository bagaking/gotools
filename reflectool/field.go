package reflectool

import "reflect"

// IsFieldExported returns whether a field are exported
// @see go/src/reflect/type.go reflect.StructField
func IsFieldExported(fieldType reflect.StructField) bool {
	return fieldType.PkgPath == ""
}
