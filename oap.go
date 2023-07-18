package oap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/philchia/agollo/v4"
	"gopkg.in/yaml.v3"
)

const (
	apolloTag          = "apollo"
	apolloNamespaceTag = "apollo_namespace"
)

type UnmarshalFunc func([]byte, interface{}) error

var registryForUnmarshal = map[string]UnmarshalFunc{
	"json": json.Unmarshal,
	"yaml": yaml.Unmarshal,
}

// You can use custom unmarshal for struct type filed.
// Package oap provides built-in support for JSON&YAML.
func SetUnmarshalFunc(name string, f UnmarshalFunc) {
	registryForUnmarshal[name] = f
}

func Decode(ptr interface{}, client agollo.Client, keyOpts map[string][]agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	return decodeStruct(ptr, client, nil, keyOpts)
}

func decodeStruct(ptr interface{}, client agollo.Client, opts []agollo.OpOption, keyOpts map[string][]agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		field := v.FieldByName(structField.Name)

		if err := decode(ptr, structField, field, client, opts, keyOpts); err != nil {
			return err
		}
	}

	return nil
}

func decode(ptr interface{}, structField reflect.StructField, field reflect.Value, client agollo.Client, opts []agollo.OpOption, keyOpts map[string][]agollo.OpOption) error {
	if !field.CanSet() {
		return nil
	}

	tag := structField.Tag
	apolloRawKey := tag.Get(apolloTag)
	apolloKeyParts := strings.Split(apolloRawKey, ",")
	apolloKey := apolloKeyParts[0]

	// OpOptions
	kopts := keyOpts[apolloKey]
	newOpts := make([]agollo.OpOption, len(opts)+len(kopts))

	copy(newOpts, opts)
	copy(newOpts[len(opts):], kopts)

	// using namespace
	if ns := tag.Get(apolloNamespaceTag); ns != "" {
		//nolint:makezero // we have already copy kopts and opts to newOpts
		newOpts = append(newOpts, agollo.WithNamespace(ns))
	}

	val := reflect.New(structField.Type)

	// nested struct fields
	if apolloKey == "" {
		if err := decodeStruct(val.Interface(), client, newOpts, keyOpts); err != nil {
			return fmt.Errorf("Decode %s error: %w", structField.Name, err)
		}

		field.Set(val.Elem())

		return nil
	}

	// get config content
	apolloVal := client.GetString(apolloKey, newOpts...)

	// set raw string, first
	if field.Kind() == reflect.String {
		field.Set(reflect.ValueOf(apolloVal).Convert(val.Elem().Type()))

		return nil
	}

	// use unmarshaller function
	if len(apolloKeyParts) > 1 {
		if unmarshallerFunc, ok := registryForUnmarshal[apolloKeyParts[1]]; ok {
			if err := unmarshallerFunc([]byte(apolloVal), val.Interface()); err != nil {
				return fmt.Errorf("%s unmarshal %s error: %w", apolloKeyParts[1], apolloKey, err)
			}

			if field.CanSet() {
				field.Set(val.Elem())
			}

			return nil
		}
	}

	// parse value via yaml
	if err := yaml.Unmarshal([]byte(apolloVal), val.Interface()); err != nil {
		return fmt.Errorf("unmarshal %s error: %w", apolloVal, err)
	}

	field.Set(val.Elem())

	return nil
}
