package plot

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Position represents the position of a plot. These positions are similar to chunk positions, in that they
// do not represent absolute coordinates, but, instead, a coordinate based on the size of plots.
type Position [2]int

// PosFromBlockPos returns a Position that reflects the position of the plot present at that position.
func PosFromBlockPos(pos cube.Pos, settings Settings) Position {
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

// Add adds a Position to the current Position and returns a new resulting Position.
func (pos Position) Add(p Position) Position {
	return Position{pos[0] + p[0], pos[1] + p[1]}
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
func (pos Position) Bounds(settings Settings) (min, max cube.Pos) {
	fullPlotSize := pathWidth + boundaryWidth + settings.PlotWidth

	baseX, baseZ := pos[0]*fullPlotSize, pos[1]*fullPlotSize
	x, z := baseX+pathWidth+1, baseZ+pathWidth+1
	return cube.Pos{x, 0, z}, cube.Pos{
		baseX + fullPlotSize - 2,
		255,
		baseZ + fullPlotSize - 2,
	}
}

// Absolute returns an absolute cube.Pos that holds the corner of the plot.
func (pos Position) Absolute(settings Settings) cube.Pos {
	fullPlotSize := pathWidth + boundaryWidth + settings.PlotWidth
	baseX, baseZ := pos[0]*fullPlotSize, pos[1]*fullPlotSize
	return cube.Pos{baseX, 0, baseZ}
}

// TeleportPosition returns an absolute mgl64.Vec3 that can be used for teleporting the player.
func (pos Position) TeleportPosition(settings Settings) mgl64.Vec3 {
	return pos.Absolute(settings).Add(cube.Pos{2, RoadHeight, 2}).Vec3Middle()
}

// Within checks if a cube.Pos is within the minimum and maximum cube.Pos passed.
func Within(pos, min, max cube.Pos) bool {
	return (pos[0] >= min[0] && pos[0] <= max[0]) &&
		(pos[1] >= min[1] && pos[1] <= max[1]) &&
		(pos[2] >= min[2] && pos[2] <= max[2])
}

// Reset resets the Plot at the Position in the world.World passed. The Settings are used to determine the
// bounds of the plot.
func (pos Position) Reset(tx *world.Tx, settings Settings) {
	base := pos.Absolute(settings).Add(cube.Pos{pathWidth + 1, 0, pathWidth + 1})
	tx.BuildStructure(base, &resetter{settings: settings})
}

// resetter is a world.Structure implements that handles the fast resetting of chunks.
type resetter struct {
	settings Settings
}

// Dimensions returns the dimensions of a plot.
func (r *resetter) Dimensions() [3]int {
	return [3]int{
		r.settings.PlotWidth,
		256,
		r.settings.PlotWidth,
	}
}

// At returns either dirt, the floor block or air, depending on the y value.
func (r *resetter) At(_, y, _ int, _ func(x int, y int, z int) world.Block) (world.Block, world.Liquid) {
	switch {
	case y < 22:
		return block.Dirt{}, nil
	case y == 22:
		return r.settings.FloorBlock, nil
	default:
		return block.Air{}, nil
	}
}
