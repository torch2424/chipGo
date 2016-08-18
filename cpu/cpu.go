package cpu

/*
    This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

    http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

//unsigned short = uint16
//unsigned Char = uint8

//Import io/ioutil for file reading/writing
//Import video for shared constants
import (
    "io/ioutil"
    video "chipGo/graphics"
)

type Cpu struct {

    //Name our cpu
    CpuName string

    //Our current Opcode, 2 bytes long, int16 = 16 bits = 2 bytes
    //uint = unsigned int
    // https://tour.golang.org/basics/11
    currentOpcode uint16

    //Capture the memory of of the CHIP-8, 4KB
    chipMemory []byte

    //Our CPU registers, they hold values to be used by the CPU
    //15 registers V0, V1, ... VE
    registers [16]uint8

    //Next we need our Index Register (I) and Program Counter (PC)
    //Index Register contains the value in arithmetic to be applied to a stored value
    //The program counter stores the address of the current instrstuion (or line) of the program currently being exectued by the CPU
    indexRegister uint16
    programCounter uint16

    /*
        For Reference: The System Memory Map

        0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
        0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
        0x200-0xFFF - Program ROM and work RAM
    */

    //Chip-8 has a 64 x 32 screen size, and only black our white display. So, create an array if a pizel is black (0) or white (1)
    graphicsDisplay [video.Height][video.Width]uint8

    //Chip-8 has timers, they simply count down to zero when set
    delayTimer uint8
    soundTimer uint8

    //For Goto, and jumps into functions, we need to have a stack, and point to where we currently are on the stack
    stack [16]uint16
    stackPointer uint16

    //Create our keypad
    keyPad [16]uint8
}

//Declare our function to load a game
func LoadGame(fileName string, cpu Cpu) Cpu {

    //Read the bytes of the file into memory
    game, err := ioutil.ReadFile(fileName)
    if(err != nil) {
        print("Failed loading game...")
        panic(err)
    } else {
        print("Game loaded! Value: " + string(cpu.chipMemory))
    }
    cpu.chipMemory = game
    return cpu
}
