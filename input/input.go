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
//Our mapping to keys (keymap). Key is the glfwKey, and the value is the index of the key on hex chip 8 keyboard
var keyMap = map[glfw.Key]int{
	//Zero
	glfw.KeyX: 0,
	//One
	glfw.Key1: 1,
	//Two
	glfw.Key2: 2,
	//Three
	glfw.Key3: 3,
	//Four
	glfw.KeyQ: 4,
	//Five
	glfw.KeyW: 5,
	//Six
	glfw.KeyE: 6,
	//Seven
	glfw.KeyA: 7,
	//Eight
	glfw.KeyS: 8,
	//Nine
	glfw.KeyD: 9,
	//A (10)
	glfw.KeyZ: 10,
	//B (11)
	glfw.KeyC: 11,
	//C (12)
	glfw.Key4: 12,
	//D (13)
	glfw.KeyR: 13,
	//E (14)
	glfw.KeyF: 14,
	//F (15)
	glfw.KeyV: 15,
}

//Array of boolean saying if key is pressed (0 - F on keypad)
var pressedKeys [16]bool

func GetKeyArray() ([16]bool, bool) {

	//Declare if we found a key that was pressed
	keyPressed := false

	//Loop through our currently pressed keys
	for i := 0; i < len(pressedKeys); i++ {
		if pressedKeys[i] == true {
			keyPressed = true
			break
		}
	}

	return pressedKeys, keyPressed
}

func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	//First use a two value assignment to check for key existance
	//https://blog.golang.org/go-maps-in-action
	_, validKey := keyMap[key]

	//Get the key's value from our keymap
	if validKey {

		//Get the keys value
		keyPressed := keyMap[key]

		if action == glfw.Press {

			//Set pressed keys to true
			pressedKeys[keyPressed] = true

		} else if action == glfw.Release {

			//Set pressed keys to false
			pressedKeys[keyPressed] = false
		}
	}
}
