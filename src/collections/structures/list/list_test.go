package list_test

import (
	"drm-blockchain/src/collections/structures/list"
	"drm-blockchain/src/utils/maybe"
	"testing"
)

func Test__ListWithElements__Should__IterateElements(t *testing.T) {
	expected := make([]int, 20)

	for i := 0; i < 20; i++ {
		expected[i] = i * i
	}

	l := list.New[int]()
	for _, elem := range expected {
		l.Add(elem)
	}

	iter := l.GetIterator()
	idx := 0
	for iter.HasNext() {
		if expected[idx] != *iter.GetNext() {
			t.Errorf("element %d of list did not match with the expected", idx)
		}
		idx++
	}
}

func Test__ListPredicates__Should__ReturnTrueIfSatisfied(t *testing.T) {
	l := list.New[int]()
	for i := 0; i < 20; i++ {
		l.Add(i * i)
	}

	even_pred := func(i int) bool { return i%2 == 0 }
	odd_pred := func(i int) bool { return i%2 != 0 }
	even := l.Where(even_pred)
	odd := l.Where(odd_pred)

	if len(even.ToArray()) != len(odd.ToArray()) {
		t.Error("expected lists with same size")
	}

	if even.Any(odd_pred) {
		t.Error("expected no odds in 'even'")
	}

	if odd.Any(even_pred) {
		t.Error("expected no even in 'odd'")
	}

	if !even.All(even_pred) {
		t.Error("expected no odds in 'even'")
	}

	if !odd.All(odd_pred) {
		t.Error("expected no even in 'odd'")
	}

	if odd.First() != maybe.Just(1) {
		t.Error("expected 1 to be the first odd")
	}

	if even.First() != maybe.Just(0) {
		t.Error("expected 1 to be the first odd")
	}
}

func Test__ComposedListQueries__Should__ReturnResultOfPredicatesApplied(t *testing.T) {
	l := list.New[int]()
	for i := 0; i < 20; i++ {
		l.Add(i)
	}

	even_pred := func(i int) bool { return i%2 == 0 }
	greaterThan10 := func(i int) bool { return i > 10 }
	lessThan15 := func(i int) bool { return i < 15 }

	result := l.Where(even_pred).Where(greaterThan10).Where(lessThan15)

	if result.Count() != 2 {
		t.Errorf("expected exactly 2 elements, got %d", result.Count())
	}

	iter := result.GetIterator()
	for iter.HasNext() {
		x := iter.GetNext()

		if *x != 12 && *x != 14 {
			t.Error("expected numbers to be 12 or 14, given the predicates provided")
		}
	}

	if !(result.All(even_pred) && result.All(greaterThan10) && result.All(lessThan15)) {
		t.Error("not all predicates match with the 'All' expected results")
	}
}

func Test__ListIterator__Should__AllowEditsInList(t *testing.T) {
	l := list.New[int]()
	for i := 0; i < 20; i++ {
		l.Add(i)
	}

	iter := l.GetIterator()
	for iter.HasNext() {
		x := iter.GetNext()

		if *x%2 == 0 {
			continue
		}

		*x = (*x) * (*x)
	}

	even_pred := func(i int) bool { return i%2 == 0 }
	greaterThan20 := func(i int) bool { return i > 20 }
	if l.Where(even_pred).Any(greaterThan20) {
		t.Error("expected only odds to have been squared")
	}

	odd_pred := func(i int) bool { return i%2 != 0 }
	iter_odd := l.Where(odd_pred).GetIterator()
	odd := 1
	for iter_odd.HasNext() {
		x := iter_odd.GetNext()

		if *x != odd*odd {
			t.Errorf("value %d was not the expected product %d * %d", *x, odd, odd)
		}

		odd += 2
	}
}
