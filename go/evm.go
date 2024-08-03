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
	"fmt"
	"math/big"
	"golang.org/x/crypto/sha3"
	"encoding/hex"
)



//*** Block ***//

type Block struct {
	Basefee  string `json:"basefee"`
	Coinbase string `json:"coinbase"`
	Timestamp string `json:"timestamp"`
	Number    string `json:"number"`
	Difficulty string `json:"difficulty"`
	Gaslimit   string `json:"gaslimit"`
	Chainid    string `json:"chainid"`
	Blockhash  int 
}

//*** State and Account ***//

type Account struct {
	Nonce    uint64
	Balance  string `json:"balance"`
	Storage  map[string]string
	CodeHash Code `json:"code"`
}

type Code struct {
	Asm string `json:"asm"`
	Bin string `json:"bin"`
}

type State struct {
	Accounts map[string]Account
}

func NewState() *State {
	return &State{
		Accounts: make(map[string]Account),
	}
}

func (s *State) GetAccount(address string) Account {
	return s.Accounts[address]
}

// func (s *State) CreateAccount(address string, balance *big.Int, codeHash string) {
// 	s.Accounts[address] = &Account{
// 		Nonce:    0,
// 		Balance:  balance,
// 		Storage:  make(map[string]string),
// 		CodeHash: codeHash,
// 	}
// }

// func (s *State) UpdateAccount(address string, balance *big.Int, nonce uint64) {
// 	account := s.Accounts[address]
// 	if account != nil {
// 		account.Balance = balance
// 		account.Nonce = nonce
// 	}
// }

//*** Memory ***//

type Memory struct {
	data  []byte
}

type Transaction struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Origin   string `json:"origin"`
	Gasprice string `json:"gasprice"`
	Value    string `json:"value"`
	Data     string `json:"data"`
}

func NewMemory(size int) *Memory {
	return &Memory {
		data: make([]byte, size),
	}
}

func (m *Memory) Mstore(value []byte, offset int) {
	m.MSIZE(offset)
	copy(m.data[offset:], value)
}

func (m *Memory) Mstore8(value byte, offset int) {
	m.MSIZE(offset - 32)
	m.data[offset] = value
}

func (m *Memory) MSIZE(offset int) int {
	for i := 0 ; offset + 32 > len(m.data) ; i++ {
		extend := make([]byte, 32)
		m.data = append(m.data, extend...)
	}
	return len(m.data)
}

func (m *Memory) MLoad(offset int, size int) []byte {
	m.MSIZE(offset)
	return m.data[offset:offset + size]
}


//****  EVM  FUNCTION  ****//

