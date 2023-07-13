package oap_test

import (
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/philchia/agollo/v4"
	"github.com/ringsaturn/oap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecode_String(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type CustomStringType string

	config := struct {
		String           string           `apollo:"foo"`
		CustomStringType CustomStringType `apollo:"bar"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("foo")).Return("{}")
	client.EXPECT().GetString(gomock.Eq("bar")).Return("hello")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)

	assert.Equal(t, "{}", config.String)
	assert.Equal(t, CustomStringType("hello"), config.CustomStringType)
}

func TestDecode_Int(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := struct {
		Int      int           `apollo:"int"`
		Int8     int8          `apollo:"int8"`
		Int16    int16         `apollo:"int16"`
		Int32    int32         `apollo:"int32"`
		Int64    int64         `apollo:"int64"`
		Uint     uint          `apollo:"uint"`
		Uint8    uint8         `apollo:"uint8"`
		Uint16   uint16        `apollo:"uint16"`
		Uint32   uint32        `apollo:"uint32"`
		Uint64   uint64        `apollo:"uint64"`
		Duration time.Duration `apollo:"duration"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("int")).Return("1")
	client.EXPECT().GetString(gomock.Eq("int8")).Return("1")
	client.EXPECT().GetString(gomock.Eq("int16")).Return("1")
	client.EXPECT().GetString(gomock.Eq("int32")).Return("1")
	client.EXPECT().GetString(gomock.Eq("int64")).Return("1")
	client.EXPECT().GetString(gomock.Eq("uint")).Return("1")
	client.EXPECT().GetString(gomock.Eq("uint8")).Return("1")
	client.EXPECT().GetString(gomock.Eq("uint16")).Return("1")
	client.EXPECT().GetString(gomock.Eq("uint32")).Return("1")
	client.EXPECT().GetString(gomock.Eq("uint64")).Return("1")
	client.EXPECT().GetString(gomock.Eq("duration")).Return("1m")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)

	assert.Equal(t, 1, config.Int)
	assert.Equal(t, int8(1), config.Int8)
	assert.Equal(t, int16(1), config.Int16)
	assert.Equal(t, int32(1), config.Int32)
	assert.Equal(t, int64(1), config.Int64)
	assert.Equal(t, uint(1), config.Uint)
	assert.Equal(t, uint8(1), config.Uint8)
	assert.Equal(t, uint16(1), config.Uint16)
	assert.Equal(t, uint32(1), config.Uint32)
	assert.Equal(t, uint64(1), config.Uint64)
	assert.Equal(t, time.Minute, config.Duration)
}

func TestDecode_Float(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := struct {
		Float32 float32 `apollo:"float32"`
		Float64 float64 `apollo:"float64"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("float32")).Return("1.1")
	client.EXPECT().GetString(gomock.Eq("float64")).Return("1.1")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)

	assert.Equal(t, float32(1.1), config.Float32)
	assert.Equal(t, float64(1.1), config.Float64)
}

func TestDecode_Bool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := struct {
		Bool bool `apollo:"bool"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("bool")).Return("true")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)

	assert.Equal(t, true, config.Bool)
}

func TestDecode_JSONStruct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	config := struct {
		User    user  `apollo:"user"`
		UserPtr *user `apollo:"user_ptr"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("user")).Return(`{"name":"Alice","age":18}`)
	client.EXPECT().GetString(gomock.Eq("user_ptr")).Return(`{"name":"Alice","age":18}`)

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)

	assert.Equal(t, user{Name: "Alice", Age: 18}, config.User)
	assert.Equal(t, &user{Name: "Alice", Age: 18}, config.UserPtr)
}

func TestDecode_YAMLStruct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type user struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	config := struct {
		User    user  `apollo:"user"`
		UserPtr *user `apollo:"user_ptr"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("user")).Return("name: Alice\nage: 18")
	client.EXPECT().GetString(gomock.Eq("user_ptr")).Return("name: Alice\nage: 18")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)

	assert.Equal(t, user{Name: "Alice", Age: 18}, config.User)
	assert.Equal(t, &user{Name: "Alice", Age: 18}, config.UserPtr)
}

func TestDecode_StructSlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type user struct {
		Name string `yaml:"name"`
		Age  int    `yaml:"age"`
	}

	t.Run("using yaml", func(t *testing.T) {
		config := struct {
			Users []user `apollo:"users,yaml"`
		}{}

		client := NewMockClient(ctrl)
		client.EXPECT().GetString(gomock.Eq("users")).Return("- name: Alice\n  age: 18")

		err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
		require.NoError(t, err)

		assert.Equal(t, []user{{Name: "Alice", Age: 18}}, config.Users)
	})

	t.Run("raw", func(t *testing.T) {
		config := struct {
			Users []user `apollo:"users"`
		}{}

		client := NewMockClient(ctrl)
		client.EXPECT().GetString(gomock.Eq("users")).Return("- name: Alice\n  age: 18")

		err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
		require.NoError(t, err)

		assert.Equal(t, []user{{Name: "Alice", Age: 18}}, config.Users)
	})
}

func TestDecode_NonTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := struct {
		Foo        string
		unexported interface{}
	}{}

	client := NewMockClient(ctrl)

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)
}

func TestDecode_WithNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := struct {
		Foo string `apollo:"foo" apollo_namespace:"ns"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("foo"), gomock.Any()).Return("bar")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)
}

func TestDecode_NestedWithNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type FooConfig struct {
		Foo string `apollo:"foo"`
	}

	config := struct {
		FooConfig `apollo_namespace:"ns"`
	}{}

	client := NewMockClient(ctrl)
	client.EXPECT().GetString(gomock.Eq("foo"), gomock.Any()).Return("bar")

	err := oap.Decode(&config, client, make(map[string][]agollo.OpOption))
	require.NoError(t, err)
}

func ExampleDecode() {
	const (
		appid = "SampleApp"
		addr  = "http://81.68.181.139:8080"
		// http://81.68.181.139/app/access_key.html?#/appid=SampleApp
		secret = ""
	)

	if secret == "" {
		return
	}

	type FooConfig struct {
		Foo string `apollo:"abc"`
	}

	config := struct {
		FooConfig `apollo_namespace:"proper"`
	}{}

	client := agollo.NewClient(&agollo.Conf{
		AppID:           appid,
		MetaAddr:        addr,
		AccesskeySecret: secret,
	})
	if err := client.Start(); err != nil {
		panic(err)
	}

	if err := client.SubscribeToNamespaces("proper"); err != nil {
		panic(err)
	}

	fmt.Println(client.GetAllKeys(agollo.WithNamespace("application")))
	fmt.Println(client.GetAllKeys(agollo.WithNamespace("proper")))

	if err := oap.Decode(&config, client, make(map[string][]agollo.OpOption)); err != nil {
		panic(err)
	}

	fmt.Println(config.Foo)
}
