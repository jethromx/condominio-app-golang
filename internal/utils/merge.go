package utils

import (
	"errors"
	"reflect"
)

// MergeStructs realiza un merge de los datos entrantes con los datos existentes
func MergeStructs(existing, updates interface{}) error {
	if existing == nil || updates == nil {
		return errors.New("existing and updates cannot be nil")
	}

	existingVal := reflect.ValueOf(existing)
	updatesVal := reflect.ValueOf(updates)

	if existingVal.Kind() != reflect.Ptr || updatesVal.Kind() != reflect.Ptr {
		return errors.New("existing and updates must be pointers")
	}

	existingVal = existingVal.Elem()
	updatesVal = updatesVal.Elem()

	if existingVal.Kind() != reflect.Struct || updatesVal.Kind() != reflect.Struct {
		return errors.New("existing and updates must be pointers to structs")
	}

	for i := 0; i < updatesVal.NumField(); i++ {
		field := updatesVal.Field(i)
		if !field.IsZero() {
			existingField := existingVal.Field(i)
			if existingField.CanSet() {
				existingField.Set(field)
			}
		}
	}
	return nil
}
