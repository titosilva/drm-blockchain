package kdirectsumg_test

import (
	"drm-blockchain/src/math/algebra/group/bigintmodg"
	"drm-blockchain/src/math/algebra/group/kdirectsumg"
	"testing"
)

func Test__KDirectSumGCombine__Should__CombineAllEntriesCorrectly(t *testing.T) {
	kds1 := kdirectsumg.New[bigintmodg.BigIntModG](3)
	kds2 := kdirectsumg.New[bigintmodg.BigIntModG](3)

	iter1 := kds1.Entries.GetIterator()
	var agg int64 = 9
	for iter1.HasNext() {
		*iter1.GetNext() = bigintmodg.FromInt64(agg, 10)
		agg += 4
	}

	iter2 := kds2.Entries.GetIterator()
	agg = 3
	for iter2.HasNext() {
		*iter2.GetNext() = bigintmodg.FromInt64(agg, 10)
		agg += 3
	}

	comb := kds1.CombineWith(kds2)
	iter_comb := comb.AsPure().Entries.GetIterator()
	agg = 12
	for iter_comb.HasNext() {
		r := *iter_comb.GetNext()
		exp := bigintmodg.FromInt64(agg, 10)

		if !r.EqualsTo(exp) {
			t.Errorf("Expected %d mod 10, got %s mod %s", agg, r.AsPure().Value.String(), r.AsPure().Mod.String())
		}

		agg += 7
	}
}
