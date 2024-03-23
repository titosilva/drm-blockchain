package ez

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type EzTest struct { // yes, you really read "Easy Test"
	t *testing.T
}

func New(t *testing.T) *EzTest {
	e := new(EzTest)
	e.t = t
	return e
}

func (ez *EzTest) AssertFalse(condition bool, msg ...string) {
	ez.assertWithLevel(1, !condition, msg...)
}

func (ez *EzTest) AssertAreEqual(val1 any, val2 any, msg ...string) {
	ez.assertWithLevel(1, reflect.DeepEqual(val1, val2), choose("Expected values to be equal, but got different values", msg))
}

func (ez *EzTest) AssertNoError(err error, msg ...string) {
	if err != nil {
		defMsg := fmt.Sprintf("Expected no error here, but got an error: %s", err.Error())
		ez.failWithMsg(1, defMsg, msg)
	}
}

func (ez *EzTest) Assert(condition bool, msg ...string) {
	if !condition {
		ez.failWithMsg(1, "Expected condition to be true, but got false", msg)
	}
}

func (ez *EzTest) assertWithLevel(callLevel int, condition bool, msg ...string) {
	if !condition {
		ez.failWithMsg(callLevel+1, "Expected condition to be true, but got false", msg)
	}
}

func (ez *EzTest) failWithMsg(callLevel int, defaultMsg string, alternativeMsg []string) {
	msg := choose(defaultMsg, alternativeMsg)

	_, absPath, line, ok := runtime.Caller(callLevel + 1)
	if ok {
		_, filename := path.Split(absPath)
		msg = fmt.Sprintf("%s:%d: %s", filename, line, msg)
	}

	ez.t.Fatal(msg)
}

func choose(defaultMsg string, alternativeMsg []string) string {
	if len(alternativeMsg) > 0 {
		return strings.Join(alternativeMsg, "; ")
	}

	return defaultMsg
}
