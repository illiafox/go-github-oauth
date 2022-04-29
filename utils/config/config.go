package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"

	"github.com/BurntSushi/toml"
)

func Read(filename string, env bool) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf Config

	err = toml.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	if env {
		err = ReadEnv(&conf)
		if err != nil {
			return nil, fmt.Errorf("environment loading: %w", err)
		}
	}

	return &conf, nil
}

func ReadEnv(conf *Config) error {
	return loadEnv(conf)
}

func loadEnv(s any) error {
	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		typeField := reflectType.Field(i)

		value := reflectValue.Field(i)
		kind := value.Kind()

		if kind == reflect.Struct {
			err := loadEnv(reflectValue.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
		}

		env, ok := typeField.Tag.Lookup("env")
		if !ok || env == "" {
			continue
		}

		v := os.Getenv(env)
		if v == "" {
			continue
		}

		switch kind {

		case reflect.String:
			value.SetString(v)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf(
					"can't parse %s with type %s: %w",
					typeField.Name,
					typeField.Type,
					err,
				)
			}

			value.SetInt(num)
		default:
			return fmt.Errorf("can't set %s", reflectType)
		}

	}

	return nil
}
