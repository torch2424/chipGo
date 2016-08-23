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
const scale int = 10

//Our colors for the display
var ColorBg = sf.Color{25, 25, 25, 255}

var ColorSprite = sf.Color{242, 242, 242, 255}

//Chip8 display size
const Width int = 64
const Height int = 32

//Debug mode
var debugMode bool

type Video struct {

	//Our window object
	Window *glfw.Window

	//Our target object we render to
	Target *sf.RenderTarget

	//Our array of Sprites(Individual pixels) we ar erendering
	pixels []*Pixel
}

//Constructor for video
func NewVideo(debug bool) Video {

	//Create our video
	video := Video{Window: GetWindow(), Target: GetTarget()}

	//Initialize our slice
	video.pixels = make([]*Pixel, 0)

	//Debug mode
	debugMode = debug

	return video
}

//Function to poll events from our glfw library
func PollEvents() {
	glfw.PollEvents()
}

//Get a window for a struct
func GetWindow() *glfw.Window {

	//Inform user of graphics initialization
	print("\nOpening Window...\n")

	//Initialize graphics library and window
	glfw.Init()
	window, err := glfw.CreateWindow(Width * scale, Height * scale, "Chip Go", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

//Get a target for a struct
func GetTarget() *sf.RenderTarget {

	//Save our target and exit
	target := sf.NewRenderTarget(sf.Vector2{float32(Width * scale), float32(Height * scale)})

	return target
}

//Return if the window is open
func IsOpen(video Video) bool {
	return !video.Window.ShouldClose()
}

//Function to draw sprites to the window
//Also, display how game should look in terminal
func Render(video Video, display [Width][Height]uint8) {

	//If debug mode, title video state
	if debugMode {
		print("Display State:\n\n")
	}

	//Loop through and create our sprites
	for i := 0; i < Height; i++ {
		//Y coordinate
		for j := 0; j < Width; j++ {
			//X Corrdinate

			if display[j][i] == 1 {

				if debugMode {
					print(1)
				}
				//Create a sprite at the location
				video.pixels = append(video.pixels, NewPixel(float32(j * scale), float32(i * scale), float32(scale), float32(scale)))
			} else if debugMode {
				print(" ")
			}

			if debugMode && j >= Width - 1 {
				print("\n")
			}
		}
	}

	//Render all of the pixels
	//Clear the screen
	video.Target.Clear(ColorBg)

	//Render all of our pixels
	for i := 0; i < len(video.pixels); i++ {
		video.pixels[i].Render(video.Target)
	}

	//Swap the buffers to show the new renders
	video.Window.SwapBuffers()
}

//Function to clear the window
func Clear(video Video) {

	//Clear our pixels
	video.pixels = nil

	//video.Target.Clear(ColorBg)
	video.Window.SwapBuffers()
}
