package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translationsEn "github.com/go-playground/validator/v10/translations/en"
	translationsRu "github.com/go-playground/validator/v10/translations/ru"
)

const (
	locale = `ru`
)

var (
	validate *validator.Validate

	universalTranslator *ut.UniversalTranslator

	translator ut.Translator

	customValidations = map[string]func(fl validator.FieldLevel) bool{
		// Override internal boolean validator with type check (for case `type SomeProperty bool`)
		`boolean`: func(fl validator.FieldLevel) bool {
			if fl.Field().Type().String() == `bool` {
				return true
			}
			value := fl.Field().String()
			_, err := strconv.ParseBool(value)
			return err == nil
		},
	}
)

func mustLoad() {
	if validate == nil {
		validate = validator.New()

		for tag, fn := range customValidations {
			if err := validate.RegisterValidation(tag, fn); err != nil {
				panic(err)
			}
		}

		validate.RegisterTagNameFunc(
			func(sf reflect.StructField) string {
				name := strings.SplitN(sf.Tag.Get("json"), ",", 2)[0]

				if name == "-" {
					return ""
				}

				return name
			},
		)

		var localeTranslator locales.Translator
		var registerTranslations func(v *validator.Validate, trans ut.Translator) (err error)
		switch locale {
		case `ru`:
			localeTranslator = ru.New()
			registerTranslations = translationsRu.RegisterDefaultTranslations
		case `en`:
			localeTranslator = en.New()
			registerTranslations = translationsEn.RegisterDefaultTranslations
		default:
			panic(fmt.Sprintf(`unknown locale "%s" for internal validator`, locale))
		}

		universalTranslator = ut.New(localeTranslator, localeTranslator)
		translator, _ = universalTranslator.GetTranslator(locale)
		if err := registerTranslations(validate, translator); err != nil {
			panic(fmt.Sprintf(`failed to register validator translations: %v`, err.Error()))
		}
	}
}

func IsInvalidValidationError(err error) bool {
	if err == nil {
		return false
	}
	_, is := err.(*validator.InvalidValidationError)
	return is
}

func AsValidationErrors(err error) validator.ValidationErrors {
	return err.(validator.ValidationErrors)
}

func AsCustomValidationTranslations(err error) validator.ValidationErrorsTranslations {
	return translateValidationErrors(AsValidationErrors(err), translator)
}

func translateValidationErrors(ve validator.ValidationErrors, ut ut.Translator) validator.ValidationErrorsTranslations {
	trans := make(validator.ValidationErrorsTranslations)

	var fe validator.FieldError

	for i := 0; i < len(ve); i++ {
		trans[fe.Field()] = ve[i].Translate(ut)
	}

	return trans
}

func Struct(s any) (validator.ValidationErrors, error) {
	mustLoad()
	err := validate.Struct(s)
	if err != nil {
		if IsInvalidValidationError(err) {
			return nil, err
		}
		return AsValidationErrors(err), nil
	}
	return nil, nil
}

func Var(field any, tag string) (validator.ValidationErrors, error) {
	mustLoad()
	err := validate.Var(field, tag)
	if err != nil {
		if IsInvalidValidationError(err) {
			return nil, err
		}
		return AsValidationErrors(err), nil
	}
	return nil, nil
}
