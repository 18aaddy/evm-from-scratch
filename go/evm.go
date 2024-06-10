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
	// "encoding/binary"
	// "encoding/hex"
	// "encoding/hex"
	"math/big"

	// "golang.org/x/text/number"
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

	}	
		
		return stack, successOrNot
}


