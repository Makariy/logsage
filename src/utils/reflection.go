package utils

import (
	"reflect"
)

func SetField(model interface{}, fieldName string, value interface{}) {
	dst := reflect.ValueOf(model).Elem()
	fieldValue := reflect.ValueOf(value)

	field := dst.FieldByName(fieldName)
	if field.IsValid() && field.CanSet() {
		field.Set(fieldValue)
	}
}
