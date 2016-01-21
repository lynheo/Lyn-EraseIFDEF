package main

type DepthStack struct {
	defineStack []string
}

func (s *DepthStack) Push(value string) {
	s.defineStack = append(s.defineStack, value)
}

func (s *DepthStack) Pop() {
	s.defineStack = s.defineStack[:len(s.defineStack)-1]
}

func (s *DepthStack) Top() string {
	return s.defineStack[len(s.defineStack)-1]
}
