package lib

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	uni "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var gpgValidator *validator.Validate

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type CustomValidator struct {
	validator *validator.Validate
}

func (c CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return err
	}
	return nil
}
func MapError(valErr error) (errs *[]ValidationError) {
	if gpgValidator == nil {
		log.Fatal("Instantiate validator first")
	}
	english := en.New()
	universalTranslator := uni.New(english, english)
	trans, _ := universalTranslator.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(gpgValidator, trans)
	if valErr == nil {
		return &[]ValidationError{}
	}
	validationErr := valErr.(validator.ValidationErrors)
	errors := []ValidationError{}
	for _, e := range validationErr {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errors = append(errors, ValidationError{Field: e.Field(), Message: translatedErr.Error()})

	}
	return &errors
}
func NewCustomValidator() *CustomValidator {
	gpgValidator = validator.New()
	gpgValidator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return &CustomValidator{gpgValidator}
}
