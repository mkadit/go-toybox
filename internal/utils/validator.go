package utils

import (
	"context"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New(validator.WithRequiredStructEnabled())

var dynamicStructValidator = map[string]map[string]string{
	"ExampleValidation": {
		"idCustomer": "required",
	},
}

// RegisterValidation Register dynamic validator
func RegisterValidation(validatorMapping map[string]map[string]string) {
	for keys := range validatorMapping {
		Validator.RegisterValidationCtx(keys, func(ctx context.Context, fl validator.FieldLevel) bool {
			m, ok := fl.Field().Interface().(map[string]interface{})
			if !ok {
				return false
			}
			vm := ctx.Value(keys).(map[string]string)

			for k, v := range vm {
				if Validator.Var(m[k], v) != nil {
					return false
				}
			}
			return true
		})
	}
}

// ValidateRequest Validate a request and check if value is empty
func ValidateRequest(request any) bool {
	v := reflect.ValueOf(request)
	t := reflect.TypeOf(request)
	stringType := reflect.TypeOf("")
	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		valueType := value.Type()
		typeOf := t.Field(i)
		if typeOf.Name == "Description" {
			continue
		}
		if valueType == stringType && value.Interface() == "" {
			return false
		}
	}
	return true
}

func Reduce[T, U any](slice []T, reducer func(U, T) U, initialValue U) U {
	result := initialValue
	for _, v := range slice {
		result = reducer(result, v)
	}
	return result
}
