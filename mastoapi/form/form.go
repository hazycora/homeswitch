package form

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var (
	validate        = validator.New(validator.WithRequiredStructEnabled())
	reValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
)

func init() {
	validate.RegisterValidation("username", validateUsername)
}

func validateUsername(fl validator.FieldLevel) bool {
	return reValidUsername.MatchString(fl.Field().String())
}

type FormError struct {
	ErrorMessage string                       `json:"error"`
	Details      map[string][]FormErrorDetail `json:"details"`
}

func (e FormError) Error() string {
	return e.ErrorMessage
}

type FormErrorDetail struct {
	Error       string `json:"error"`
	Description string `json:"description"`
}

func ValidateForm(form interface{}) (err error) {
	err = validate.Struct(form)
	if err == nil {
		return
	}

	if _, invalid := err.(*validator.InvalidValidationError); invalid {
		log.Error().Err(err).Msg("Error validating form")
		return
	}

	val := reflect.ValueOf(form)
	formType := val.Type()

	formError := FormError{
		ErrorMessage: "Validation failed",
		Details:      make(map[string][]FormErrorDetail),
	}

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		structField, ok := formType.FieldByName(field)
		if !ok {
			log.Error().Str("field", field).Msg("Error validating form")
			continue
		}
		jsonName := structField.Tag.Get("json")

		tag := err.Tag()
		var mastoAPIErrorCode string
		var mastoAPIErrorDescription string

		switch tag {
		case "required":
			mastoAPIErrorCode = "ERR_BLANK"
			mastoAPIErrorDescription = "can't be blank"
		case "username":
			mastoAPIErrorCode = "ERR_INVALID"
			mastoAPIErrorDescription = "must contain only letters, numbers and underscores"
		case "max":
			mastoAPIErrorCode = "ERR_TOO_LONG"
			mastoAPIErrorDescription = "is too long (maximum is whatever characters)"
		}

		if _, ok := formError.Details[jsonName]; !ok {
			formError.Details[jsonName] = []FormErrorDetail{}
		}

		formError.Details[jsonName] = append(formError.Details[jsonName], FormErrorDetail{
			Error:       mastoAPIErrorCode,
			Description: mastoAPIErrorDescription,
		})
	}
	err = formError
	return
}
