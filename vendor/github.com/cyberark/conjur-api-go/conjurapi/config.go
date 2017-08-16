package conjurapi

import (
	"fmt"
	"reflect"
	"strings"
	"os"
)

type Config struct {
	Account        string `validate:"required" env:"CONJUR_ACCOUNT"`
	ApplianceURL   string `validate:"required" env:"CONJUR_APPLIANCE_URL"`
}

func (c *Config) validate() (error) {
	v := reflect.ValueOf(*c)
	errors := []string{}

	const tagName = "validate"
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		tag := f.Tag.Get(tagName)

		switch tag {
		case "required":
			val := v.Field(i).Interface()
			if val.(string) == "" {
				errors = append(errors, fmt.Sprintf("%s is required.", f.Name))
			}
		default:
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errors, "\n"))
}

func LoadFromEnv(c interface{}) {
	const tagName = "env"

	vElem := reflect.ValueOf(c).Elem()
	vType := vElem.Type()


	for i := 0; i < vElem.NumField(); i++ {
		typeField := vType.Field(i)
		elemField := vElem.Field(i)
		tag := typeField.Tag.Get(tagName)

		switch elemField.Interface().(type) {
		case string:
			elemField.SetString(os.Getenv(tag))
		default:
			continue
		}

	}
}
