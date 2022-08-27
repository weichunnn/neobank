package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/weichunnn/neobank/util"
)

// reflection to examine types at runtime
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835
	// used interface to cater for invalid types to only validate for strings
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
