package csvp

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"

	"github.com/bagaking/gotools/annotation"
	"github.com/bagaking/gotools/reflectool"
	"github.com/bagaking/gotools/strs"
)

type ColAnnotation struct {
	Col    int
	Parser string
	Param  string
}

func (ColAnnotation) TagName() string {
	return "csv"
}

var (
	csvAnStructCache sync.Map
	template         = ColAnnotation{}
)

func ParseByCol(outSlicePointer interface{}, reader LineReader) error {
	itr := func() (line interface{}, err error) { return reader.Read() }

	elemType, err := reflectool.GetSliceElementType(outSlicePointer)
	if err != nil {
		return fmt.Errorf("invalid input, %w", err)
	}

	elemSpawner := reflectool.NewSpawnerFromType(elemType)
	mapper := func(in interface{}) (interface{}, error) {
		line, ok := in.([]string)
		if !ok {
			return nil, fmt.Errorf("invalid input %v", in)
		}
		v := elemSpawner.Spawn()

		if e := ParseLineByCol(v, line); e != nil {
			return nil, e
		}
		return v, nil
	}
	csvReaderExitValidator := func(iv interface{}, err error) (bool, error) {
		if err == io.EOF {
			return true, nil
		}
		return false, err
	}

	return reflectool.Iterator(itr).WriteTo(outSlicePointer,
		reflectool.ItrMapper(mapper),
		reflectool.ItrExitValidator(csvReaderExitValidator),
	)
}

func ParseLineByCol(data interface{}, line []string) (err error) {
	csvAnStruct, err := csvAnnotationsFor(data)
	if err != nil {
		return err
	}

	if err = reflectool.ForEachField(data, func(fCtx reflectool.FieldContext) error {
		a := csvAnStruct.Get(fCtx.Path, template.TagName())
		if a == nil {
			return nil
		}

		aCSV := a.(*ColAnnotation)
		if aCSV.Col < 0 || aCSV.Col >= len(line) {
			return fmt.Errorf("csv column %d out of range for field %s, line has %d columns", aCSV.Col, fCtx.Path, len(line))
		}
		valStr := line[aCSV.Col]
		var value interface{}
		parser := aCSV.Parser
		switch {
		case "" == parser || strs.StartsWith(parser, "plain"):
			value, err = strs.Conv2PlainType(valStr, fCtx.Type)
		case strs.StartsWith(parser, "time"):
			if aCSV.Param == "" {
				value, err = time.Parse(time.RFC3339, valStr)
			} else {
				value, err = time.Parse(aCSV.Param, valStr)
			}
		}

		if err != nil {
			return err
		}
		fCtx.Value.Set(reflect.ValueOf(value))

		return nil
	},
		reflectool.ForEachFieldOptions.OnlyExported(),
		reflectool.ForEachFieldOptions.Drill(-1),
	); err != nil {
		return err
	}

	return nil
}

func csvAnnotationsFor(data interface{}) (*annotation.StructAnnotations, error) {
	ty := reflect.TypeOf(data)
	if ty == nil {
		return nil, fmt.Errorf("analyze model failed, data is nil")
	}
	for ty.Kind() == reflect.Ptr {
		ty = ty.Elem()
	}

	if cached, ok := csvAnStructCache.Load(ty); ok {
		return cached.(*annotation.StructAnnotations), nil
	}

	holder, err := annotation.Analyze(data, template)
	if err != nil {
		return nil, fmt.Errorf("%w, analyze model failed", err)
	}
	actual, _ := csvAnStructCache.LoadOrStore(ty, holder)
	return actual.(*annotation.StructAnnotations), nil
}
