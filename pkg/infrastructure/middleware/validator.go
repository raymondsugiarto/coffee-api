package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response"
	"github.com/raymondsugiarto/coffee-api/pkg/shared/response/status"
)

type XValidator struct {
	validator *validator.Validate
}

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

// This is the validator instance
// for more information see: https://github.com/go-playground/validator
var validate = validator.New()

func (v XValidator) Validate(data interface{}) error {
	var err error
	validationErrors := []error{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem response.ErrorResponse

			elem.Field = err.Field()   // Export struct field name
			elem.Tag = err.Tag()       // Export struct tag
			elem.Value = err.Value()   // Export field value
			elem.Param = err.Param()   // Export struct tag parameter (if exists) if any
			elem.Message = err.Error() // Export validator error message

			validationErrors = append(validationErrors, elem)
		}
		err = status.New(status.BadRequest, validationErrors...)
	}

	return err
}

var (
	AppValidator *XValidator
)

func SetupValidator() {
	AppValidator = &XValidator{
		validator: validate,
	}
}
