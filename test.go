package main

import (
	"math/big"
	"fmt"
)

func main() {
	value , _ := new(big.Int).SetString("455867356320691211509944977504407603390036387149619137164185182714736811808", 10)
	stack := []*big.Int{big.NewInt(0), value}

	m := NewMemory(0)

	fmt.Println("400000 in bytes:", big.NewInt(400000).Bytes())

	op := 0x00
	i := 0

	
	for i = 0; i < 10; i++ {
		
		switch i {
		case 0:
			op = 0x52
			fmt.Println("Stack round", i , ":", stack)
		case 1:
			stack = []*big.Int{big.NewInt(61)}
			op = 0x51
		case 2: 
			return
			// op = 0x59

		case 4:
			return 
		}

		//35. MSTORE
		if op == 0x52 || op == 0x53{
			if len(stack) < 2 {
				// return nil, false
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
				fmt.Println("Memory:", m)
			}
		}

		//36. MLOAD
		if op == 0x51 {
			offset := stack[0]

			value := m.MLoad(int(offset.Int64()))

			stack = append([]*big.Int{new(big.Int).SetBytes(value)}, stack[1:]...)
		}

		//37. MSIZE 
		if op == 0x59 {
			stack = append([]*big.Int{big.NewInt(int64(m.MSIZE(-32)))}, stack...)
		}

		fmt.Println("Stack round", i , ":", stack)

	}

}

type Memory struct {
	data  []byte
	// offsetMax int
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
	fmt.Println("Current memory length:", len(m.data))
	return len(m.data)
}

func (m *Memory) MLoad(offset int) []byte {
	m.MSIZE(offset)
	return m.data[offset:offset + 32]
}