package plot

import (
	"github.com/df-mc/dragonfly/dragonfly/world"
)

// Position represents the position of a plot. These positions are similar to chunk positions, in that they
// do not represent absolute coordinates, but, instead, a coordinate based on the size of plots.
type Position [2]int

// PosFromBlockPos returns a Position that reflects the position of the plot present at that position.
func PosFromBlockPos(pos world.BlockPos, settings Settings) Position {
	fullPlotSize := pathWidth + boundaryWidth + settings.PlotWidth
	// Integers are truncated down, so negative numbers will be wrong. We need to account for those.

	if pos[0] < 0 && mod(int32(pos[0]), int32(fullPlotSize)) != 0 {
		// Negative number that would be truncated, causing the value to be one higher than required.
		pos[0] -= fullPlotSize
	}
	if pos[2] < 0 && mod(int32(pos[2]), int32(fullPlotSize)) != 0 {
		// Negative number that would be truncated, causing the value to be one higher than required.
		pos[2] -= fullPlotSize
	}
	return Position{pos[0] / fullPlotSize, pos[2] / fullPlotSize}
}

// PosFromHash returns a Position by a byte hash created using Position.Hash(). PosFromHash panics if the
// length of the hash is not 8.
func PosFromHash(h []byte) Position {
	if len(h) != 8 {
		panic("position hash must be 8 bytes long")
	}
	return Position{
		int(int32(uint32(h[0]) | uint32(h[1])<<8 | uint32(h[2])<<16 | uint32(h[3])<<24)),
		int(int32(uint32(h[4]) | uint32(h[5])<<8 | uint32(h[6])<<16 | uint32(h[7])<<24)),
	}
}

// Hash creates a hash of the position and returns it. This hash is unique per Position and may be used to do
// lookups in databases.
func (pos Position) Hash() []byte {
	a, b := int32(pos[0]), int32(pos[1])
	return []byte{
		byte(a), byte(a >> 8), byte(a >> 16), byte(a >> 24),
		byte(b), byte(b >> 8), byte(b >> 16), byte(b >> 24),
	}
}

// Bounds returns the bounds of the Plot present at this position. Blocks may only be edited within these
// block positions.
func (pos Position) Bounds(settings Settings) (min, max world.BlockPos) {
	fullPlotSize := pathWidth + boundaryWidth + settings.PlotWidth

	baseX, baseZ := pos[0]*fullPlotSize, pos[1]*fullPlotSize
	x, z := baseX+pathWidth+1, baseZ+pathWidth+1
	return world.BlockPos{x, 0, z}, world.BlockPos{
		baseX + fullPlotSize - 2,
		255,
		baseZ + fullPlotSize - 2,
	}
}

// Within checks if a world.BlockPos is within the minimum and maximum world.BlockPos passed.
func Within(pos, min, max world.BlockPos) bool {
	return (pos[0] >= min[0] && pos[0] <= max[0]) &&
		(pos[1] >= min[1] && pos[1] <= max[1]) &&
		(pos[2] >= min[2] && pos[2] <= max[2])
}
