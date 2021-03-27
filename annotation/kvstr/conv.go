package kvstr

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// ConvStrToPlainType takes any string str and convert it into any plain types (possibly truncated)
func ConvStrToPlainType(str string, p reflect.Type) (interface{}, error) {
	switch kind := p.Kind(); kind {
	case reflect.String:
		return str, nil
	case reflect.Bool:
		return strings.ToLower(str) != "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valI64, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		if kind == reflect.Int64 {
			return valI64, nil
		}
		return reflect.ValueOf(valI64).Convert(p).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		valUI64, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
		if kind == reflect.Uint64 {
			return valUI64, nil
		}
		return reflect.ValueOf(valUI64).Convert(p).Interface(), nil
	case reflect.Float64:
		return strconv.ParseFloat(str, 10)
	case reflect.Float32:
		valF64, err := strconv.ParseFloat(str, 10)
		if err != nil {
			return nil, err
		}
		return float32(valF64), nil
	case reflect.Complex64:
		return nil, errors.New("complex value are not supported")
	default:
		// reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer, reflect.Uintptr:
		return nil, errors.New("unsupported type")
	}
}
