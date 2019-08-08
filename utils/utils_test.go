package utils

import "testing"
import "github.com/stretchr/testify/assert"

func TestInterect(t *testing.T) {
	l1 := []string{"f", "e", "d"}
	l2 := []string{"a", "b", "c"}
	l3 := IntersectSlices(l1, l2)
	assert.EqualValues(t, l3, []string(nil))

	l1 = []string{"a", "e", "d"}
	l2 = []string{"a", "b", "c"}
	l3 = IntersectSlices(l1, l2)
	assert.EqualValues(t, l3, []string{"a"})

	l1 = []string{"z", "a", "e", "d"}
	l2 = []string{"a", "b", "c"}
	l3 = IntersectSlices(l1, l2)
	assert.EqualValues(t, l3, []string{"a"})

	l1 = []string{"a"}
	l2 = []string{"a"}
	l3 = IntersectSlices(l1, l2)
	assert.EqualValues(t, l3, []string{"a"})

	l1 = []string{"a", "b", "c"}
	l2 = []string{"a", "b", "c"}
	l3 = IntersectSlices(l1, l2)
	assert.EqualValues(t, l3, []string{"a", "b", "c"})

	l1 = []string{"a", "b", "c"}
	l2 = []string{"c", "a", "b"}
	l3 = IntersectSlices(l1, l2)
	assert.EqualValues(t, l3, []string{"c", "a", "b"})
}
