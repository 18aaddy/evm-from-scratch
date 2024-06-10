// Package evm is an **incomplete** implementation of the Ethereum Virtual
// Machine for the "EVM From Scratch" course:
// https://github.com/w1nt3r-eth/evm-from-scratch
//
// To work on EVM From Scratch In Go:
//
// - Install Golang: https://golang.org/doc/install
// - Go to the `go` directory: `cd go`
// - Edit `evm.go` (this file!), see TODO below
// - Run `go test ./...` to run the tests
package evm

import (
	"math/big"

)

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
	var stack []*big.Int
	pc := 0
	successOrNot := true

	for pc < len(code) {
		op := code[pc]
		pc++

		// TODO: Implement the EVM here!

		// 1. STOP 
		if (op == 0x00) {
			break
		}

		//2. PUSH0
		if (op == 0x5f) {
			stack = append([]*big.Int{big.NewInt(0)}, stack...)
		}
		
		// General PUSH function
		Push := func(increment int) {
			var toAppendInBytes []byte

			if (pc + increment > len(code)) {    // In case given byte code is too short
				successOrNot = false
				return
			}
			toAppendInBytes = code[pc:pc+increment]

			stack = append([]*big.Int{new(big.Int).SetBytes(toAppendInBytes)}, stack...)
			pc += increment 
		}
		
		//3. PUSH1 - PUSH32
		if (0x60 <= op && op <= 0x7f) {
			increment := int(op - 0x60) + 1
			Push(increment)
		}

		//5. POP
		if (op == 0x50) {
			if (len(stack) < 1) {
				successOrNot = false
				break
			}
			stack = stack[1:]
		}

		//6. ADD
		if (op == 0x01) {
			if (len(stack) < 2) {        // In case stack does not have 2 numbers to add
				successOrNot = false
				break
			}
			
			number1 := stack[0]
			number2 := stack[1]
			result := new(big.Int)
			max_uint := result.Exp(big.NewInt(2), big.NewInt(256), nil)

			sum := new(big.Int).Add(number1, number2)

			sum.Mod(sum, max_uint)

			stack = append([]*big.Int{sum}, stack[2:]...)
		}

		//7. MUL 
		if (op == 0x02) {
			if (len(stack) < 2) {        // In case stack does not have 2 numbers to multiply
				successOrNot = false
				break
			}
			number1 := stack[0]
			number2 := stack[1]
			result := new(big.Int)
			stack = append([]*big.Int{new(big.Int).Mod((new(big.Int).Mul(number1, number2)), result.Exp(big.NewInt(2), big.NewInt(256), nil))}, stack[2:]...)
		}

		//8. SUB
		if (op == 0x03) {
			if (len(stack) < 2) {        // In case stack does not have 2 numbers to subtract
				successOrNot = false
				break
			}
			number1 := stack[0]
			number2 := stack[1]
			answer := new(big.Int).Sub(number1, number2)

			result := new(big.Int)
			max_uint := result.Exp(big.NewInt(2), big.NewInt(256), nil)

			answer.Mod(answer, max_uint)
			
			stack = append([]*big.Int{answer}, stack[2:]...)
		}

		//9. DIV
		if (op == 04) {
			if (len(stack) < 2) {        // In case stack does not have 2 numbers to divide
				successOrNot = false
				break
			}
			var answer *big.Int

			number1 := stack[0]
			number2 := stack[1]

			if (number2.Cmp(big.NewInt(0)) == 0) {
				answer = big.NewInt(0)       // If someone tries to divide by zero
			} else {
				answer = new(big.Int).Div(number1, number2)
			}

			result := new(big.Int)
			max_uint := result.Exp(big.NewInt(2), big.NewInt(256), nil)

			answer.Mod(answer, max_uint)
			stack = append([]*big.Int{answer}, stack[2:]...)
		}

		//10. MOD
		if (op == 0x06) {
			if (len(stack) < 2) {        // In case stack does not have 2 numbers to modulate
				successOrNot = false
				break
			}

			var answer *big.Int

			number1 := stack[0]
			number2 := stack[1]

			if (number2.Cmp(big.NewInt(0)) == 0) {
				answer = big.NewInt(0)
			} else {
				answer = new(big.Int).Mod(number1, number2)
			}

			stack = append([]*big.Int{answer}, stack[2:]...)
		}

		//11. ADDMOD
		if (op == 0x08) {
			if (len(stack) < 3) {        // In case stack does not have 3 numbers to add mod
				successOrNot = false
				break
			}

			var answer *big.Int

			number1 := stack[0]
			number2 := stack[1]
			number3 := stack[2]

			result := new(big.Int)
			max_uint := result.Exp(big.NewInt(2), big.NewInt(256), nil)

			sum := new(big.Int).Add(number1, number2)
			sum.Mod(sum, max_uint)

			if (number3.Cmp(big.NewInt(0)) == 0) {
				successOrNot = false
				break
			} else {
				answer = new(big.Int).Mod(sum, number3)
			}

			stack = append([]*big.Int{answer}, stack[3:]...)
		}

		//12. MULMOD
		if (op == 0x09) {
			if (len(stack) < 3) {        // In case stack does not have 3 numbers to mul mod
				successOrNot = false
				break
			}

			var answer *big.Int

			number1 := stack[0]
			number2 := stack[1]
			number3 := stack[2]

			result := new(big.Int)
			max_uint := result.Exp(big.NewInt(2), big.NewInt(256), nil)

			product := new(big.Int).Mul(number1, number2)

			if (number3.Cmp(big.NewInt(0)) == 0) {
				successOrNot = false
				break
			} else {
				answer = new(big.Int).Mod(product, number3)
			}

			answer.Mod(answer, max_uint)

			stack = append([]*big.Int{answer}, stack[3:]...)
		}

		//13. EXP
		if (op == 0x0a) {
			if (len(stack) < 2) {        // In case stack does not have 3 numbers to exponentiate
				successOrNot = false
				break
			}		
			
			number1 := stack[0]
			number2 := stack[1]
			
			result := new(big.Int)
			max_uint := result.Exp(big.NewInt(2), big.NewInt(256), nil)

			answer := new(big.Int).Exp(number1, number2, max_uint)

			stack = append([]*big.Int{answer}, stack[2:]...)
		}
	}	
		
		return stack, successOrNot
}


