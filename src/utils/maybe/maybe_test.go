package maybe_test

import (
	. "drm-blockchain/src/utils/maybe"
	"testing"
)

func Test__GetNothing__Should__Error(t *testing.T) {
	a := Nothing[int]()
	_, err := a.Get()

	if err == nil {
		t.Error("Nothing should not be able to get")
	}
}

func Test__GetJust__Should__ReturnValue(t *testing.T) {
	a := Just(3)
	val, err := a.Get()

	if err != nil {
		t.Error("Just should be able to get")
	}

	if val != 3 {
		t.Error("Unexpected value on Just Get")
	}
}

func Test__IsNothingNothing__Should__True(t *testing.T) {
	a := Nothing[int]()

	if !a.IsNothing() {
		t.Error("Just should return false on IsNothing")
	}
}

func Test__IsNothingJust__Should__False(t *testing.T) {
	a := Just(3)

	if a.IsNothing() {
		t.Error("Just should return false on IsNothing")
	}
}
