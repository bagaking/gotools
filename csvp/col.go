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
	once        = sync.Once{}
	csvAnStruct *annotation.StructAnnotations
	template    = ColAnnotation{}
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
	if csvAnStruct == nil {
		holder, err := annotation.Analyze(data, template)
		if err != nil {
			return fmt.Errorf("%w, analyze model failed", err)
		}
		once.Do(func() {
			csvAnStruct = holder
		})
	}

	if err = reflectool.ForEachField(data, func(field *reflect.Value, fieldType reflect.StructField) error {
		a := csvAnStruct.Get(fieldType.Name, template.TagName())
		if a == nil {
			return nil
		}

		aCSV := a.(*ColAnnotation)
		valStr := line[aCSV.Col]
		var value interface{}
		parser := aCSV.Parser
		switch {
		case "" == parser || strs.StartsWith(parser, "plain"):
			value, err = strs.Conv2PlainType(valStr, fieldType.Type)
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
		field.Set(reflect.ValueOf(value))

		return nil
	}); err != nil {
		return err
	}

	return nil
}
