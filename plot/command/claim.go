package command

import (
	"github.com/df-mc/dragonfly/dragonfly/cmd"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/df-mc/dragonfly/dragonfly/world"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"reflect"
)

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
		output.Error("You are not currently in a plot")
		return
	}
	if current, err := h.DB().Plot(pos); err == nil {
		output.Errorf("This plot is already claimed by %v", current.OwnerName)
		return
	}

	newPlot := &plot.Plot{OwnerName: p.Name(), Owner: p.UUID(), Name: p.Name() + "'s Plot"}
	if err := h.DB().StorePlot(pos, newPlot); err != nil {
		output.Errorf("Failed claiming plot, please try again later (%v)", err)
		return
	}
	output.Printf(text.Green()("Successfully claimed the plot!"))
}

// claim ...
type claim string

// Type ...
func (claim) Type() string {
	return "claim"
}

// Options ...
func (claim) Options() []string {
	return []string{"claim"}
}

// SetOption ...
func (claim) SetOption(string, reflect.Value) {}
