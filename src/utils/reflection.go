package utils

import (
	"github.com/jinzhu/copier"
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

func MarshalModelToForm[Model any, Form any](model *Model) (*Form, error) {
	var form Form
	err := copier.Copy(&form, model)
	if err != nil {
		return nil, err
	}
	return &form, err
}
