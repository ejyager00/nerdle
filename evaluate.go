package main

import (
	"math"
	"strconv"
	"strings"
)

func precedence(x, y rune) int {
	// * and / have greater precedence than + and -
	if x == '+' || x == '-' {
		if y == '+' || y == '-' {
			return 0
		} else if y == '*' || y == '/' {
			return -1
		}
	} else if x == '*' || x == '/' {
		if y == '+' || y == '-' {
			return 1
		} else if y == '*' || y == '/' {
			return 0
		}
	}
	return 0
}

func shuntingYard(expression string) string {
	// convert infix to postfix
	ops_stack := make([]rune, len(expression))
	ops_last := -1
	var output string = ""
	prior := 0
	for i, c := range expression {
		if strings.ContainsRune("+-*/", c) {
			if c == '-' && (i == 0 || strings.ContainsRune("+-*/", rune(expression[i-1]))) {
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
			if i != len(expression)-1 && strings.ContainsRune("0123456789", rune(expression[i+1])) {
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

func postfixCalc(expression string) float64 {
	// evaluate postfix expression
	stack := make([]float64, len(expression))
	stack_loc := -1
	last_open := 0
	for i, c := range expression {
		if c == '-' && last_open == i-1 {
		} else if strings.ContainsRune("+-*/", c) {
			stack_loc -= 1
			switch c {
			case '+':
				stack[stack_loc] = stack[stack_loc] + stack[stack_loc+1]
			case '-':
				stack[stack_loc] = stack[stack_loc] - stack[stack_loc+1]
			case '*':
				stack[stack_loc] = stack[stack_loc] * stack[stack_loc+1]
			case '/':
				stack[stack_loc] = stack[stack_loc] / stack[stack_loc+1]
			}
		} else if c == ')' {
			stack_loc++
			new, _ := strconv.Atoi(expression[last_open+1 : i])
			stack[stack_loc] = float64(new)
		} else if c == '(' {
			last_open = i
		}
	}
	return stack[0]
}

func IsEqual(equation string) bool {
	// returns the truth value of the equation
	sides := strings.Split(equation, "=")
	left := sides[0]
	ans, _ := strconv.Atoi(sides[1])
	left = shuntingYard(left)
	val := postfixCalc(left)
	if math.Mod(val, 1) != 0 {
		return false
	}
	return int(val) == ans
}
