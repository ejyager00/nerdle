package main

import (
	"errors"
	"strconv"
	"strings"
)

const NUMBERS = "0123456789"
const OPERATORS = "+-*/"

var OP_PRECDENCE = map[rune]int{
	'*': 1,
	'/': 1,
	'+': 0,
	'-': 0,
}

func validChars(equation string, include_equal bool) bool {
	var letters string
	if include_equal {
		letters = NUMBERS + OPERATORS + "="
	} else {
		letters = NUMBERS + OPERATORS
	}
	for _, c := range equation {
		if !strings.ContainsRune(letters, c) {
			return false
		}
	}
	return true
}

func precedence(x, y rune) int {
	if !strings.ContainsRune(OPERATORS, x) || !strings.ContainsRune(OPERATORS, y) {
		return 0
	}
	return OP_PRECDENCE[x] - OP_PRECDENCE[y]
}

func shuntingYard(expression string) string {
	// convert infix to postfix
	ops_stack := make([]rune, len(expression))
	ops_last := -1
	var output string = ""
	prior := 0
	for i, c := range expression {
		if strings.ContainsRune(OPERATORS, c) {
			if c == '-' && (i == 0 || strings.ContainsRune(OPERATORS, rune(expression[i-1]))) {
				prior++
			} else {
				if ops_last != -1 && precedence(c, ops_stack[ops_last]) <= 0 {
					for ops_last >= 0 && precedence(c, ops_stack[ops_last]) <= 0 {
						output += string(ops_stack[ops_last])
						ops_last--
					}
				}
				ops_last++
				ops_stack[ops_last] = c
			}
		} else {
			if i != len(expression)-1 && strings.ContainsRune(NUMBERS, rune(expression[i+1])) {
				prior++
			} else {
				output += "(" + expression[i-prior:i+1] + ")"
				prior = 0
			}
		}
	}
	for i := ops_last; i >= 0; i-- {
		output += string(ops_stack[i])
	}
	return output
}

func postfixCalc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, errors.New("empty expression")
	}
	// evaluate postfix expression
	stack := make([]float64, len(expression))
	stack_loc := -1
	last_open := 0
	last_closed := -1
	for i, c := range expression {
		if c == '-' && last_open > last_closed {
		} else if strings.ContainsRune(OPERATORS, c) {
			stack_loc -= 1
			if stack_loc < 0 {
				return 0, errors.New("missing operands")
			}
			switch c {
			case '+':
				stack[stack_loc] = stack[stack_loc] + stack[stack_loc+1]
			case '-':
				stack[stack_loc] = stack[stack_loc] - stack[stack_loc+1]
			case '*':
				stack[stack_loc] = stack[stack_loc] * stack[stack_loc+1]
			case '/':
				if stack[stack_loc+1] == 0 {
					return 0, errors.New("divide by zero")
				}
				stack[stack_loc] = stack[stack_loc] / stack[stack_loc+1]
			}
		} else if c == ')' {
			stack_loc++
			new, err := strconv.Atoi(expression[last_open+1 : i])
			if err != nil {
				return 0, errors.New("double negative")
			}
			stack[stack_loc] = float64(new)
			last_closed = i
		} else if c == '(' {
			last_open = i
		}
	}
	return stack[0], nil
}

func IsEqual(equation string) (bool, error) {
	// returns the truth value of the equation
	if !validChars(equation, true) {
		return false, errors.New("invalid character")
	}
	sides := strings.Split(equation, "=")
	if len(sides) < 2 {
		return false, errors.New("no equality")
	} else if len(sides) > 2 {
		return false, errors.New("multiple equalities")
	}
	left := shuntingYard(sides[0])
	right := shuntingYard(sides[1])
	left_val, err := postfixCalc(left)
	if err != nil {
		return false, err
	}
	right_val, err := postfixCalc(right)
	if err != nil {
		return false, err
	}
	return left_val == right_val, nil
}

func LeadingZeros(puzzle string) bool {
	for i, c := range puzzle {
		if c == '0' && i != len(puzzle)-1 && strings.ContainsRune(NUMBERS, rune(puzzle[i+1])) {
			return true
		}
	}
	return false
}

func NegativeZeros(puzzle string) bool {
	searching_for_num := false
	for i, c := range puzzle {
		if c == '0' && i > 0 && puzzle[i-1] == '-' && (i < 2 || strings.ContainsRune(OPERATORS+"=", rune(puzzle[i-2]))) {
			searching_for_num = true
		} else if searching_for_num {
			if strings.ContainsRune(NUMBERS, c) && c != '0' {
				searching_for_num = false
			} else if strings.ContainsRune(OPERATORS+"=", c) {
				return true
			}
		}
	}
	return false
}

func ContainsOperator(puzzle string) []bool {
	expressions := strings.Split(puzzle, "=")
	results := make([]bool, len(expressions))
	for i, e := range expressions {
		results[i] = false
		for _, c := range OPERATORS {
			if strings.ContainsRune(e, c) {
				results[i] = true
				break
			}
		}
	}
	return results
}

func MultiplicationByZero(puzzle string) (bool, error) {
	_, err := IsEqual(puzzle)
	if err != nil {
		return false, err
	}
	expressions := strings.Split(puzzle, "=")
	for _, e := range expressions {
		e = shuntingYard(e)
		// evaluate postfix expression
		stack := make([]float64, len(e))
		stack_loc := -1
		last_open := 0
		last_closed := -1
		for i, c := range e {
			if c == '-' && last_open > last_closed {
			} else if strings.ContainsRune(OPERATORS, c) {
				stack_loc -= 1
				switch c {
				case '+':
					stack[stack_loc] = stack[stack_loc] + stack[stack_loc+1]
				case '-':
					stack[stack_loc] = stack[stack_loc] - stack[stack_loc+1]
				case '*':
					if stack[stack_loc] == 0 || stack[stack_loc+1] == 0 {
						return true, nil
					}
					stack[stack_loc] = stack[stack_loc] * stack[stack_loc+1]
				case '/':
					stack[stack_loc] = stack[stack_loc] / stack[stack_loc+1]
				}
			} else if c == ')' {
				stack_loc++
				new, _ := strconv.Atoi(e[last_open+1 : i])
				stack[stack_loc] = float64(new)
				last_closed = i
			} else if c == '(' {
				last_open = i
			}
		}
	}
	return false, nil
}

func DivisionOfZero(puzzle string) (bool, error) {
	_, err := IsEqual(puzzle)
	if err != nil {
		return false, err
	}
	expressions := strings.Split(puzzle, "=")
	for _, e := range expressions {
		e = shuntingYard(e)
		// evaluate postfix expression
		stack := make([]float64, len(e))
		stack_loc := -1
		last_open := 0
		last_closed := -1
		for i, c := range e {
			if c == '-' && last_open > last_closed {
			} else if strings.ContainsRune(OPERATORS, c) {
				stack_loc -= 1
				switch c {
				case '+':
					stack[stack_loc] = stack[stack_loc] + stack[stack_loc+1]
				case '-':
					stack[stack_loc] = stack[stack_loc] - stack[stack_loc+1]
				case '*':
					stack[stack_loc] = stack[stack_loc] * stack[stack_loc+1]
				case '/':
					if stack[stack_loc] == 0 {
						return true, nil
					}
					stack[stack_loc] = stack[stack_loc] / stack[stack_loc+1]
				}
			} else if c == ')' {
				stack_loc++
				new, _ := strconv.Atoi(e[last_open+1 : i])
				stack[stack_loc] = float64(new)
				last_closed = i
			} else if c == '(' {
				last_open = i
			}
		}
	}
	return false, nil
}
