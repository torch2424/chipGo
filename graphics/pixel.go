package graphics

import (
	"github.com/tedsta/gosfml"
	"math/rand"
)

var spriteColor sf.Color

type Pixel struct {
	position sf.Vector2
	size     sf.Vector2
}

func NewPixel(x, y, w, h float32) *Pixel {
	return &Pixel{sf.Vector2{x, y}, sf.Vector2{w, h}}
}

func (pixel Pixel) Render(target *sf.RenderTarget, ranColorMode bool) {

	//Check for random color mode
	if ranColorMode {
		spriteColor = sf.Color{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
	} else {
		spriteColor = ColorSprite
	}

	var verts [4]sf.Vertex
	verts[0] = sf.Vertex{sf.Vector2{},
		spriteColor,
		sf.Vector2{}}
	verts[1] = sf.Vertex{sf.Vector2{0, pixel.size.Y},
		spriteColor,
		sf.Vector2{0, pixel.size.Y}}
	verts[2] = sf.Vertex{pixel.size,
		spriteColor,
		pixel.size}
	verts[3] = sf.Vertex{sf.Vector2{pixel.size.X, 0},
		spriteColor,
		sf.Vector2{pixel.size.X, 0}}

	states := sf.RenderStates{sf.BlendAlpha, sf.IdentityTransform(), nil}
	states.Transform.Translate(pixel.position)
	target.Render(verts[:], sf.Quads, states)
}
