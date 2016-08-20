package graphics

import (
	"github.com/tedsta/gosfml"
)

type Pixel struct {
	position sf.Vector2
	size sf.Vector2
}

func NewPixel(x, y, w, h float32) *Pixel {
	return &Pixel{sf.Vector2{x, y}, sf.Vector2{w, h}}
}

func (pixel Pixel) Render(target *sf.RenderTarget) {
	var verts [4]sf.Vertex
	verts[0] = sf.Vertex{sf.Vector2{},
		ColorSprite,
		sf.Vector2{}}
	verts[1] = sf.Vertex{sf.Vector2{0, pixel.size.Y},
		ColorSprite,
		sf.Vector2{0, pixel.size.Y}}
	verts[2] = sf.Vertex{pixel.size,
		ColorSprite,
		pixel.size}
	verts[3] = sf.Vertex{sf.Vector2{pixel.size.X, 0},
		ColorSprite,
		sf.Vector2{pixel.size.X, 0}}

	states := sf.RenderStates{sf.BlendAlpha, sf.IdentityTransform(), nil}
	states.Transform.Translate(pixel.position)
	target.Render(verts[:], sf.Quads, states)
}
