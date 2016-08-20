package main

/*
    This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

    Requires Go < 1.6, because graphics library breaks on Go 1.6

    http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

import (
    "runtime"
    "path/filepath"
    cpu "chipGo/cpu"
    graphics "chipGo/graphics"
)

//Our CPU
var chipCpu cpu.Cpu
//Our graphics
var video graphics.Video

func main() {

    // This is needed to arrange that main() runs on main thread.
    // See documentation for functions that are only allowed to be called from the main thread.
    runtime.LockOSThread()

    //Inform user we are starting!
    print("Starting chipGo!\n");

    //Test our graphics
    video := graphics.NewVideo()

    //Initialize our CPU. Input is handled by opcode.go in cpu package
    chipCpu := cpu.NewCpu("chipCpu")
    print("Cpu initialized, name: " + chipCpu.CpuName + "\n")

    //Load our game into the cpu
    gamePath, _ := filepath.Abs("./games/UFO")
    chipCpu = cpu.LoadGame(gamePath, chipCpu)


    //Run the game while the video is open
    for graphics.IsOpen(video) {

        //Use the Cpu Clock to see if we should run an instruction

        //Check for if our cpu clock timer has ticked
        select {
        case <- chipCpu.Clock:
            //Timer ticked
            //Run the instruction
            chipCpu = cpu.EmulateCycle(chipCpu)

            //Render our display
            graphics.PollEvents()
            if chipCpu.ShouldRender {
                graphics.Render(video, chipCpu.GraphicsDisplay)
            }
            if chipCpu.ClearScreen {
                graphics.Clear(video)
                chipCpu = cpu.ClearGraphics(chipCpu)
            }
            chipCpu.ShouldRender = false
            chipCpu.ClearScreen = false
            break
        }
    }
}
