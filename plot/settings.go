package plot

import "github.com/df-mc/dragonfly/dragonfly/world"

// Settings holds the settings for a plot Generator. These settings may be changed in order to change the
// appearance of the plots generated.
type Settings struct {
	// FloorBlock is the block on the floor of each plot. The floor may be changed later, but plots will have
	// this floor by default.
	FloorBlock world.Block
	// BoundaryBlock is the block used to surround plots with. These blocks cannot be changed by an individual
	// player.
	BoundaryBlock world.Block
	// RoadBlockOuter is the outer block of the pattern on the road. These blocks cannot be changed by any
	// player.
	RoadBlockOuter world.Block
	// RoadBlockInner is the inner block of the pattern on the road. These blocks cannot be changed by any
	// player.
	RoadBlockInner world.Block
	// PlotWidth is the width in blocks that each plot generated will be.
	PlotWidth int
}
