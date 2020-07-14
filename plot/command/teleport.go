package command

import (
	"github.com/df-mc/dragonfly/dragonfly/cmd"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"reflect"
)

// Teleport implements a /plot tp command which may be used to teleport to a specific plot owned by the
// player.
type Teleport struct {
	Sub tp
	// Number is the number of the plot to teleport to. These numbers may be found by running /p list.
	Number int `name:"number"`
}

// Run ...
func (t Teleport) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	h, _ := plot.LookupHandler(p)

	plotPositions := h.PlotPositions()
	if t.Number < 1 || t.Number > len(plotPositions) {
		output.Errorf("Unknown plot with number %v. Use /p list to get a list of plots to teleport to.", t.Number)
		return
	}
	pl := h.Plots()[t.Number-1]
	pos := plotPositions[t.Number-1]

	p.Teleport(pos.TeleportPosition(h.Settings()))
	output.Printf(text.Green()(pl.ColourToFormat()("■"), "Successfully teleported to your plot."))
}

// tp ...
type tp string

// Type ...
func (tp) Type() string {
	return "tp"
}

// Options ...
func (tp) Options() []string {
	return []string{"tp"}
}

// SetOption ...
func (tp) SetOption(string, reflect.Value) {}
