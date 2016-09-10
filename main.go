package main

/*
   This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

   Requires Go < 1.6, because graphics library breaks on Go 1.6

   http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

import (
	"bufio"
	cpu "github.com/torch2424/chipGo/cpu"
	graphics "github.com/torch2424/chipGo/graphics"
	input "github.com/torch2424/chipGo/input"
	audio "github.com/torch2424/chipGo/sound"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

//Our CPU
var chipCpu cpu.Cpu

//Our graphics
var video graphics.Video

//Our Sound
var sound audio.AudioPlayer

//Number of opcodes to pass in debug mode
var skipDebug int

//Command Line Parser (Kingpin) Setup
var (
	app       = kingpin.New("ChipGo", "A cjip 8 emulator written in Go")
	gamePath  = kingpin.Arg("game", "Relative filepath to the game you would like to play. e.g: games/BRIX").Required().String()
	debugMode = kingpin.Flag("debug", "Debug mode. Step through the emulator per opcode, and displays status of cpu, as well as a graphics mapping.").Short('d').Bool()
	gameSpeed = kingpin.Flag("speed", "Clock speed of the game. Increase this to make the game run faster, decrease to make the game run slower. Minimum is 1. Which will execute 1 opcode per second").Default("600").Int()
	gameScale = kingpin.Flag("scale", "Increase in Scale of the game. Original Chip-8 had a 64x32 display. Scale=10 would make the display 640x320").Default("10").Int()
	partyMode = kingpin.Flag("party", "Party Mode. Who knew Emulation could get so trippy mayne?").Short('p').Bool()
)

func main() {

	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()

	//Print our banner
	printBanner()

	//Parse our input
	kingpin.Parse()

	//Check if our input file exists
	_, err := os.Stat(*gamePath)
	if err != nil {
		// no such file or dir
		fmt.Println("File Not Found: ", *gamePath)
		print("\n")
		panic(err)
	}

	//Inform user we are starting!
	print("Starting chipGo!\n")

	//Test our graphics
	video := graphics.NewVideo(*gameScale, *debugMode, *partyMode)

	//Start our sound
	sound := audio.NewAudioPlayer(*debugMode)

	//Set our input handler
	video.Window.SetKeyCallback(input.KeyCallback)

	//Initialize our CPU. Input is handled by opcode.go in cpu package
	chipCpu := cpu.NewCpu("chipCpu", *gameSpeed, *debugMode)
	print("Cpu initialized...\n")

	//Load the game
	loadGame, _ := filepath.Abs(*gamePath)
	chipCpu = cpu.LoadGame(loadGame, chipCpu)

	//Set skip debug checks
	skipDebug = 0

	//Run the game while the video is open
	for graphics.IsOpen(video) {

		//Poll for events
		graphics.PollEvents()

		//Use the Cpu Clock to see if we should run an instruction
		//Check for if our cpu clock timer has ticked
		select {
		case <-chipCpu.Clock.C:

			//Timer ticked
			//Run the instruction
			chipCpu = cpu.EmulateCycle(chipCpu)

			//Render our display
			//using go function to call in other thread using goRoutines
			if chipCpu.ShouldRender {
				graphics.Render(video, chipCpu.GraphicsDisplay)
			}
			if chipCpu.ClearScreen {
				graphics.Clear(video)
				chipCpu = cpu.ClearGraphics(chipCpu)
			}
			chipCpu.ShouldRender = false
			chipCpu.ClearScreen = false

			//Play any sounds
			if cpu.ShouldPlaySound(chipCpu) {
				audio.PlayBlip(sound)
			}

			//Exit the case
			break
		}

		//If debug mode wait for user input to continue
		if *debugMode {
			if skipDebug < 1 {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Debug Mode On. Enter a number of opcodes to execute before pausing. Or, Press enter to continue...\n")
				text, _ := reader.ReadString('\n')

				//Remove newline from text
				text = strings.Replace(text, "\n", "", -1)

				//Try to parse input to set debug check
				parseResult, err := strconv.Atoi(text)
				if err != nil {
					print("\n\nDid not find an int, continuing debug stepping...\n\n")
				} else {
					skipDebug = int(parseResult)
				}

				fmt.Print("\n")
				print("\n\n\n")
			} else {
				skipDebug--
			}
		}
	}
}

//Function to pring program banner
func printBanner() {
	print("\n\n")
	fmt.Println("   ________    _       ______    ")
	fmt.Println("  / ____/ /_  (_)___  / ____/___ ")
	fmt.Println(" / /   / __ // / __ |/ / __/ __ /")
	fmt.Println("/ /___/ / / / / /_/ / /_/ / /_/ /")
	fmt.Println("|____/_/ /_/_/ .___/|____/|____/ ")
	fmt.Println("            /_/                  ")
	print("\n\n")

}

//Function to print usage of the program
func printUsage() {

	//Print the usage
	print("USAGE:\n\n")
	print("*rom file: the path to the rom you would like to play\n")
	print("-d: flag to enable debug mode\n")

	//Inform asterisk means required
	print("\n* - fields marked with asterisk are required to play\n")

	//Spacing to end on a new line
	print("\n\n")
}
