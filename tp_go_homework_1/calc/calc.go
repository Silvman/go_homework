package main

import (
	"./stack"
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

func plus(s *stack.Stack) error {
	a, _ := s.Pop()
	if b, err := s.Pop(); err == nil {
		s.Push(a + b)
	} else {
		return err
	}

	return nil
}

func minus(s *stack.Stack) error {
	a, _ := s.Pop()
	if b, err := s.Pop(); err == nil {
		s.Push(b - a)
	} else {
		return err
	}

	return nil
}

func div(s *stack.Stack) error {
	a, _ := s.Pop()
	if b, err := s.Pop(); err == nil {
		s.Push(b / a)
	} else {
		return err
	}

	return nil
}

func mult(s *stack.Stack) error {
	a, _ := s.Pop()
	if b, err := s.Pop(); err == nil {
		s.Push(a * b)
	} else {
		return err
	}

	return nil
}

func equal(s *stack.Stack, writer io.Writer) error {
	if a, err := s.Pop(); err == nil {
		fmt.Fprintf(writer, "Result = %d", a)
	} else {
		return err
	}

	return nil
}

func calc(input io.Reader, output io.Writer) error {
	sc := bufio.NewReader(input)
	st := new(stack.Stack)

	for {
		var s string
		var d int

		if _, err := fmt.Fscanf(sc, "%1s", &s); err == nil { // %s - чтобы обрасывать пробелы и другие лишние символы
			r := []rune(s)[0]
			if unicode.IsDigit(r) {
				sc.UnreadRune() // если сосканировали цифру как символ, возвращаем её обратно в буфер для последующего перескана
				if _, err := fmt.Fscanf(sc, "%d", &d); err == nil {
					st.Push(d)
				}
			} else {
				var err error

				switch r {
				case '=':
					err = equal(st, output)
				case '+':
					err = plus(st)
				case '-':
					err = minus(st)
				case '/':
					err = div(st)
				case '*':
					err = mult(st)
				}

				if err != nil {
					return err
				}
			}
		} else if err == io.EOF {
			break
		}
	}

	return nil
}

func main() {
	if err := calc(os.Stdin, os.Stdout); err != nil {
		fmt.Println(err)
	}

	return
}
