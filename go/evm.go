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
	// "fmt"
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

		//14. SIGN EXTEND 
		if (op == 0x0b) {
			if len(stack) < 2 {
				return nil, false
			}

			// Pop b from the stack

			// Pop x from the stack

			// Ensure b is within the bounds of our bit width (0-31 for a 256-bit number)

			// Pop b and x from the stack
			b := stack[0]
			x := stack[1]
			stack = stack[2:]

			// Calculate the sign extension mask
			bInt := int(b.Int64())
			if bInt >= 32 {
				return stack, false
			}
			bits := (bInt + 1) * 8
			signBit := new(big.Int).Lsh(big.NewInt(1), uint(bits-1))

			// Check if the sign bit is set
			if x.Cmp(signBit) >= 0 {
				// If the sign bit is set, extend with 1s
				extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
				extended.Sub(extended, big.NewInt(1))
				extended.Lsh(extended, uint(bits))
				x.Or(x, extended)
			} else {
				// Ensure higher bits are zero
				mask := new(big.Int).Lsh(big.NewInt(1), uint(bits))
				mask.Sub(mask, big.NewInt(1))
				x.And(x, mask)
			}

			// Push the result back onto the stack
			stack = append([]*big.Int{x}, stack...)

		}

		//15. SDIV
		if (op == 0x05) {
			if len(stack) < 2 {
				return nil, false
			}

			if stack[1].Cmp(big.NewInt(0)) == 0 {
				stack = stack[2:]
				stack = append([]*big.Int{big.NewInt(0)}, stack...)
			} else {
				value1 := stack[0].Int64()
				int8Value1 := int8(value1)
				value2 := stack[1].Int64()
				int8Value2 := int8(value2)

				value := int8Value1 / int8Value2

				bits := 8

				// Check if the sign bit is set
				if value < 0 {
					value8 := new(big.Int).Add(big.NewInt(int64(256)), big.NewInt(int64(value)))
					// If the sign bit s set, extend with 1s
					extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
					extended.Sub(extended, big.NewInt(1))
					extended.Lsh(extended, uint(bits))
					value8.Or(value8, extended)
					stack = stack[2:]
					stack = append([]*big.Int{value8}, stack...)
				} else {
					stack = stack[2:]
					stack = append([]*big.Int{big.NewInt(int64(value))}, stack...)
				}

			}
				
		}

		//16. SMOD
		if (op == 0x07) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]
			
			if number2.Cmp(big.NewInt(0)) == 0 {
				stack = append([]*big.Int{big.NewInt(0)}, stack[2:]... )
			} else {
			
				number1int8 := int8(number1.Int64())
				number2int8 := int8(number2.Int64())
				
				bits := 8
				value := number1int8 % number2int8

				if value < 0 {
					value8 := new(big.Int).Add(big.NewInt(int64(256)), big.NewInt(int64(value)))
					// If the sign bit s set, extend with 1s
					extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
					extended.Sub(extended, big.NewInt(1))
					extended.Lsh(extended, uint(bits))
					value8.Or(value8, extended)
					stack = stack[2:]
					stack = append([]*big.Int{value8}, stack...)				

				} else {
					stack = stack[2:]
					stack = append([]*big.Int{big.NewInt(int64(value))}, stack...)
				}
			}

		}

		LessThan := func (number1 *big.Int, number2 *big.Int) {
			intNumber1 := number1.Int64()
			intNumber2 := number2.Int64()

			if intNumber1 < intNumber2 {
				stack = append([]*big.Int{big.NewInt(1)}, stack[2:]...)
			} else {
				stack = append([]*big.Int{big.NewInt(0)}, stack[2:]...)
			}
		}

		GreaterThan := func (number1 *big.Int, number2 *big.Int) {
			intNumber1 := number1.Int64()
			intNumber2 := number2.Int64()

			if intNumber1 > intNumber2 {
				stack = append([]*big.Int{big.NewInt(1)}, stack[2:]...)
			} else {
				stack = append([]*big.Int{big.NewInt(0)}, stack[2:]...)
			}
		}

		EqualTo := func (number1 *big.Int, number2 *big.Int) {
			intNumber1 := number1.Int64()
			intNumber2 := number2.Int64()

			if intNumber1 == intNumber2 {
				stack = append([]*big.Int{big.NewInt(1)}, stack[2:]...)
			} else {
				stack = append([]*big.Int{big.NewInt(0)}, stack[2:]...)
			}
		}

		//17. LT
		if (op == 0x10 || op == 0x12) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]

			LessThan(number1, number2)			
		}
		
		//18. GT
		if (op == 0x11 || op == 0x13) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]

			GreaterThan(number1, number2)			
		}		

		//19. EQ
		if (op == 0x14) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]

			EqualTo(number1, number2)
		}

		//20. ISZERO
		if (op == 0x15) {
			if len(stack) < 1 {
				return nil, false
			}

			number1 := stack[0]

			if number1.Cmp(big.NewInt(0)) == 0 {
				stack = append([]*big.Int{big.NewInt(1)}, stack[1:]...)
			} else {
				stack = append([]*big.Int{big.NewInt(0)}, stack[1:]...)
			}
		}

		//21. NOT
		if (op == 0x19) {
			if len(stack) < 1 {
				return nil, false
			}

			number1 := stack[0]

			mask := new(big.Int).Lsh(big.NewInt(1), 256) 
			mask.Sub(mask, big.NewInt(1))               
		
			result := new(big.Int).Xor(number1, mask)
		
			stack = append([]*big.Int{result}, stack[1:]...)
		}

		//22. AND
		if (op == 0x16) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]

			result := new(big.Int).And(number1, number2)
		
			stack = append([]*big.Int{result}, stack[2:]...)
		}

		//23. OR
		if (op == 0x17) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]

			result := new(big.Int).Or(number1, number2)
		
			stack = append([]*big.Int{result}, stack[2:]...)
		}

		//24. XOR
		if (op == 0x18) {
			if len(stack) < 2 {
				return nil, false
			}

			number1 := stack[0]
			number2 := stack[1]

			result := new(big.Int).Xor(number1, number2)
		
			stack = append([]*big.Int{result}, stack[2:]...)
		}

		//25. SHL
		if (op == 0x1b) {
			if len(stack) < 2 {
				return nil, false
			}
			
			number1 := stack[0]
			number2 := stack[1]

			result := new(big.Int).Lsh(number1, uint(number2.Int64()))
			
			
			stack = append([]*big.Int{result}, stack[2:]...)
		}

	}	
		
		return stack, successOrNot
}


