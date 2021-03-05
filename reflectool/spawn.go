package reflectool

import "reflect"

type (
	Spawner func() interface{}
)

func NewSpawner(model interface{}) (spawner Spawner) {
	if ty := reflect.TypeOf(model); ty.Kind() == reflect.Ptr {
		return func() interface{} {
			return reflect.New(reflect.ValueOf(model).Elem().Type()).Interface()
		}
	} else {
		return func() interface{} { return reflect.New(ty).Interface() }
	}
}

func (sp Spawner) Spawn() interface{} {
	return sp()
}
