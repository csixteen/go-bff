package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const MEM_SIZE = 30000

var OPERATORS = map[byte]bool{
	'<': true,
	'>': true,
	'+': true,
	'-': true,
	',': true,
	'.': true,
	'[': true,
	']': true,
}

func readChar() byte {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return []byte(input)[0]
}

// ------------------------------------------------------------------
// `Stack` structure and methods

type stack struct {
	s []int
}

func NewStack() *stack {
	return &stack{
		s: make([]int, 0),
	}
}

func (s *stack) Push(v int) {
	s.s = append(s.s, v)
}

func (s *stack) Pop() (int, error) {
	l := len(s.s)
	if l == 0 {
		return -1, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) Top() (int, error) {
	l := len(s.s)
	if l == 0 {
		return -1, errors.New("Empty Stack")
	}

	return s.s[l-1], nil
}

// -------------------------------------------------------------------
// `Program` structure and methods

type Program struct {
	Code   []byte // Brainfuck code
	Ip     int    // Instruction pointer
	Cursor int    // Memory cursor
	Stack  *stack // Program stack
	Cells  []byte // Cells
}

func NewProgram(code []byte, n int) *Program {
	return &Program{
		Code:   code,
		Ip:     0,
		Cursor: 0,
		Stack:  NewStack(),
		Cells:  make([]byte, n),
	}
}

func (p *Program) executeInstruction() {
	offset := 1

	switch p.Code[p.Ip] {
	case '>':
		p.Cursor += 1
	case '<':
		p.Cursor -= 1
	case '+':
		p.Cells[p.Cursor] += 1
	case '-':
		p.Cells[p.Cursor] -= 1
	case ',':
		p.Cells[p.Cursor] = readChar()
	case '.':
		fmt.Printf("%c", p.Cells[p.Cursor])
	case '[':
		if p.Cells[p.Cursor] > 0 {
			p.Stack.Push(p.Ip)
		} else {
			i := p.Ip + 1
			skipped := 0

			for i < len(p.Code) {
				switch p.Code[i] {
				case '[':
					skipped += 1
				case ']':
					if skipped == 0 {
						i += 1
						break
					} else {
						skipped -= 1
					}
				}

				i += 1
			}

			offset = i - p.Ip
		}
	case ']':
		if p.Cells[p.Cursor] == 0 {
			p.Stack.Pop()
		} else {
			top, _ := p.Stack.Top()
			offset = top - p.Ip + 1
		}
	}

	p.Ip += offset
}

func (p *Program) execute() {
	for p.Ip < len(p.Code) {
		p.executeInstruction()
	}
}

// --------------------------------------------------------------

// Loads code from file and returns byte array that
// only contains valid Brainfuck operators.
func loadCode(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	code, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var contents []byte
	for _, char := range code {
		if OPERATORS[char] {
			contents = append(contents, char)
		}
	}
	return contents
}

// ---------------------------------------------------------------

// Entry point
func main() {
	if len(os.Args) != 2 {
		panic("Usage: <program> <file.bf>")
	}

	p := NewProgram(loadCode(os.Args[1]), MEM_SIZE)
	p.execute()
}
