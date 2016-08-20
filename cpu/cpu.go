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
    "time"
    graphics "chipGo/graphics"
)

//Font set that is loaded into chip 8 memory on initialization
var fontSet = [80]byte {
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Cpu struct {

    //Name our cpu
    CpuName string

    //boolean if we should render
    ShouldRender bool

    //Boolean for clearing the screen
    ClearScreen bool

    //Our current Opcode, 2 bytes long, int16 = 16 bits = 2 bytes
    //uint = unsigned int
    // https://tour.golang.org/basics/11
    currentOpcode uint16

    //Capture the memory of of the CHIP-8, 4KB
    chipMemory [4096]byte

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
    GraphicsDisplay [graphics.Width][graphics.Height]uint8

    //Chip-8 has timers, they simply count down to zero when set
    //Timer speed used to increase or decrease clock speed
    delayTimer uint8
    soundTimer uint8
    timerSpeed float32
    Clock <-chan time.Time

    //For Goto, and jumps into functions, we need to have a stack, and point to where we currently are on the stack
    stack [16]uint16
    stackPointer int

    //Create our keypad
    keyPad [16]uint8
}

//Function to construct a new CPU
func NewCpu(cpuName string) Cpu {

    cpu := Cpu{CpuName: cpuName, stackPointer: -1}
    return cpu
}

//Declare our function to load a game
func LoadGame(fileName string, cpu Cpu) Cpu {

    //Read the bytes of the file into memory
    game, err := ioutil.ReadFile(fileName)
    if err != nil {
        print("Failed loading game...\n")
        panic(err)
    } else {
        print("Game loaded!\n")
    }

    //Set our values to the initial state
    cpu.programCounter = 0x200
    cpu.currentOpcode = 0
    cpu.indexRegister = 0
    cpu.stackPointer = 0

    //Reset timers (60 cycles per second)
    cpu.delayTimer = 60
    cpu.soundTimer = 60
    cpu.timerSpeed = 1.5
    //Find our clock speed
    clockSpeed := time.Duration(60.0 * cpu.timerSpeed)
    cpu.Clock = time.Tick(time.Second / clockSpeed)

    //Load the Chip 8 fontset into memory
    for i := 0; i < 80; i++ {
        cpu.chipMemory[i] = fontSet[i]
    }

    //Load the game into memory
    for i := 0; i < len(game); i++ {
        cpu.chipMemory[i + 512] = game[i]
    }

    return cpu
}

//Function to grab an opcode to interpret
func EmulateCycle(cpu Cpu) Cpu {

    //Reset our video booleans
    cpu.ShouldRender = false
    cpu.ClearScreen = false

    //Get and Decode opcode found in package's opcode.go

    //Get our currently pressed keys

    //Count down our timers
    if cpu.delayTimer > 0 {
        cpu.delayTimer--
    }
    if cpu.soundTimer > 0 {
        cpu.soundTimer--
    }


    //Grab the opcode
    cpu.currentOpcode = GetOpcode(cpu)

    //Decode the Opcode
    cpu = DecodeOpcode(cpu)

    //Finally increase the program counter by two
    cpu.programCounter = cpu.programCounter + 2

    return cpu
}

//Function to reset our graphics display
func ClearGraphics(cpu Cpu) Cpu {

    for i := 0; i < graphics.Width; i++ {
        for j := 0; j < graphics.Height; j++ {
            cpu.GraphicsDisplay[i][j] = uint8(0)
        }
    }

    return cpu
}
