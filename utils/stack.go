package utils

type Stack struct {
	stack []interface{}
	count int
}

func (s *Stack) Push(obj interface{}) {
	if len(s.stack) <= s.count {
		s.stack = append(s.stack, obj)
		s.count++
	} else {
		s.stack[s.count] = obj
		s.count++
	}
}

func (s *Stack) Pop() interface{} {
	if s.count == 0 {
		return nil
	}
	s.count--
	return s.stack[s.count]
}
func (s *Stack) Count() int {
	return s.count
}
func (s *Stack) Peek() interface{} {
	if s.count == 0 {
		return nil
	}
	return s.stack[s.count-1]
}
