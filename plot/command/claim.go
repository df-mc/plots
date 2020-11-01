package command

import (
	"github.com/df-mc/dragonfly/dragonfly/block"
	"github.com/df-mc/dragonfly/dragonfly/block/colour"
	"github.com/df-mc/dragonfly/dragonfly/cmd"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/df-mc/dragonfly/dragonfly/world"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Claim implements the claim command.
type Claim struct {
	Sub claim
}

// Run ...
func (Claim) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	h, _ := plot.LookupHandler(p)

	blockPos := world.BlockPosFromVec3(p.Position())
	pos := plot.PosFromBlockPos(blockPos, h.Settings())

	min, max := pos.Bounds(h.Settings())

	if !plot.Within(blockPos, min, max) {
		output.Error("You are not currently in a plot.")
		return
	}
	if current, err := h.DB().Plot(pos); err == nil {
		output.Errorf("This plot is already claimed by %v.", current.OwnerName)
		return
	}
	plots := h.Plots()
	if len(plots) >= h.Settings().MaximumPlots {
		output.Errorf("You have reached the maximum amount of plot claims. (%v/%v)", len(plots), h.Settings().MaximumPlots)
		return
	}
	c := generateRandomColour(plots)

	newPlot := &plot.Plot{OwnerName: p.Name(), Owner: p.UUID(), Colour: c.String()}
	if err := h.DB().StorePlot(pos, newPlot); err != nil {
		output.Errorf("Failed claiming plot, please try again later. (%v)", err)
		return
	}
	if err := h.SetPlotPositions(append(h.PlotPositions(), pos)); err != nil {
		output.Errorf("Failed claiming plot, please try again later. (%v)", err)
		return
	}
	b := block.Concrete{Colour: c}
	for x := -1; x < h.Settings().PlotWidth+1; x++ {
		for z := -1; z < h.Settings().PlotWidth+1; z++ {
			if x == -1 || x == h.Settings().PlotWidth || z == -1 || z == h.Settings().PlotWidth {
				p.World().SetBlock(min.Add(world.BlockPos{x, 22, z}), b)
			}
		}
	}
	f := newPlot.ColourToFormat()
	output.Printf(text.Colourf("<%v>■</%v> <green>Successfully claimed the plot. (%v/%v)</green>", f, f, len(plots)+1), h.Settings().MaximumPlots)
}

// generateRandomColour generates a random colour based on the colours of existing plots. Where possible, a
// colour that has not yet been used will be selected.
func generateRandomColour(existing []*plot.Plot) colour.Colour {
	if len(existing) >= 16 {
		return colour.All()[rand.Intn(16)]
	}
	for {
		found := true
		c := colour.All()[rand.Intn(16)]
		for _, p := range existing {
			if c.String() == p.Colour {
				// We generated a colour that we already used before, so try again.
				found = false
				break
			}
		}
		if found {
			return c
		}
	}
}

// claim ...
type claim string

// SubName ...
func (claim) SubName() string {
	return "claim"
}
