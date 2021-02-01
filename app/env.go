package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type Env struct{}

func NewEnv() Env {
	return Env{}
}

func (e Env) Load(config interface{}) error {
	configVal := reflect.ValueOf(config)
	if configVal.Kind() != reflect.Ptr {
		return errors.New("config must be a pointer")
	}

	if configVal.IsNil() {
		return errors.New("config cannot be nil")
	}

	elem := configVal.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("config must be a struct")
	}

	numFields := elem.NumField()
	configType := elem.Type()

	for i := 0; i < numFields; i++ {
		field := configType.Field(i)

		envName, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}

		val := os.Getenv(envName)
		if val == "" {
			var ok bool
			val, ok = field.Tag.Lookup("default")
			if !ok {
				return fmt.Errorf("missing value for tag %s", field.Name)
			}
		}

		err := e.setFieldValue(field, elem.Field(i), val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e Env) setFieldValue(field reflect.StructField, value reflect.Value, rawValue string) error {
	kind := field.Type.Kind()
	switch kind {
	case reflect.String:
		value.SetString(rawValue)
	case reflect.Int, reflect.Int64:
		num, err := strconv.Atoi(rawValue)
		if err != nil {
			return err
		}
		value.SetInt(int64(num))
	case reflect.Bool:
		boolean, err := strconv.ParseBool(rawValue)
		if err != nil {
			return err
		}
		value.SetBool(boolean)
	default:
		return fmt.Errorf("unexpected field type: %s", kind)
	}
	return nil
}
