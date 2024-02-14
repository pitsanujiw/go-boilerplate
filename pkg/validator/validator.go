package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	// v is the validator object used to validate incoming data.
	// From their docs:
	// Validate is designed to be thread-safe and used as a singleton instance.
	// It caches information about your struct and validations,
	// in essence only parsing your validation tags once per struct type.
	// Using multiple instances neglects the benefit of caching.
	v          *validator.Validate
	translator ut.Translator

	once sync.Once
)

func Struct(s any) error {
	once.Do(initValidator)

	var (
		errValidationErrors validator.ValidationErrors
		err                 = v.Struct(s)
	)

	switch {
	case errors.As(err, &errValidationErrors):

		// validator.ValidationErrors is a slice of errors, so we store it in a
		// slice and combine it into one error
		errMsgs := make([]string, len(errValidationErrors))

		for i, e := range errValidationErrors {
			errMsgs[i] = e.Translate(translator)
		}

		return errors.New(strings.Join(errMsgs, ", "))
	case err != nil:
		// this occurs if validation isn't properly set up or used
		return err
	default:
		return nil
	}
}

// initValidator constructs our global validator with default translator for friendlier error messages.
// Errors are not expected to be thrown at this stage.
func initValidator() {
	v = validator.New()
	translator, _ = ut.New(en.New()).GetTranslator("en")

	// if json tags are present we use that value instead, so it is more user-friendly
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// if tag key says it should be ignored, use field name
		if name == "-" {
			return fld.Name
		}
		return name
	})

	v.RegisterValidation("non_special_characters", func(fl validator.FieldLevel) bool {
		value := fl.Field().Interface().(string)
		return regexp.MustCompile(`^[a-zA-Z0-9-_ ]*$`).MatchString(value)
	})
}
