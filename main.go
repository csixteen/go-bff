package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

type Stack []int

func (s Stack) Push(v int) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, int, error) {
	l := len(s)
	if l == 0 {
		return s, -1, errors.New("Empty stack")
	}
	return s[:l-1], s[l-1], nil
}

func readChar() byte {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return []byte(input)[0]
}

func executeInstruction(code []byte, ip int, stack Stack, mem []byte, cursor int) (Stack, int, int) {
	op := code[ip]
	offset := 1

	switch op {
	case '>':
		cursor += 1
	case '<':
		cursor -= 1
	case '+':
		mem[cursor] += 1
	case '-':
		mem[cursor] -= 1
	case ',':
		mem[cursor] = readChar()
	case '.':
		fmt.Printf("%c", mem[cursor])
	case '[':
		if mem[cursor] > 0 {
			stack = stack.Push(ip)
		} else {
			if len(stack) > 0 {
				stack, _, _ = stack.Pop()
			}
			offset = strings.Index(string(code[ip+1:]), "]")
			if offset == -1 {
				panic("malformed program")
			}
		}
	case ']':
		if mem[cursor] == 0 {
			stack, _, _ = stack.Pop()
		} else {
			stack, offset, _ = stack.Pop()
			offset -= ip
		}
	}

	return stack, ip + offset, cursor
}

func execute(code []byte) {
	memory := make([]byte, MEM_SIZE)
	cursor := 0 // memory cursor
	stack := make(Stack, 0)
	ip := 0 // Instruction Pointer
	for ip < len(code) {
		stack, ip, cursor = executeInstruction(code, ip, stack, memory, cursor)
	}
}

// Loads code from file and returns byte array that
// only contains valid Brainfuck operators.
func loadCode(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	c, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var contents []byte
	for _, char := range c {
		if _, ok := OPERATORS[char]; ok {
			contents = append(contents, char)
		}
	}
	return contents
}

func main() {
	if len(os.Args) != 2 {
		panic("Usage: <program> <file.bf>")
	}
	execute(loadCode(os.Args[1]))
}
