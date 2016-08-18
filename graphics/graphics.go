package graphics

//Using SFML rewritten in Go
//See here for dependencies needed to be installed
//https://github.com/tedsta/gosfml

//Also, for vagrant development, use X Forwarding:
//http://computingforgeeks.com/how-to-enable-and-use-ssh-x11-forwarding-on-vagrant-instances/
//Need to use "vagrant ssh -- -X chipGo" to test on host machine

import (
	"github.com/go-gl/glfw3/v3.1/glfw"
	"github.com/tedsta/gosfml"
)

//Our graphics scale for the window
const scale int = 20

//Chip8 display size
const Width int = 64
const Height int = 32

var window *glfw.Window
var target           *sf.RenderTarget
var p1Score, p2Score int

//Create our screen
func GetWindow() {

	//Inform user of graphics initialization
	print("Opening Window...\n")

	//Initialize graphics library and window
	glfw.Init()
	window, err := glfw.CreateWindow(Width * scale, Height * scale, "Chip Go", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	//Save our target and exit
	target = sf.NewRenderTarget(sf.Vector2{float32(0), float32(0)})
}
