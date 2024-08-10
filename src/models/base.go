package models

import (
	"database/sql/driver"
	"errors"
)

type UserGettable interface {
	GetUser() *User
}

type UserSettable interface {
	SetUser(*User)
}

func (modelID ModelID) Value() (driver.Value, error) {
	result := int64(modelID)
	return result, nil
}

func (modelID *ModelID) Scan(value interface{}) error {
	if value == nil {
		*modelID = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		*modelID = ModelID(v)
		return nil
	case int32:
		*modelID = ModelID(v)
		return nil
	case uint:
		*modelID = ModelID(v)
		return nil
	case uint64:
		*modelID = ModelID(v)
		return nil
	default:
		return errors.New("unsupported type for ModelID")
	}
}