// Run runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte, transaction Transaction, block Block, state State) ([]*big.Int ,bool) {
	var stack []*big.Int
	pc := 0
	successOrNot := true

	m := NewMemory(0)

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
		
			shiftAmount := stack[0]
			value := stack[1]
		
			shiftInt := shiftAmount.Uint64()
			result := new(big.Int).Lsh(value, uint(shiftInt))
		
			mask := new(big.Int).Lsh(big.NewInt(1), 256) 
			mask.Sub(mask, big.NewInt(1))               
			result.And(result, mask)
		
			stack = append([]*big.Int{result}, stack[2:]...)
		}

		//26. SHR
		if (op == 0x1c) {
			if len(stack) < 2 {
				return nil, false
			}
		
			shiftAmount := stack[0]
			value := stack[1]
		
			shiftInt := shiftAmount.Uint64()
			result := new(big.Int).Rsh(value, uint(shiftInt))
		
			mask := new(big.Int).Lsh(big.NewInt(1), 256) 
			mask.Sub(mask, big.NewInt(1))               
			result.And(result, mask)
		
			stack = append([]*big.Int{result}, stack[2:]...)
		}

		//27. SAR
		if op == 0x1d {
			if len(stack) < 2 {
				return nil, false
			}
		
			shiftAmount := stack[0]
			value := stack[1]
			stack = stack[2:]
		
			shiftInt := int(shiftAmount.Uint64())
		
			if shiftInt >= 256 {

				if negative_converter(value).Sign() < 0 {
					result := new(big.Int).Lsh(big.NewInt(1), 256)
					result.Sub(result, big.NewInt(1)) 
					stack = append([]*big.Int{result}, stack...)
				} else {
					stack = append([]*big.Int{big.NewInt(0)}, stack...)
				}

			} else {
				result := new(big.Int).Rsh(value, uint(shiftInt))
		
				if negative_converter(value).Sign() < 0 {

					signBit := new(big.Int).Lsh(big.NewInt(1), 256-uint(shiftInt))
					signBit.Sub(signBit, big.NewInt(1))

					extended := new(big.Int).Lsh(big.NewInt(1), uint(shiftInt))
					extended.Sub(extended, big.NewInt(1))
					extended.Lsh(extended, uint(256-uint(shiftInt)))
					result.Or(result, extended)
				}
		
				stack = append([]*big.Int{result}, stack...)

				fmt.Println(stack)
			}
		}
		
		//28. BYTE
		if op == 0x1a {
			if len(stack) < 2 {
				return nil, false
			}

			shift := stack[0]
			value := stack[1]
			
			shiftInt := shift.Int64()

			if shiftInt < 32 {

				extended := new(big.Int).Lsh(big.NewInt(1), 32*8)
				mask := new(big.Int).Lsh(big.NewInt(1), (32 - uint(shiftInt)) * 8)

				extended.Sub(extended, mask)
				value.And(extended.Not(extended), value)

				value.Rsh(value, 8 * uint(31 - shiftInt))

				stack = append([]*big.Int{value}, stack[2:]...)
			
			} else {
				stack = append([]*big.Int{big.NewInt(0)}, stack[2:]...)
			}
		}

		//29. DUP1-16 
		if 0x80 <= op && op <= 0x8f {
			element := op - 0x80 + 1

			if len(stack) < int(element) {
				return nil, false
			}

			value := stack[int(element) - 1]
			
			stack = append([]*big.Int{value}, stack...)
		}

		//30. SWAP1-16
		if 0x90 <= op && op <= 0x9f {
			swap := int(op - 0x90 + 1) 

			topItem := stack[0]
			swapItem := stack[swap]

			if len(stack) < swap {
				return nil, false
			}
			
			stack[swap] = topItem
			stack[0] = swapItem
		}

		//31. INVALID
		if op == 0xfe {
			return nil, false
		}

		//32. PC
		if op == 0x58 {
			stack = append([]*big.Int{big.NewInt(int64(pc - 1))}, stack...)
		}

		//33. GAS
		if op == 0x5a {			// TODO: Can add actual Gas functionality
			max_uint := new(big.Int).Lsh(big.NewInt(1), 256)
			max_uint.Sub(max_uint, big.NewInt(1))
			stack = append([]*big.Int{max_uint}, stack...)
		} 

		//34. JUMP, JUMP 1 
		if op == 0x56 || op == 0x57 {
			
			value := stack[0]
			jumpOrNot := true
			
			if op == 0x57 {
				if len(stack) < 2{
					return nil, false
				}
				jumpOrNot = stack[1].Cmp(big.NewInt(0)) == 0  // If true, means don't jump
				stack = stack[1:]
			}
			
			stack = stack[1:]

			if !jumpOrNot || op == 0x56 { 
				if int(value.Int64()) > len(code) - 1 {
					return nil, false
				}

				pc = int(value.Int64())
				op = code[pc]

				if op != 0x5b {
					return nil, false
				} else {
					for i := 0; i<pc; i++ {
						if code[i] == 0x00 {
							return nil, false
						}
						
						if 0x60 <= code[i] && code[i] <= 0x7f {
							increment := int(code[i] - 0x60) + 1
							if i < pc && pc <= i + increment {
								return nil, false
							}
						}
					}
				}
			}

		}

		//35. MSTORE
		if op == 0x52 || op == 0x53{
			if len(stack) < 2 {
				return nil, false
			}

			offset := stack[0]
			offsetInt := int(offset.Int64())
			value := stack[1]

			stack = stack[2:]
			valueBytes := value.Bytes()

			if op == 0x53 {
				m.Mstore8(valueBytes[0], offsetInt)
			} else {

				if len(valueBytes) < 32 {
					padding := make([]byte, 32 - len(valueBytes))
					valueBytes = append(padding, valueBytes...)
				}

				m.Mstore(valueBytes, offsetInt)
			}
		}

		//36. MLOAD
		if op == 0x51 {
			offset := stack[0]

			value := m.MLoad(int(offset.Int64()), 32)

			stack = append([]*big.Int{new(big.Int).SetBytes(value)}, stack[1:]...)
		}

		//37. MSIZE 
		if op == 0x59 {
			stack = append([]*big.Int{big.NewInt(int64(m.MSIZE(-32)))}, stack...)
		}

		//38. SHA3
		if op == 0x20 {
			offset := stack[0].Int64()
			size := stack[1].Int64()

			data := m.MLoad(int(offset), int(size))
			
			hasher := sha3.NewLegacyKeccak256()
			hasher.Write(data)
			
			hash := hasher.Sum(nil)

			stack = append([]*big.Int{new(big.Int).SetBytes(hash)}, stack[2:]...)
		}

		//39. ADDRESS
		if op == 0x30 {
			if len(transaction.To) == 0 {
				return nil, false
			}

			address, _ := new(big.Int).SetString(transaction.To[2:], 16)

			stack = append([]*big.Int{address}, stack...)
		}

		//40. CALLER
		if op == 0x33 {
			if len(transaction.From) == 0 {
				return nil, false
			}

			address, _ := new(big.Int).SetString(transaction.From[2:], 16)

			stack = append([]*big.Int{address}, stack...)
		}

		//41. ORIGIN
		if op == 0x32 {
			if len(transaction.Origin) == 0 {
				return nil, false
			}
			
			address, _ := new(big.Int).SetString(transaction.Origin[2:], 16)

			stack = append([]*big.Int{address}, stack...)
		}

		//42. GASPRICE
		if op == 0x3a {
			if len(transaction.Gasprice) == 0 {
				return nil, false
			}
			
			price, _ := new(big.Int).SetString(transaction.Gasprice[2:], 16)

			stack = append([]*big.Int{price}, stack...)
		}

		//43. BASEFEE
		if op == 0x48 {
			if len(block.Basefee) == 0 {
				return nil, false
			}

			basefee, _ := new(big.Int).SetString(block.Basefee[2:], 16)

			stack = append([]*big.Int{basefee}, stack...)
		}

		//44. COINBASE
		if op == 0x41 {
			if len(block.Coinbase) == 0 {
				return nil, false
			}

			coinbase, _ := new(big.Int).SetString(block.Coinbase[2:], 16)

			stack = append([]*big.Int{coinbase}, stack...)
		}

		//45. TIMESTAMP
		if op == 0x42 {
			if len(block.Timestamp) == 0 {
				return nil, false
			}

			timestamp, _ := new(big.Int).SetString(block.Timestamp[2:], 16)

			stack = append([]*big.Int{timestamp}, stack...)
		}

		//46. NUMBER
		if op == 0x43 {
			if len(block.Number) == 0 {
				return nil, false
			}

			number, _ := new(big.Int).SetString(block.Number[2:], 16)

			stack = append([]*big.Int{number}, stack...)
		}

		//47. DIFFICULTY
		if op == 0x44 {
			if len(block.Difficulty) == 0 {
				return nil, false
			}

			difficulty, _ := new(big.Int).SetString(block.Difficulty[2:], 16)

			stack = append([]*big.Int{difficulty}, stack...)
		}

		//48. GASLIMIT
		if op == 0x45 {
			if len(block.Gaslimit) == 0 {
				return nil, false
			}

			gaslimit, _ := new(big.Int).SetString(block.Gaslimit[2:], 16)

			stack = append([]*big.Int{gaslimit}, stack...)
		}

		//49. CHAINID
		if op == 0x46 {
			if len(block.Chainid) == 0 {
				return nil, false
			}

			chainid, _ := new(big.Int).SetString(block.Chainid[2:], 16)

			stack = append([]*big.Int{chainid}, stack...)
		}

		//50. BLOCKHASH
		if op == 0x40 {
			// number := stack[0]
			stack = stack[1:]

			// currentNumber, _ := new(big.Int).SetString(block.Number[2:], 16)
						
			//** Normally, you would find a block with the block number given in the stack and use its blockhash

			//** if block is between current block - 256 and current block , then it is valid

			//** What is implemented here is just to fulfill the test
			
			// if currentNumber.Cmp(number) < 1 && (currentNumber.Sub(currentNumber, number)).Cmp(big.NewInt(256)) > 0 {
				// return nil, false 
			// }

			stack = append([]*big.Int{big.NewInt(int64(block.Blockhash))}, stack...)
		}

		// 51. BALANCE
		if op == 0x31 {
			if len(stack) < 1 {
				return nil, false
			}

			address := stack[0]
			stack = stack[1:]

			accounts := state.Accounts

			account := accounts["0x" + address.Text(16)]

			balance := account.Balance

			if len(balance) == 0 {
				stack = append([]*big.Int{big.NewInt(0)}, stack...)
			} else {
				balanceInt, _ := new(big.Int).SetString(balance[2:], 16)

				stack = append([]*big.Int{balanceInt}, stack...)
			}
		}

		//52. CALLVALUE
		if op == 0x34 {
			valueString := transaction.Value
			
			if len(valueString) == 0 {
				return nil, false
			}

			valueInt, _ := new(big.Int).SetString(valueString[2:], 16)

			stack = append([]*big.Int{valueInt}, stack...)
		}

		//53. CALLDATALOAD
		if op == 0x35 {
			if len(stack) < 1 {
				return nil, false 
			}

			offset := int(stack[0].Int64())
			stack = stack[1:]

			data := transaction.Data
			dataBytes , _ := hex.DecodeString(data)
			
			var requiredData []byte;

			if offset + 32 <= len(dataBytes) {
				requiredData = dataBytes[offset:offset+32]
			} else {
				for i:=0; offset+32>=len(dataBytes); i++ { 
					dataBytes = append(dataBytes, byte(0))
				}
				
				requiredData = dataBytes[offset:offset+32]
			}
			_ = requiredData

			requiredDataInt := new(big.Int).SetBytes(requiredData)
			stack = append([]*big.Int{requiredDataInt}, stack...)
		}

		//54. CALLDATASIZE
		if op == 0x36 {
			data := transaction.Data
			dataBytes , _ := hex.DecodeString(data)
			
			stack = append([]*big.Int{big.NewInt(int64(len(dataBytes)))}, stack...)
		}

		//55. CALLDATACOPY
		if op == 0x37 {
			if len(stack) < 3 {
				return nil, false
			}

			memoryOffset := int(stack[0].Int64())
			calldataOffset := int(stack[1].Int64())
			sizeCalldata := int(stack[2].Int64())

			stack = stack[3:]

			data := transaction.Data
			dataBytes , _ := hex.DecodeString(data)
			
			dataCopy := dataBytes[calldataOffset : calldataOffset + sizeCalldata]
			m.Mstore(dataCopy, memoryOffset)
		}

		//56. CODESIZE
		if op == 0x38 {
			stack = append([]*big.Int{big.NewInt(int64(len(code)))}, stack...)
		}
	}	
		
		return stack, successOrNot
}


func negative_converter(number *big.Int) *big.Int {   // TODO: Can use this function more in earlier opcodes 
	bitWidth := 256
	msb := new(big.Int).Lsh(big.NewInt(1), uint(bitWidth-1))

	if new(big.Int).And(number, msb).Cmp(big.NewInt(0)) != 0 {
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitWidth))
		mask.Sub(mask, big.NewInt(1)) 

		number.Xor(number, mask) 
		number.Add(number, big.NewInt(1)) 
		number.Mul(number, big.NewInt(-1))
	}

	return number
}
