package csvp

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/bagaking/gotools/annotation"
	"github.com/bagaking/gotools/reflectool"
	"github.com/bagaking/gotools/strs"
)

type CSVAnnotation struct {
	Col    int
	Parser string
	Param  string
}

func (CSVAnnotation) TagName() string {
	return "csv"
}

var (
	once        = sync.Once{}
	csvAnStruct *annotation.StructAnnotations
	template    = CSVAnnotation{}
)

func ParseLine(data interface{}, line []string) (err error) {
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

		aCSV := a.(*CSVAnnotation)
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
		fmt.Println(fieldType.Name, value)

		return nil
	}); err != nil {
		return err
	}

	return nil
}
