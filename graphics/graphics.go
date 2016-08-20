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
var ColorBg = sf.Color{0, 0, 0, 0}

var ColorSprite = sf.Color{255, 255, 255, 255}

//Chip8 display size
const Width int = 64
const Height int = 32

type Video struct {

	//Our window object
	Window *glfw.Window

	//Our target object we render to
	Target *sf.RenderTarget

	//Our array of Sprites(Individual pixels) we ar erendering
	pixels []*Pixel
}

//Constructor for video
func NewVideo() Video{

	video := Video{Window: GetWindow(), Target: GetTarget()}

	//Initialize our slice
	video.pixels = make([]*Pixel, 0)

	return video
}

//Function to poll events from our glfw library
func PollEvents() {
	glfw.PollEvents()
}

//Get a window for a struct
func GetWindow() *glfw.Window {

	//Inform user of graphics initialization
	print("Opening Window...\n")

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
func Render(video Video, display [Width][Height]uint8) {

	//Loop through and create our sprites
	for i := 0; i < Height; i++ {
		//Y coordinate
		for j := 0; j < Width; j++ {
			//X Corrdinate

			if display[j][i] == 1 {
				print(1)
				//Create a sprite at the location
				video.pixels = append(video.pixels, NewPixel(float32(j * scale), float32(i * scale), float32(scale), float32(scale)))
			} else {
				print(" ")
			}

			if j >= Width - 1 {
				print("\n")
			}
		}
	}

	print("\n\n\n\n")

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
