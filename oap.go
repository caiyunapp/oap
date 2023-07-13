package oap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/philchia/agollo/v4"
	"gopkg.in/yaml.v3"
)

type UnmarshalFunc func([]byte, interface{}) error

var registryForUnmarshal = map[string]UnmarshalFunc{
	"json": json.Unmarshal,
	"yaml": yaml.Unmarshal,
}

// You can use custom unmarshal for strcut type filed.
// Predfined JSON&YAML.
func SetUnmarshalFunc(name string, f UnmarshalFunc) {
	registryForUnmarshal[name] = f
}

func Decode(ptr interface{}, client agollo.Client, keyOpts map[string][]agollo.OpOption) error {
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		tag := structField.Tag
		apolloRawKey := tag.Get("apollo")
		apolloKeyParts := strings.Split(apolloRawKey, ",")
		apolloKey := apolloKeyParts[0]

		apolloVal := client.GetString(apolloKey, keyOpts[apolloKey]...)
		val := reflect.New(structField.Type)

		// use unmarshaller function
		if len(apolloKeyParts) > 1 {
			if unmarshallerFunc, ok := registryForUnmarshal[apolloKeyParts[1]]; ok {
				if err := unmarshallerFunc([]byte(apolloVal), val.Interface()); err != nil {
					return fmt.Errorf("%s unmarshal %s error: %w", apolloKeyParts[1], apolloKey, err)
				}

				v.FieldByName(structField.Name).Set(val.Elem())
			}
		}

		if err := yaml.Unmarshal([]byte(apolloVal), val.Interface()); err != nil {
			return fmt.Errorf("unmarshal %s error: %w", apolloVal, err)
		}

		v.FieldByName(structField.Name).Set(val.Elem())
	}

	return nil
}
