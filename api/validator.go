package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/weichunnn/neobank/util"
)

// reflection to examine types at runtime
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// https://github.com/gin-gonic/examples/blob/master/custom-validation/server.go
	// used interface to cater for invalid types to only validate for strings
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
