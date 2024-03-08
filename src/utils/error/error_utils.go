package errorutils

import (
	"errors"
	"fmt"
)

func NewfWithInner(inner error, format string, vals ...any) error {
	return errors.Join(Newf(format, vals...), inner)
}

func Newf(format string, vals ...any) error {
	str := fmt.Sprintf(format, vals...)
	return errors.New(str)
}
