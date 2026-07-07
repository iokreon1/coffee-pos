package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	v := validator.New()

	// Daftarkan fungsi untuk ambil nama field dari JSON tag
	// supaya error key pakai "product_name" bukan "ProductName"
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) map[string]string {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = humanizeError(e)
	}
	return errors
}

func humanizeError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		if e.Kind().String() == "string" {
			return fmt.Sprintf("minimal %s karakter", e.Param())
		}
		return fmt.Sprintf("minimal %s", e.Param())
	case "max":
		if e.Kind().String() == "string" {
			return fmt.Sprintf("maksimal %s karakter", e.Param())
		}
		return fmt.Sprintf("maksimal %s", e.Param())
	case "oneof":
		return fmt.Sprintf("harus salah satu dari: %s", e.Param())
	case "uuid4":
		return "format UUID tidak valid"
	case "len":
		return fmt.Sprintf("harus tepat %s karakter", e.Param())
	case "numeric":
		return "harus berupa angka"
	default:
		return "tidak valid"
	}
}
