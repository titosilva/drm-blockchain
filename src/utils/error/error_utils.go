package errorutils

import (
	"errors"
	"fmt"
)

func NewfWithInner(inner error, format string, vals ...any) error {
	return errors.Join(Newf(format, vals...), inner)
}

func NewWithInner(inner error, msg string) error {
	return errors.Join(New(msg), inner)
}

func Newf(format string, vals ...any) error {
	msg := fmt.Sprintf(format, vals...)
	return New(msg)
}

func New(msg string) error {
	return errors.New(msg)
}
