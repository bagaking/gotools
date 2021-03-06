package reflectool

import (
	"reflect"

	"github.com/bagaking/gotools/procast"
)

type (
	forEachFieldConfig struct {
		onlyExported bool
	}
	forEachFieldOption func(conf *forEachFieldConfig)

	FieldHandler func(field *reflect.Value, fieldType reflect.StructField) error
)

var ForEachFieldOptions = &forEachFieldConfig{}

func (conf *forEachFieldConfig) pipe(options []forEachFieldOption) *forEachFieldConfig {
	for _, fn := range options {
		fn(conf)
	}
	return conf
}

func (conf *forEachFieldConfig) OnlyExported() forEachFieldOption {
	return func(conf *forEachFieldConfig) { conf.onlyExported = true }
}

func ForEachField(target interface{}, fn FieldHandler, options ...forEachFieldOption) (err error) {
	defer procast.Recover(func(e error) { err = e }, "foreach execute failed")

	conf := (&forEachFieldConfig{}).pipe(options)

	r := reflect.ValueOf(target)
	elem := r
	if r.Kind() == reflect.Ptr {
		elem = r.Elem()
	}
	rType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field, fieldType := elem.Field(i), rType.Field(i)

		if conf.onlyExported && IsFieldExported(fieldType) {
			continue
		}

		if err := fn(&field, fieldType); err != nil {
			return err
		}
	}
	return nil
}
