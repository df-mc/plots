package plot

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
)

// WorldHandler handles events of the world.World, making sure liquids don't spread out of plots.
type WorldHandler struct {
	world.NopHandler
	settings Settings
}

// NewWorldHandler returns a new WorldHandler instance using the world.World and Settings passed.
func NewWorldHandler(settings Settings) *WorldHandler {
	return &WorldHandler{settings: settings}
}

// HandleLiquidFlow prevents liquid from flowing out of a plot.
func (w *WorldHandler) HandleLiquidFlow(ctx *world.Context, _, into cube.Pos, _ world.Liquid, _ world.Block) {
	fullPlotSize := int32(pathWidth + boundaryWidth + w.settings.PlotWidth)
	relativeX, relativeZ := mod(int32(into[0]), fullPlotSize), mod(int32(into[2]), fullPlotSize)

	if relativeX <= 5 || relativeZ <= 5 || relativeX >= fullPlotSize-1 || relativeZ >= fullPlotSize-1 {
		ctx.Cancel()
	}
}
