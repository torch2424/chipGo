package main

/*
    This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

    http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

import (
    "runtime"
    cpu "chipGo/cpu"
    video "chipGo/graphics"
)

//Our CPU
var chipCpu cpu.Cpu

func main() {

    // This is needed to arrange that main() runs on main thread.
    // See documentation for functions that are only allowed to be called from the main thread.
    runtime.LockOSThread()

    //Test our graphics
    video.GetWindow()
    //TODO set up input

    //Initialize our CPU
    print("Starting chipGo!\n");
    chipCpu := cpu.Cpu{CpuName: "chipGo"}
    print("Cpu initialized, name: " + chipCpu.CpuName + "\n")
    //Load our game into the cpu
    cpu.LoadGame("./pong", chipCpu)

}
