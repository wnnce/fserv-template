package tool

import (
	"errors"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gofiber/fiber/v2"
)

var defaultStructValidator *StructValidator

func init() {
	defaultStructValidator = NewStruckValidator()
}

type StructValidator struct {
	validate *validator.Validate
	trans    ut.Translator
}

func Validator() *StructValidator {
	return defaultStructValidator
}

func (self *StructValidator) Engine() any {
	return self.validate
}

func (self *StructValidator) Validate(out any) error {
	err := self.validate.Struct(out)
	if err == nil {
		return nil
	}
	var validateErrors validator.ValidationErrors
	if !errors.As(err, &validateErrors) || len(validateErrors) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Struct parameter error")
	}
	return fiber.NewError(fiber.StatusBadRequest, validateErrors[0].Translate(self.trans))
}

func NewStruckValidator() *StructValidator {
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	return &StructValidator{
		validate: validate,
		trans:    trans,
	}
}
