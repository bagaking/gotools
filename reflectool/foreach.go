package reflectool

import (
	"fmt"
	"reflect"

	"github.com/bagaking/gotools/procast"
)

type (
	forEachFieldConfig struct {
		onlyExported bool
		drillLevel   int
		// runtime val for drill
		visited map[reflect.Value]bool
		prefix  string
	}
	forEachFieldOption func(conf *forEachFieldConfig)

	FieldContext struct {
		Value *reflect.Value
		reflect.StructField
		Path string
	}

	FieldHandler func(ctx FieldContext) error
)

var ForEachFieldOptions = &forEachFieldConfig{}

func (*forEachFieldConfig) pipe(options []forEachFieldOption) *forEachFieldConfig {
	conf := &forEachFieldConfig{}
	for _, fn := range options {
		fn(conf)
	}
	return conf
}

func (*forEachFieldConfig) OnlyExported() forEachFieldOption {
	return func(conf *forEachFieldConfig) { conf.onlyExported = true }
}

func (*forEachFieldConfig) Drill(level int) forEachFieldOption {
	return func(conf *forEachFieldConfig) { conf.drillLevel = level }
}

func (*forEachFieldConfig) Override(fn func(cfg *forEachFieldConfig)) forEachFieldOption {
	return func(conf *forEachFieldConfig) { fn(conf) }
}

func ForEachField(target interface{}, fn FieldHandler, options ...forEachFieldOption) (err error) {
	defer procast.Recover(func(e error) { err = e }, "foreach execute failed")

	conf := (&forEachFieldConfig{}).pipe(options)

	var r reflect.Value
	if v, ok := target.(reflect.Value); ok {
		r = v
	} else {
		r = reflect.ValueOf(target)
	}

	elem := r
	if r.Kind() == reflect.Ptr {
		elem = r.Elem()
	}
	rType := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		field, fieldType := elem.Field(i), rType.Field(i)

		if conf.onlyExported && !IsFieldExported(fieldType) {
			continue
		}

		fieldCtx := FieldContext{
			StructField: fieldType,
			Value:       &field,
			Path:        fmt.Sprintf("%s.%s", conf.prefix, fieldType.Name),
		}
		// even the field is a structure to drill, the fn will be performed
		if err := fn(fieldCtx); err != nil {
			return err
		}

		if conf.drillLevel != 0 {
			if field.Kind() == reflect.Ptr {
				field = field.Elem()
			}
			if field.Kind() == reflect.Struct {
				if err := ForEachField(field, fn, conf.Override(func(cfg *forEachFieldConfig) {
					cfg.onlyExported = conf.onlyExported
					cfg.drillLevel = conf.drillLevel - 1
					cfg.visited = conf.visited
					cfg.prefix = fieldCtx.Path
				})); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
