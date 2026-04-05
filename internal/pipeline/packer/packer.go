package packer

type Rect struct {
	X, Y, W, H int
}

type PlacedSprite struct {
	Name string
	Rect Rect
	Path string
}
