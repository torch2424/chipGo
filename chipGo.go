package main

/*
    This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

    Requires Go < 1.6, because graphics library breaks on Go 1.6

    http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

import (
    "os"
    "bufio"
    "fmt"
    "flag"
    "runtime"
    "path/filepath"
    cpu "chipGo/cpu"
    graphics "chipGo/graphics"
    audio "chipGo/sound"
    input "chipGo/input"
)

//Our CPU
var chipCpu cpu.Cpu
//Our graphics
var video graphics.Video
//Our Sound
var sound audio.AudioPlayer

//Settings
var debugMode *bool

func main() {

    // This is needed to arrange that main() runs on main thread.
    // See documentation for functions that are only allowed to be called from the main thread.
    runtime.LockOSThread()

    //Print our banner
    printBanner()

    //Get our flags
    debugMode = flag.Bool("d", false, "debug mode")
    flag.Parse()

    //Get user input
    inputArgs := flag.Args()

    //Check if we have args
    var gamePath string
    if len(inputArgs) < 1 {
        printUsage()
        os.Exit(0);
    } else {
        //Set the fist argument as the file path
        gamePath = inputArgs[0]
    }

    //Check if our input file exists
    _, err := os.Stat(gamePath)
    if err != nil {
        // no such file or dir
        fmt.Println("File Not Found: ", gamePath)
        print("\n")
        panic(err)
    }

    //Inform user we are starting!
    print("Starting chipGo!\n")

    //Test our graphics
    video := graphics.NewVideo(*debugMode)

    //Start our sound
    sound := audio.NewAudioPlayer(*debugMode)

    //Set our input handler
	video.Window.SetKeyCallback(input.KeyCallback)

    //Initialize our CPU. Input is handled by opcode.go in cpu package
    chipCpu := cpu.NewCpu("chipCpu", *debugMode)
    print("Cpu initialized, name: " + chipCpu.CpuName + "\n")

    //Load the game
    loadGame, _ := filepath.Abs(gamePath)
    chipCpu = cpu.LoadGame(loadGame, chipCpu)


    //Run the game while the video is open
    for graphics.IsOpen(video) {

        //Poll for events
        graphics.PollEvents()

        //Use the Cpu Clock to see if we should run an instruction
        //Check for if our cpu clock timer has ticked
        select {
        case <- chipCpu.Clock:

            //Timer ticked
            //Run the instruction
            chipCpu = cpu.EmulateCycle(chipCpu)

            //Render our display
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
            break
        }

        //If debug mode wait for user input to continue
        if *debugMode {
            reader := bufio.NewReader(os.Stdin)
            fmt.Print("Debug Mode On. Press enter to continue...\n")
            text, _ := reader.ReadString('\n')
            fmt.Print(text);
            print("\n\n\n")
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
