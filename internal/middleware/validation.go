package middleware

import (
	"github.com/garaevmir/avitocoinstore/internal/model"
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("valid_item", func(fl validator.FieldLevel) bool {
		item := fl.Field().String()
		_, exists := model.Items[item]
		return exists
	})
	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
