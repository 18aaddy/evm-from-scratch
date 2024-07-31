package main

import (
	"math/big"
	"fmt"
)

func LSH(value, shiftAmount *big.Int) *big.Int {
	// Perform the left shift
	shiftInt := shiftAmount.Uint64()
	result := new(big.Int).Lsh(value, uint(shiftInt))

	// Apply the 256-bit mask
	mask := new(big.Int).Lsh(big.NewInt(1), 256) // 2^256
	mask.Sub(mask, big.NewInt(1))               // 2^256 - 1
	result.And(result, mask)

	return result
}

func main() {
	value, _ := new(big.Int).SetString("123456789012345678901234567890", 10)
	shiftAmount := big.NewInt(8) // Example shift amount

	result := LSH(value, shiftAmount)

	max_uint := new(big.Int).Lsh(big.NewInt(1), 256)
	max_uint.Sub(max_uint, big.NewInt(1))
	max_uint_2 := 1 >> 256
	_ = max_uint_2

	fmt.Printf("Value: %s\n", value.String())
	fmt.Printf("Shift Amount: %s\n", shiftAmount.String())
	fmt.Printf("Result: %s\n", result.String())
	fmt.Println("Testing something:", max_uint.Text(2))

	op := 0x1d
	number := new(big.Int).Sub(max_uint, big.NewInt(7))

	stack := []*big.Int{big.NewInt(4), number}
	fmt.Println("Current Stack:", stack)

	fmt.Println("Negative Converter:", negative_converter(number))

	
	if op == 0x1d {
		// if len(stack) < 2 {
			// 	return nil, false
			// }
			
		shiftAmount := stack[0]
		value := stack[1]
		stack = stack[2:]
		
		value = negative_converter(value)
			
		fmt.Println("Shift Amount:",shiftAmount)
		fmt.Println("Value:")
		fmt.Println("Value First Digit:", value.Text(2))
	
		shiftInt := int(shiftAmount.Uint64())


	
		if shiftInt >= 256 {
			if value.Sign() < 0 {
				// If the original number is negative, fill with all 1s
				result := new(big.Int).Lsh(big.NewInt(1), 256)
				result.Sub(result, big.NewInt(1)) // 2^256 - 1
				stack = append([]*big.Int{result}, stack...)
			} else {
				// If the original number is positive, fill with 0s
				stack = append([]*big.Int{big.NewInt(0)}, stack...)
			}
		} else {
			result := new(big.Int).Rsh(value, uint(shiftInt))
	
			if value.Sign() < 0 {
				// Create a mask for the sign extension
				signBit := new(big.Int).Lsh(big.NewInt(1), 256-uint(shiftInt))
				signBit.Sub(signBit, big.NewInt(1))
				// Shift the mask to the left to fill with 1s
				extended := new(big.Int).Lsh(big.NewInt(1), uint(shiftInt))
				extended.Sub(extended, big.NewInt(1))
				extended.Lsh(extended, uint(256-uint(shiftInt)))
				
				fmt.Println("Extended:", extended)
				fmt.Println("Extended digits:", len(extended.Text(2)))

				result.Or(result, extended)
			}
	
			stack = append([]*big.Int{result}, stack...)

			fmt.Println("Final Answer:",stack[0].Text(2))
			fmt.Println("Final Answer digits:",len(stack[0].Text(2)))
		}
	}

	fmt.Println("Negative Converter on positive:", negative_converter(big.NewInt(24)))
	fmt.Println("Negative Converter on negative:", negative_converter(big.NewInt(-24)))
	fmt.Println("Length of this string is:", len("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff602a1a"))
}

func negative_converter(number *big.Int) *big.Int {
	bitWidth := 256
	msb := new(big.Int).Lsh(big.NewInt(1), uint(bitWidth-1))

	if new(big.Int).And(number, msb).Cmp(big.NewInt(0)) != 0 {
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitWidth))
		mask.Sub(mask, big.NewInt(1)) // 2^bitWidth - 1

		number.Xor(number, mask) // Invert the bits
		number.Add(number, big.NewInt(1)) // Add 1
		number.Mul(number, big.NewInt(-1)) // Negate the number
	}

	return number
}