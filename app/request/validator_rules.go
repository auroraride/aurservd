package request

import (
	"github.com/auroraride/aurservd/app/model/types"
	"github.com/go-playground/validator/v10"
)

func validateEnum(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(types.Enum)
	return value.IsValid()
}
