package int32g_test

import (
	"drm-blockchain/src/math/algebra/group/int32g"
	"testing"
)

func Test_zero_plus_any__should__any(t *testing.T) {
	var v int32 = 2
	if int32g.Zero().CombineWith(int32g.From(v)) != int32g.From(v) {
		t.Errorf("Zero combined with From(%d) should equals From(%d)", v, v)
	}
}
