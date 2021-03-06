package annotation

import (
	"reflect"

	"github.com/bagaking/gotools/reflectool"
	"github.com/bagaking/gotools/strs"
)

type (
	IAnnotation interface {
		TagName() string
	}

	AnnMap map[string]IAnnotation

	StructAnnotations struct {
		Spawners          map[string]reflectool.Spawner
		Field2Annotations map[string]AnnMap // field => annotation_name => annotation
		Setter            func(val interface{}, desc IAnnotation) (interface{}, error)
	}
)

func ExtractFromTag(annotate IAnnotation, tag string) error {
	_, err := strs.KVStr(tag).ReflectTo(annotate)
	return err
}

func (sa StructAnnotations) Spawn(annotation IAnnotation) IAnnotation {
	return sa.Spawners[annotation.TagName()].Spawn().(IAnnotation)
}

func (sa StructAnnotations) Get(fieldName string, tagName string) IAnnotation {
	if sa.Field2Annotations[fieldName] == nil {
		return nil
	}

	return sa.Field2Annotations[fieldName][tagName]
}

func Analyze(prototype interface{}, annotations ...IAnnotation) (*StructAnnotations, error) {
	ret := &StructAnnotations{
		Spawners:          make(map[string]reflectool.Spawner),
		Field2Annotations: make(map[string]AnnMap),
	}
	for _, annotation := range annotations {
		ret.Spawners[annotation.TagName()] = reflectool.NewSpawner(annotation)
	}

	if err := reflectool.ForEachField(prototype, func(field *reflect.Value, fieldType reflect.StructField) error {
		for _, annotation := range annotations {
			tagName := annotation.TagName()
			tagContent, ok := fieldType.Tag.Lookup(tagName)
			if !ok {
				return nil
			}
			anData := ret.Spawn(annotation)
			if err := ExtractFromTag(anData, tagContent); err != nil {
				return err
			}
			if ret.Field2Annotations[fieldType.Name] == nil {
				ret.Field2Annotations[fieldType.Name] = make(AnnMap)
			}
			ret.Field2Annotations[fieldType.Name][tagName] = anData
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}