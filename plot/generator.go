package plot

import (
	"github.com/df-mc/dragonfly/dragonfly/block"
	"github.com/df-mc/dragonfly/dragonfly/world"
	"github.com/df-mc/dragonfly/dragonfly/world/chunk"
)

// Generator implements a generator for a plot world. The settings of the generator are configurable,
// allowing for different results depending on the fields set.
type Generator struct {
	floor, boundary, road uint32
	width                 int
}

// NewGenerator returns a new plot Generator with the Settings passed.
func NewGenerator(s Settings) *Generator {
	floor, _ := world.BlockRuntimeID(s.FloorBlock)
	boundary, _ := world.BlockRuntimeID(s.BoundaryBlock)
	roadOuter, _ := world.BlockRuntimeID(s.RoadBlock)
	return &Generator{
		floor:    floor,
		boundary: boundary,
		road:     roadOuter,
		width:    s.PlotWidth,
	}
}

// dirt holds the runtime ID of a dirt block.
var dirt, _ = world.BlockRuntimeID(block.Dirt{})

const (
	// RoadHeight is a rough Y position of the height of the road where a player can be safely teleported.
	RoadHeight = 24
	// Base Y of 20 blocks. The floor will start at a y level of 20.
	baseY = 20
	// Path width of 5 blocks, excluding the boundary blocks.
	pathWidth = 5
	// Boundary width of 2 blocks, 1 block around all sides.
	boundaryWidth = 2
)

// GenerateChunk generates a chunk for a plot world.
func (g *Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
	// The full plot size ends up being the width of the path and the boundary width added to the actual width
	// of plots.
	fullPlotSize := int32(pathWidth + boundaryWidth + g.width)

	// Grab the absolute coordinates of the chunk position.
	baseX, baseZ := pos[0]<<4, pos[1]<<4

	// Iterate over all possible local coordinates within the chunk.
	for localX := int32(0); localX < 16; localX++ {
		for localZ := int32(0); localZ < 16; localZ++ {
			// Create the absolute coordinates of the current block.
			x, z := baseX+localX, baseZ+localZ

			localX8, localZ8 := uint8(localX), uint8(localZ)

			// Because plots are all the same, we need to base our X and Z on the full plot size, so that we
			// can put specific blocks at specific offsets.
			relativeX, relativeZ := mod(x, fullPlotSize), mod(z, fullPlotSize)

			switch {
			case relativeX < 5 || relativeZ < 5:
				// Road blocks.
				g.fill(chunk, localX8, localZ8, baseY)
				chunk.SetRuntimeID(localX8, baseY+1, localZ8, 0, g.road)
			case relativeX == 5 || relativeZ == 5 || relativeX == fullPlotSize-1 || relativeZ == fullPlotSize-1:
				// Boundary blocks.
				g.fill(chunk, localX8, localZ8, baseY+1)
				chunk.SetRuntimeID(localX8, baseY+2, localZ8, 0, g.boundary)
			default:
				// Normal plot floor blocks.
				g.fill(chunk, localX8, localZ8, baseY+1)
				chunk.SetRuntimeID(localX8, baseY+2, localZ8, 0, g.floor)
			}
		}
	}
}

// mod does a modulo operation on a and b but always returns a positive integer.
func mod(a, b int32) int32 {
	return (a%b + b) % b
}

// fill fills the column at a specific x and z in the chunk passed up to a specific height with dirt blocks.
func (g *Generator) fill(chunk *chunk.Chunk, x, z uint8, height uint8) {
	for y := uint8(0); y <= height; y++ {
		chunk.SetRuntimeID(x, y, z, 0, dirt)
	}
}
