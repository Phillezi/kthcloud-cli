package validation

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var once sync.Once

func Validate(body any) error {
	once.Do(registerCustomValidators)

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("failed to access gin validator")
	}
	return v.Struct(body)
}
