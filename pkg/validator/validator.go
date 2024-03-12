package validator

import (
	"reflect"
	"smm/pkg/constants"
	"smm/pkg/response"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func InitValidator() {
	if validate == nil {
		validate = validator.New()
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			jsonParts := 2
			name := strings.SplitN(fld.Tag.Get("json"), ",", jsonParts)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

func Validate(data interface{}) error {
	errs := validate.Struct(data)
	var errStr string
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			errStr = errStr + err.Field() + ": " + err.ActualTag()
			if err.Param() != "" {
				errStr = errStr + " = " + err.Param() + " | "
			} else {
				errStr += " | "
			}
		}
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Code: constants.ErrCodeAppBadRequest, Err: errStr})
	}
	return nil
}
