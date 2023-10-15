package bigintmodg_test

import (
	"drm-blockchain/src/math/algebra/group/bigintmodg"
	"testing"
)

func Test_sum_classes_mod(t *testing.T) {
	a := bigintmodg.FromInt64(3, 5)
	b := bigintmodg.FromInt64(4, 5)
	r := a.CombineWith(b)

	if !r.EqualsTo(bigintmodg.FromInt64(2, 5)) {
		t.Error("Comparison failed")
	}
}
