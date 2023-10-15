package direct_product

import (
	"drm-blockchain/src/math/algebra/group"
	"drm-blockchain/src/math/algebra/group/int32g"
)

var _ group.Elem[DirectProductG[int32g.Int32G, int32g.Int32G]] = DirectProductG[int32g.Int32G, int32g.Int32G]{}

type DirectProductG[X any, Y any] struct {
	g group.Elem[X]
	h group.Elem[Y]
}

func From[G any, H any](g group.Elem[G], h group.Elem[H]) DirectProductG[G, H] {
	return DirectProductG[G, H]{g: g, h: h}
}

func (x DirectProductG[X, Y]) CombineWith(y DirectProductG[X, Y]) DirectProductG[X, Y] {
	return DirectProductG[X, Y]{
		g: x.g.CombineWith(y.g).AsGroup(),
		h: x.h.CombineWith(y.h).AsGroup(),
	}
}
