package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/anggadarkprince/crud-employee-go/utilities"
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Validator *validator.Validate
var Trans ut.Translator

func Init() {
	eng := english.New()
	uni := ut.New(eng, eng)
	var found bool
    Trans, found = uni.GetTranslator("en")
	if !found {
		fmt.Println("Translation not found")
	}

	Validator = validator.New(validator.WithRequiredStructEnabled())

	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })

	err := en_translations.RegisterDefaultTranslations(Validator, Trans)
	if err != nil {
		fmt.Println("Register translation error")
	}

	Validator.RegisterValidation("gender", func(fl validator.FieldLevel) bool {
		g := fl.Field().String()
		return g == "Male" || g == "Female"
	})

	Validator.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		username := fl.Field().String()
		pattern := `^[A-Za-z0-9._-]+$`
		matched, _ := regexp.MatchString(pattern, username)
		return matched
	})
}

func FormatValidationErrors(err error) map[string]string {
    errors := make(map[string]string)
    
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrors {
            errors[e.Field()] = utilities.Capitalize(e.Translate(Trans))
        }
    }
    
    return errors
}