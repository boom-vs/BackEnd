package internal

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func InsecureMarshal(data interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	newEncoder := json.NewEncoder(buffer)
	newEncoder.SetEscapeHTML(false)
	err := newEncoder.Encode(data)
	if err != nil {
		return []byte{}, err
	}
	return []byte(strings.ReplaceAll(buffer.String(), "\n", "")), nil
}

func GetHash(date []byte) string {
	return fmt.Sprintf("%X", sha1.Sum(date))
}

func StructToMap(obj interface{}) map[string]interface{} {
	objValue := reflect.ValueOf(obj)
	objType := objValue.Type()

	data := make(map[string]interface{})

	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldName := objType.Field(i).Name
		data[fieldName] = field.Interface()
	}

	return data
}

func Copier(source, destination interface{}) error {
	//Exit if destination not ptr becuase can't write to struct
	if reflect.TypeOf(destination).Kind() != reflect.Ptr {
		return errors.New(
			fmt.Sprintf("destination variable does not contain a reference to an objec %s",
				reflect.TypeOf(destination).Name()))
	}

	typeOfSource := reflect.TypeOf(source)
	valueOfSource := reflect.ValueOf(source)

	for typeOfSource.Kind() == reflect.Ptr {
		if valueOfSource.IsNil() {
			return errors.New("source variable is nil")
		}
		typeOfSource = typeOfSource.Elem()
		valueOfSource = valueOfSource.Elem()
	}

	typeOfDestination := reflect.TypeOf(destination)
	valueOfDestination := reflect.ValueOf(destination)
	for typeOfDestination.Kind() == reflect.Ptr {
		typeOfDestination = typeOfDestination.Elem()
		valueOfDestination = valueOfDestination.Elem()
	}

	for i := 0; i < typeOfSource.NumField(); i++ {
		fieldTypeOfSource := typeOfSource.Field(i)
		fieldValueOfSource := valueOfSource.Field(i)

		//If ptr nil
		if fieldTypeOfSource.Type.Kind() == reflect.Ptr && fieldValueOfSource.IsNil() {
			continue
		}

		//Excluded for gorm
		if fieldTypeOfSource.Type.String() == "gorm.Model" {
			err := Copier(fieldValueOfSource.Interface(), destination)
			if err != nil {
				return err
			}
			continue
		}

		foundField, exist := typeOfDestination.FieldByName(fieldTypeOfSource.Name)
		if !exist {
			continue
		}

		if foundField.Type == fieldTypeOfSource.Type {
			valueOfDestination.FieldByName(fieldTypeOfSource.Name).Set(fieldValueOfSource)
			continue
		}

		//For slices
		if fieldTypeOfSource.Type.Kind() == foundField.Type.Kind() && foundField.Type.Kind() == reflect.Slice {
			sliceTypeOfSource := fieldTypeOfSource.Type.Elem()
			sliceTypeOfDestination := foundField.Type.Elem()

			for sliceTypeOfSource.Kind() == reflect.Ptr {
				sliceTypeOfSource = sliceTypeOfSource.Elem()
			}

			for sliceTypeOfDestination.Kind() == reflect.Ptr {
				sliceTypeOfDestination = sliceTypeOfDestination.Elem()
			}

			if typeOfSource != sliceTypeOfSource || typeOfDestination != sliceTypeOfDestination {
				continue
			}

			element := valueOfDestination.FieldByName(fieldTypeOfSource.Name)

			for n := 0; n < fieldValueOfSource.Len(); n++ {
				newElement := reflect.New(sliceTypeOfDestination).Elem()
				err := Copier(fieldValueOfSource.Index(n).Interface(), newElement.Addr().Interface())
				if err != nil {
					return err
				}
				element.Set(reflect.Append(element, newElement.Addr()))
			}
		}
	}

	return nil
}

func GormUpdateOrCreate(db *gorm.DB, values interface{}) *gorm.DB {
	var ID uint64 = reflect.ValueOf(values).Elem().FieldByName("ID").Uint()

	if ID == 0 {
		return db.Create(values)
	}
	return db.Updates(values)
}

func MapToStrcut(source, destination interface{}) error {
	mapToStructConfig := &mapstructure.DecoderConfig{
		ErrorUnset: true,
		Result:     &destination,
	}

	mapToStruct, err := mapstructure.NewDecoder(mapToStructConfig)
	if err != nil {
		return errors.New("MapToStruct " + err.Error())
	}

	err = mapToStruct.Decode(source)
	if err != nil {
		return errors.New("MapToStruct " + err.Error())
	}

	return nil
}
