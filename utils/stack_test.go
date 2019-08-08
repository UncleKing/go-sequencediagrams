package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStack_Push(t *testing.T) {
	a := 1
	b := 2
	c := 3
	d := 4

	s := Stack{}
	s.Push(&a)
	assert.Equal(t, 1, s.Count(), "Count should be 1")

	as, ok := s.Peek().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 1, *as, "Peek should return 1")

	s.Push(&b)
	assert.Equal(t, 2, s.Count(), "Count should be 2")

	as, ok = s.Peek().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 2, *as, "Peek should return 2")

	s.Push(&c)
	assert.Equal(t, 3, s.Count(), "Count should be 3")

	as, ok = s.Peek().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 3, *as, "Peek should return 3")

	s.Push(&d)
	assert.Equal(t, 4, s.Count(), "Count should be 4")

	as, ok = s.Peek().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 4, *as, "Peek should return 4")

	as, ok = s.Pop().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 4, *as, "Peek should return 4")
	assert.Equal(t, &d, as, "pointers should be same.")

	as, ok = s.Pop().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 3, *as, "Peek should return 4")
	assert.Equal(t, &c, as, "pointers should be same.")

	as, ok = s.Pop().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 2, *as, "Peek should return 4")
	assert.Equal(t, &b, as, "pointers should be same.")

	as, ok = s.Pop().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 1, *as, "Peek should return 4")
	assert.Equal(t, &a, as, "pointers should be same.")

	bs := s.Pop()
	assert.Nil(t, bs, "stack is empty should return nil")

	// push something back again after empting the stack.
	s.Push(&d)
	assert.Equal(t, 1, s.Count(), "Count should be 1")

	as, ok = s.Pop().(*int)
	assert.Equal(t, true, ok, "Peek should be of type *int")
	assert.Equal(t, 4, *as, "Peek should return 4")
	assert.Equal(t, &d, as, "pointers should be same.")

}
func TestEmptyStack(t *testing.T) {
	s := Stack{}
	assert.Equal(t, 0, s.Count(), "Empty stack should return 0")
	assert.Nil(t, s.Peek(), "Empty stack Peek should return nil")
}
