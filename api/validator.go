package api

import (
	"utils"

	"github.com/go-playground/validator/v10"
)

var validCoin validator.Func = func(fl validator.FieldLevel) bool {
	if coin, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedCoin(coin)
	}
	return false
}
