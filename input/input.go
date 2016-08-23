package input

/*
    This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

    http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

//This is helper class for handling Input in go
//Using glfw to send key events to the window, which we can then record
//https://github.com/tedsta/gosfml/blob/master/examples/pong/main.go

//Imports
import (
    "github.com/go-gl/glfw3/v3.1/glfw"
)

//How we are mapping our keys to chip go keyboard
//Key pad is counting in hex
// Keypad                   Keyboard
// +-+-+-+-+                +-+-+-+-+
// |1|2|3|C|                |1|2|3|4|
// +-+-+-+-+                +-+-+-+-+
// |4|5|6|D|                |Q|W|E|R|
// +-+-+-+-+       =>       +-+-+-+-+
// |7|8|9|E|                |A|S|D|F|
// +-+-+-+-+                +-+-+-+-+
// |A|0|B|F|                |Z|X|C|V|
// +-+-+-+-+                +-+-+-+-+

//Our keys we will be watching
//Left side is keypad, right side is keyboard
var keyZero = []glfw.Key{glfw.KeyX}
var keyOne = []glfw.Key{glfw.Key1}
var keyTwo = []glfw.Key{glfw.Key2}
var keyThree = []glfw.Key{glfw.Key3}
var keyFour = []glfw.Key{glfw.KeyQ}
var keyFive = []glfw.Key{glfw.KeyW}
var keySix = []glfw.Key{glfw.KeyE}
var keySeven = []glfw.Key{glfw.KeyA}
var keyEight = []glfw.Key{glfw.KeyS}
var keyNine = []glfw.Key{glfw.KeyD}
var keyA = []glfw.Key{glfw.KeyZ}
var keyB = []glfw.Key{glfw.KeyC}
var keyC = []glfw.Key{glfw.Key4}
var keyD = []glfw.Key{glfw.KeyR}
var keyE = []glfw.Key{glfw.KeyF}
var keyF = []glfw.Key{glfw.KeyV}

//Place our keys into an array
var keyArray = [][]glfw.Key{keyZero, keyOne, keyTwo, keyThree, keyFour, keyFive, keySix, keySeven, keyEight, keyNine, keyA, keyB, keyC, keyD, keyE, keyF}

//Array of boolean saying if key is pressed (0 - F on keypad)
var pressedKeys [16]bool

func GetKeyArray() ([16]bool, bool) {

    //Declare if we found a key that was pressed
    keyPressed := false
    for i:= 0; i < len(pressedKeys); i++ {
        if pressedKeys[i] == true {
            keyPressed = true
            break
        }
    }


    return pressedKeys, keyPressed
}


func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {

        //Loop through our keys to turn on
        for i := 0; i < len(keyArray); i++ {
            for j:= 0; j < len(keyArray[i]); j++ {
                if key == keyArray[i][j] {
                    pressedKeys[i] = true;
                    return;
                }
            }
        }

	} else if action == glfw.Release {

        //Loop through our keys to turn off
        for i := 0; i < len(keyArray); i++ {
            for j:= 0; j < len(keyArray[i]); j++ {
                if key == keyArray[i][j] {
                    pressedKeys[i] = false;
                    return;
                }
            }
        }
	}
}
