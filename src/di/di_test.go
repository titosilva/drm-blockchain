package di_test

import (
	"drm-blockchain/src/di"
	"testing"
)

type STest struct {
	content int
}

type STest2 struct {
	content int
}

func Test__DISingleton__ShouldReturn__SameInstanceAlways(t *testing.T) {
	ctx := new(di.DIContext)
	di.AddSingleton[STest](ctx, new(STest))

	test := di.GetService[STest](ctx)
	test.content = 212031

	test2 := di.GetService[STest](ctx)

	if test.content != test2.content {
		t.Error("Different instances provided for STest")
	}
}

func Test__DIFactory__ShouldReturn__NewInstanceEveryTime(t *testing.T) {
	ctx := new(di.DIContext)
	di.AddFactory[STest2](ctx, func() *STest2 {
		return new(STest2)
	})

	st2 := di.GetService[STest2](ctx)
	st2.content = 1000

	st3 := di.GetService[STest2](ctx)
	st3.content = 2000

	if st2.content == st3.content {
		t.Error("Same instance provided for STest2")
	}

	if st2.content != 1000 {
		t.Error("st2 content has changed")
	}
}
