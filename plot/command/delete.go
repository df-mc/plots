package command

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Delete implements a /p delete command, which may be used to clear a plot and delete the claim.
type Delete struct {
	Delete cmd.SubCommand `cmd:"delete"`
}

// Run ...
func (d Delete) Run(source cmd.Source, output *cmd.Output, tx *world.Tx) {
	p := source.(*player.Player)
	h, _ := plot.LookupHandler(p)

	blockPos := cube.PosFromVec3(p.Position())
	pos := plot.PosFromBlockPos(blockPos, h.Settings())

	min, max := pos.Bounds(h.Settings())

	if !plot.Within(blockPos, min, max) {
		output.Error("You are not currently in a plot.")
		return
	}
	current, err := h.DB().Plot(pos)
	if err != nil || current.Owner != p.UUID() {
		output.Errorf("You cannot delete this plot because you do not own it.")
		return
	}
	plots := h.Plots()

	if err := h.DB().RemovePlot(pos); err != nil {
		output.Errorf("Failed deleting plot, please try again later. (%v)", err)
		return
	}
	newPositions := make([]plot.Position, 0, len(plots)-1)
	for _, plotPos := range h.PlotPositions() {
		if plotPos == pos {
			continue
		}
		newPositions = append(newPositions, plotPos)
	}
	if err := h.SetPlotPositions(newPositions); err != nil {
		output.Errorf("Failed deleting plot, please try again later. (%v)", err)
		return
	}
	pos.Reset(tx, h.Settings())
	for x := -1; x < h.Settings().PlotWidth+1; x++ {
		for z := -1; z < h.Settings().PlotWidth+1; z++ {
			if x == -1 || x == h.Settings().PlotWidth || z == -1 || z == h.Settings().PlotWidth {
				tx.SetBlock(min.Add(cube.Pos{x, 22, z}), h.Settings().BoundaryBlock, opts)
			}
		}
	}
	f := current.ColourToFormat()
	output.Printf(text.Colourf("<%v>■</%v> <green>Successfully deleted the plot. (%v/%v)</green>", f, f, len(plots)-1, h.Settings().MaximumPlots))
}
