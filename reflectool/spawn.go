package reflectool

import "reflect"

type (
	Spawner func() interface{}
)

func NewSpawner(model interface{}) (spawner Spawner) {
	return NewSpawnerFromType(reflect.TypeOf(model))
}

func NewSpawnerFromType(ty reflect.Type) (spawner Spawner) {
	if ty.Kind() == reflect.Ptr {
		return func() interface{} { return reflect.New(ty.Elem()).Interface() } // todo: .Addr() ? test this
	} else {
		return func() interface{} { return reflect.New(ty).Interface() }
	}
}

func (sp Spawner) Spawn() interface{} {
	return sp()
}
