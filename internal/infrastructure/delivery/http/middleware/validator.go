package middleware

import (
	"net/http"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
	translator ut.Translator
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()
	
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")
	
	enTranslations.RegisterDefaultTranslations(validate, trans)

	return &CustomValidator{
		validator: validate,
		translator: trans,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)

	if err == nil {
		return nil
	}
	
	errs := err.(validator.ValidationErrors)
	return echo.NewHTTPError(http.StatusBadRequest, errs.Translate(cv.translator))
}