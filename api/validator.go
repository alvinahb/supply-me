package api

import (
	"github.com/alvinahb/supply-me/util"
	"github.com/go-playground/validator/v10"
)

var validProductCategory validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if category, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedProductCategory(category)
	}

	return false
}

// TODO: Add product category column and validator
