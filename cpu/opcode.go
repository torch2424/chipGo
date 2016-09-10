package cpu

/*
   This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

   http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

//This is helper class for decoding and handling opCodes for chip-8

//Imports
import (
	graphics "github.com/torch2424/chipGo/graphics"
	input "github.com/torch2424/chipGo/input"
	"fmt"
	"math/rand"
)

//Function to return an opCode
func GetOpcode(cpu Cpu) uint16 {
	opCode := uint16(cpu.chipMemory[cpu.programCounter])<<8 | uint16(cpu.chipMemory[cpu.programCounter+1])

	//print(fmt.Sprintf("Read opcode: 0x%X\n", opCode));
	return opCode
}

/*
   Function to Decode a chip-8 opCode
   See here for opCodes: https://en.wikipedia.org/wiki/CHIP-8#Opcode_table
   Please note, Values like NNN are variables, and are supposed to be replaced with the hex value

   Much thanks to https://github.com/ejholmes/chip8/blob/master/chip8.go
   Definitely helped in understanding the operations, and what they meant
*/
func DecodeOpcode(cpu Cpu) Cpu {

	//Using bitwise & and | in order to grab nibbles from our two byte opCode. E.g & 0xF000 will return the first nibble

	//Save our opCode for easier reference
	opCode := cpu.currentOpcode

	//Switch statements for every opCode
	//Outer switch finds 4 bits, inner switch decodes last 4 bits of code
	switch opCode & 0xF000 {
	case 0x0000:
		switch opCode & 0x00FF {
		case 0x00E0:
			//Clear the screen
			cpu.ClearScreen = true
			break
		case 0x00EE:
			//Exit from subroutine
			//To do this, we need to set the program counter to the top of the stack, and then subtract one from the stack pointer
			cpu.programCounter = cpu.stack[cpu.stackPointer]
			cpu.stackPointer--
			break
		default:
			noOpcode(opCode)
		}
	case 0x1000:
		//Jump to the adress at the last 3 nibbles (NNN)
		cpu.programCounter = opCode & 0x0FFF

		//Skip program counter since we jumped
		cpu.skipProgramCounter = true
		break
	case 0x2000:
		//Call subroutine in the last 3 nibbles (NNN)

		//Increment our stack ponter
		cpu.stackPointer++

		//Place the current operation on the stack
		cpu.stack[cpu.stackPointer] = cpu.programCounter
		cpu.programCounter = opCode & 0x0FFF

		//Skip program counter since we jumped
		cpu.skipProgramCounter = true
		break
	case 0x3000:
		//Skip to next instruction if Register X equals last byte
		regX := (opCode & 0x0F00) >> 8
		lastByte := byte(opCode)

		//Skip instruction by increasing program counter
		if cpu.registers[regX] == lastByte {
			cpu.programCounter = cpu.programCounter + 2
		}
		break
	case 0x4000:
		//Same as 0x3000 but not equal
		regX := (opCode & 0x0F00) >> 8
		lastByte := byte(opCode)

		//Skip instruction by increasing program counter
		if cpu.registers[regX] != lastByte {
			cpu.programCounter = cpu.programCounter + 2
		}
		break
	case 0x5000:
		//Skip to the next instruction if Register X equals register Y
		regX := (opCode & 0x0F00) >> 8
		regY := (opCode & 0x00F0) >> 4

		//Skip instruction by increasing program counter
		if cpu.registers[regX] == cpu.registers[regY] {
			cpu.programCounter = cpu.programCounter + 2
		}
		break
	case 0x6000:
		//Set regX to last Byte
		regX := (opCode & 0x0F00) >> 8
		lastByte := byte(opCode)

		cpu.registers[regX] = lastByte
		break
	case 0x7000:
		//Add last byte to register x
		regX := (opCode & 0x0F00) >> 8
		lastByte := byte(opCode)

		cpu.registers[regX] = cpu.registers[regX] + lastByte
		break
	case 0x8000:
		//This is all regx and regY operations, storing here for easy access
		regX := (opCode & 0x0F00) >> 8
		regY := (opCode & 0x00F0) >> 4

		switch opCode & 0x000F {
		case 0x0000:
			//Set value of Regx to regY
			cpu.registers[regX] = cpu.registers[regY]
			break
		case 0x0001:
			//Set regX to regX bitwise OR regY
			cpu.registers[regX] = cpu.registers[regX] | cpu.registers[regY]
			break
		case 0x0002:
			//Set regX to regX bitwise AND regY
			cpu.registers[regX] = cpu.registers[regX] & cpu.registers[regY]
			break
		case 0x0003:
			//Set regX to regX bitwise XOR (Exclusive or) regY
			cpu.registers[regX] = cpu.registers[regX] ^ cpu.registers[regY]
			break
		case 0x0004:
			//regx = Add regX and regY. RegF (Carry flag) is set to 1 if there is a carry. 0 if there is not
			result := uint16(cpu.registers[regX]) + uint16(cpu.registers[regY])

			var carryFlag byte
			if result > 0xFF {
				carryFlag = 1
			} else {
				carryFlag = 0
			}

			//Carry flag is in last register
			cpu.registers[15] = carryFlag

			cpu.registers[regX] = uint8(result & 0xFF)
			break
		case 0x0005:
			//regx = Subtract regX and regY. RegF (Carry flag) is set to 1 if there is NOT a borrow. 0 if there is not

			var carryFlag byte
			if cpu.registers[regX] > cpu.registers[regY] {
				carryFlag = 1
			} else {
				carryFlag = 0
			}

			result := uint16(cpu.registers[regX]) - uint16(cpu.registers[regY])

			//Carry flag is in last register
			cpu.registers[15] = carryFlag

			cpu.registers[regX] = uint8(result & 0xFF)
			break
		case 0x0006:
			//Shifts regX right by one. carry flag is set to the value of the least significant bit of regX before the shift.

			//Carry flag is in last register
			var carryFlag byte
			if (cpu.registers[regX] & 0x01) == 0x01 {
				carryFlag = 1
			} else {
				carryFlag = 0
			}

			cpu.registers[15] = carryFlag

			cpu.registers[regX] = cpu.registers[regX] >> 1
			break
		case 0x0007:
			//regx = Subtract regY and regX. RegF (Carry flag) is set to 1 if there is NOT a borrow. 0 if there is not

			var carryFlag byte
			if cpu.registers[regY] > cpu.registers[regX] {
				carryFlag = 1
			} else {
				carryFlag = 0
			}

			result := uint16(cpu.registers[regY]) - uint16(cpu.registers[regX])

			//Carry flag is in last register
			cpu.registers[15] = carryFlag

			cpu.registers[regX] = uint8(result & 0xFF)
			break
		case 0x000E:
			//Shifts regX left by one. carry flag is set to the value of the most significant bit of regX before the shift.

			//Carry flag is in last register
			var carryFlag byte
			if (cpu.registers[regX] & 0x80) == 0x80 {
				carryFlag = 1
			} else {
				carryFlag = 0
			}

			cpu.registers[15] = carryFlag

			cpu.registers[regX] = cpu.registers[regX] << 1
			break
		}
	case 0x9000:
		//Skip instruction if Regx != RegY
		regX := (opCode & 0x0F00) >> 8
		regY := (opCode & 0x00F0) >> 4

		if cpu.registers[regX] != cpu.registers[regY] {
			cpu.programCounter = cpu.programCounter + 2
		}
		break
	case 0xA000:
		//Set index register to last three nibbles
		lastThree := opCode & 0x0FFF

		cpu.indexRegister = lastThree
		break
	case 0xB000:
		//Jump to the adress in last three nibbles, plus Register 0
		lastThree := opCode & 0x0FFF

		cpu.programCounter = uint16(cpu.registers[0]) + lastThree

		//Skip program counter since we jumped
		cpu.skipProgramCounter = true
		break
	case 0xC000:
		//Set register X to bitwise and of last byte and random number
		regX := (opCode & 0x0F00) >> 8
		lastByte := byte(opCode)
		//255 because highest number in a byte
		ranByte := rand.Intn(255)

		cpu.registers[regX] = lastByte & byte(ranByte)
		break
	case 0xD000:
		//Sprites stored in memory at location in index register (I), 8bits wide. Wraps around the screen. If when drawn, clears a pixel, last register (carry flag) is set to 1 otherwise it is zero. All drawing is XOR drawing (i.e. it toggles the screen pixels). Sprites are drawn starting at position regX, regY. N is the number of 8bit rows that need to be drawn. If N is greater than 1, second line continues at position regX, regY+1, and so on.

		//This gets really confusing see the section Handling graphics and input on http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/

		//Get our register indexes
		regX := (opCode & 0x0F00) >> 8
		regY := (opCode & 0x00F0) >> 4

		//Get our X coordinate
		xCoorBase := cpu.registers[regX]
		yCoorBase := cpu.registers[regY]

		//Height of the sprite is the last nibble
		spriteHeight := opCode & 0x000F

		//Memory read to create the sprite. Starting at index Register to spriteHeight
		//The colon in the array index [] is a slice, it will return a sub array in the range
		spriteRegisters := cpu.chipMemory[cpu.indexRegister : cpu.indexRegister+spriteHeight]

		//Go through our graphics array to set the values of the sprite
		//Creating a boolean to check for collision (if a pixel was already on)
		var collision bool
		for i := 0; i < len(spriteRegisters); i++ {
			//Y Axis (Column)
			//Our current pixel byte we retrieve from memory
			//This pixel byte is a row of 8 pixel values
			pixelRow := spriteRegisters[i]

			for j := 0; j < 8; j++ {
				//X-axis (Rows), each sprite can only by a byte wide (8 bits)

				//Go bit by bit to see if the pixel is already set to one, if it is, there is a collision, else, simply set it
				//E.g 0x80 = 1000 000, so pixel and bit checks if first bit of pixel is 1
				bit := 0x80 >> uint8(j)

				//Check if the pixel has the bit on
				if pixelRow&byte(bit) != 0 {
					//Turn the pixel on

					//Get our true x and y corrdinates
					//If greater than width or height, over flow back around to zero
					xCoor := xCoorBase + uint8(j)
					yCoor := yCoorBase + uint8(i)

					//Get out graphics as unsigned ints
					uIntWidth := uint8(graphics.Width)
					uIntHeight := uint8(graphics.Height)

					for xCoor >= uIntWidth {
						xCoor = xCoor - uIntWidth
					}

					for yCoor >= uIntHeight {
						yCoor = yCoor - uIntHeight
					}

					//First check if the pixel was already on
					if cpu.GraphicsDisplay[xCoor][yCoor] == 1 {
						if DebugMode {
							fmt.Printf("Collision! found at %d, %d\n", xCoor, yCoor)
						}
						collision = true
					}

					//Set the pixel to on using XOR, ^= in go
					//XOR is true if values are 1 and 0. if 1 and 1, then zero. if Zero and Zero, then Zero
					cpu.GraphicsDisplay[xCoor][yCoor] ^= 1
				}
			}

		}

		//Set our carry flag to true or false depending on collision
		if collision {
			cpu.registers[15] = 1
		} else {
			cpu.registers[15] = 0
		}

		//Set Should Render to true
		cpu.ShouldRender = true
		break
	case 0xE000:
		//Check for key presses at regX
		regX := (opCode & 0x0F00) >> 8
		regKey := cpu.registers[regX]

		//Get keys
		cpu.keyPad, _ = input.GetKeyArray()

		switch opCode & 0x000F {
		case 0x000E:
			//Skips to the next instruction if the Key stored in RegX is pressed
			if cpu.keyPad[regKey] == true {
				cpu.programCounter = cpu.programCounter + 2
			}
			break
		case 0x0001:
			//skips if not pressed
			if cpu.keyPad[regKey] == false {
				cpu.programCounter = cpu.programCounter + 2
			}
			break
		}
	case 0xF000:
		//All going to be RegX manipulations
		regX := (opCode & 0x0F00) >> 8
		//Switch Based off of third nibble
		switch opCode & 0x00FF {
		case 0x0007:
			//Set reg X to the value of the delay timer
			cpu.registers[regX] = cpu.delayTimer
			break
		case 0x000A:
			//Wait for a key press, and then set the pressed key to regX
			var keyPressed bool
			cpu.keyPad, keyPressed = input.GetKeyArray()

			if keyPressed {
				//Loop to find which key was pressed
				var keyIndex uint8
				for i := 0; i < len(cpu.keyPad); i++ {
					if cpu.keyPad[i] == true {
						keyIndex = uint8(i)
						i = len(cpu.keyPad)
					}
				}

				//Set the key index to register X
				cpu.registers[regX] = keyIndex
			} else {
				//Come back to this opcode, since we are waiting for a key press
				cpu.programCounter = cpu.programCounter - 2
			}
			break
		case 0x0015:
			//Set the delay timer to regX
			cpu.delayTimer = cpu.registers[regX]
			break
		case 0x0018:
			//Set the sound timer to regX
			cpu.soundTimer = cpu.registers[regX]
			break
		case 0x001E:
			//Add RegX to index register
			cpu.indexRegister = cpu.indexRegister + uint16(cpu.registers[regX])
			break
		case 0x0029:
			//Sets indexRegister to the location of the sprite for the character in regX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.

			//Not quite sure why it is 0x05 being multiplied
			cpu.indexRegister = uint16(cpu.registers[regX]) * uint16(0x05)
			break
		case 0x033:
			// I = index register. Stores the binary-coded decimal representation of regX, with the most significant of three digits at the address in index register, the middle digit at indexregoster plus 1, and the least significant digit at indexregister plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)

			//Hundreds
			cpu.chipMemory[cpu.indexRegister] = cpu.registers[regX] / 100

			//tenths
			cpu.chipMemory[cpu.indexRegister+1] = (cpu.registers[regX] / 10) % 10

			//Single
			cpu.chipMemory[cpu.indexRegister+2] = (cpu.registers[regX] % 100) % 10
			break
		case 0x0055:
			//Store Register zero tozero to regX including regX starting at address indexregister

			//Loop zero to regX
			for i := uint16(0); i <= regX; i++ {
				cpu.chipMemory[cpu.indexRegister+i] = cpu.registers[i]
			}
			break
		case 0x0065:
			//Same as above, but fill the registers instead of storing
			//Loop zero to regX
			for i := uint16(0); i <= regX; i++ {
				cpu.registers[i] = cpu.chipMemory[cpu.indexRegister+i]
			}
			break
		}
	default:
		noOpcode(opCode)
	}

	//Return the cpu
	return cpu
}

func noOpcode(opCode uint16) {
	err := fmt.Sprintf("\nChipGo Error! Unrecognized opCode! Decimal: %d, Hex: 0x%X\n", opCode, opCode)
	panic(err)
}
