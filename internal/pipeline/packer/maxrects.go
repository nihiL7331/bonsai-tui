package packer

import "math"

// via github.com/juj/RectangleBinPack

type MaxRectsPacker struct {
	width     int
	height    int
	freeRects []Rect
}

func NewPacker(width int, height int) *MaxRectsPacker {
	return &MaxRectsPacker{
		width:     width,
		height:    height,
		freeRects: []Rect{{X: 0, Y: 0, W: width, H: height}},
	}
}

// finds the best free rectangle for a given width and height
// returns the placement 'Rect' and a score (lower = better)
func (p *MaxRectsPacker) findPos(width int, height int) (Rect, int) {
	bestNode := Rect{}
	bestShortSideFit := math.MaxInt32

	for _, freeRect := range p.freeRects {
		if freeRect.W >= width && freeRect.H >= height {
			leftoverX := freeRect.W - width
			leftoverY := freeRect.H - height

			shortSide := min(leftoverX, leftoverY)

			if shortSide < bestShortSideFit {
				bestNode.X = freeRect.X
				bestNode.Y = freeRect.Y
				bestNode.W = width
				bestNode.H = height
				bestShortSideFit = shortSide
			}
		}
	}

	return bestNode, bestShortSideFit
}

// takes an existing space and a placed sprite
// if they overlap, splits the free space into smaller chunks
func splitFreeNode(freeNode Rect, usedNode Rect) []Rect {
	if usedNode.X >= freeNode.X+freeNode.W || usedNode.X+usedNode.W <= freeNode.X ||
		usedNode.Y >= freeNode.Y+freeNode.H || usedNode.Y+usedNode.H <= freeNode.Y {
		return []Rect{freeNode}
	}

	var newRects []Rect

	if usedNode.Y > freeNode.Y && usedNode.Y < freeNode.Y+freeNode.H {
		newRects = append(newRects, Rect{
			X: freeNode.X, Y: freeNode.Y, W: freeNode.W, H: usedNode.Y - freeNode.Y,
		})
	}

	if usedNode.Y+usedNode.H < freeNode.Y+freeNode.H {
		newRects = append(newRects, Rect{
			X: freeNode.X, Y: usedNode.Y + usedNode.H,
			W: freeNode.W, H: freeNode.Y + freeNode.H - (usedNode.Y + usedNode.H),
		})
	}

	if usedNode.X > freeNode.X && usedNode.X < freeNode.X+freeNode.W {
		newRects = append(newRects, Rect{
			X: freeNode.X, Y: freeNode.Y, W: usedNode.X - freeNode.X, H: freeNode.H,
		})
	}

	if usedNode.X+usedNode.W < freeNode.X+freeNode.W {
		newRects = append(newRects, Rect{
			X: usedNode.X + usedNode.W, Y: freeNode.Y,
			W: freeNode.X + freeNode.W - (usedNode.X + usedNode.W), H: freeNode.H,
		})
	}

	return newRects
}

// attempts to pack a sprite of the given dims
// returns the placement rect (or rect of width 0)
func (p *MaxRectsPacker) Insert(width int, height int) Rect {
	bestNode, score := p.findPos(width, height)

	if score == math.MaxInt32 {
		return Rect{}
	}

	var nextFreeRects []Rect
	for _, freeRect := range p.freeRects {
		splits := splitFreeNode(freeRect, bestNode)
		nextFreeRects = append(nextFreeRects, splits...)
	}
	p.freeRects = nextFreeRects

	p.pruneFreeList()

	return bestNode
}

// removes free rects that are completely within another
func (p *MaxRectsPacker) pruneFreeList() {
	var pruned []Rect
	for i := 0; i < len(p.freeRects); i++ {
		isContained := false
		for j := 0; j < len(p.freeRects); j++ {
			if i != j && isContainedBy(p.freeRects[i], p.freeRects[j]) {
				isContained = true
				break
			}
		}
		if !isContained {
			pruned = append(pruned, p.freeRects[i])
		}
	}
	p.freeRects = pruned
}

// generic AABB check
func isContainedBy(a Rect, b Rect) bool {
	return a.X >= b.X && a.Y >= b.Y &&
		a.X+a.W <= b.X+b.W && a.Y+a.H <= b.Y+b.H
}
