package input

/*
    This project will be a CHIP-8 emulator in go, for a basic understanding of laearning how to code an emulator

    http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/
*/

//This is helper class for handling Input in go
//http://stackoverflow.com/questions/27198193/how-to-react-to-keypress-events-in-go

//Imports
import (
    "azul3d.org/engine/keyboard"
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

//Declare our watchers for the class
var watcher = keyboard.NewWatcher()
// Query for the map containing information about all keys
var status = watcher.States()

//Our keys we will be watching
//Left side is keypad, right side is keyboard
var keyZero = []keyboard.State{status[keyboard.X]}
var keyOne = []keyboard.State{status[keyboard.One]}
var keyTwo = []keyboard.State{status[keyboard.Two]}
var keyThree = []keyboard.State{status[keyboard.Three]}
var keyFour = []keyboard.State{status[keyboard.Q]}
var keyFive = []keyboard.State{status[keyboard.W]}
var keySix = []keyboard.State{status[keyboard.E]}
var keySeven = []keyboard.State{status[keyboard.A]}
var keyEight = []keyboard.State{status[keyboard.S]}
var keyNine = []keyboard.State{status[keyboard.D]}
var keyA = []keyboard.State{status[keyboard.Z]}
var keyB = []keyboard.State{status[keyboard.C]}
var keyC = []keyboard.State{status[keyboard.Four]}
var keyD = []keyboard.State{status[keyboard.R]}
var keyE = []keyboard.State{status[keyboard.F]}
var keyF = []keyboard.State{status[keyboard.V]}

//Place our keys into an array
var keyArray = [][]keyboard.State{keyZero, keyOne, keyTwo, keyThree, keyFour, keyFive, keySix, keySeven, keyEight, keyNine, keyA, keyB, keyC, keyD, keyE, keyF}

func GetKeyArray() ([15]uint8, bool) {

    //Declare our key array
    var pressedKeys [15]uint8

    //Declar if we found a key that was pressed
    var keyPressed bool

    //Loop through our 2d Key array
    for i := 0; i < len(keyArray); i++ {
        for j := 0; j < len(keyArray[i]); j++ {
            if keyArray[i][j] == keyboard.Down {
                pressedKeys[uint8(i)] = uint8(i)
                keyPressed = true
                //Break from the inner loop
                j = len(keyArray[i])
            }
        }
    }

    return pressedKeys, keyPressed
}
