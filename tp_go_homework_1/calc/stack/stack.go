package stack

import "fmt"

type Stack struct {
	buffer []int
	count  int
}

func (s *Stack) Push(v int) {
	//fmt.Printf("before %d: %d %d \n", v, len(s.buffer), s.count)

	if s.count == len(s.buffer) {
		s.buffer = append(s.buffer, v)
	} else {
		s.buffer[s.count] = v
	}

	s.count++
}

func (s *Stack) Pop() (v int, err error) {
	if s.Empty() {
		return 0, fmt.Errorf("stack is empty")
	}

	v = s.buffer[s.count-1]

	if s.count > 0 {
		s.count--
	}

	return v, nil
}

func (s *Stack) Top() (int, error) {
	if s.Empty() {
		return 0, fmt.Errorf("stack is empty")
	}

	return s.buffer[s.count-1], nil
}

func (s *Stack) Empty() bool {
	return s.count == 0
}

