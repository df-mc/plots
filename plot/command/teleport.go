package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
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

	f := pl.ColourToFormat()
	output.Printf(text.Colourf("<%v>■</%v> <green>Successfully teleported to your plot.</green>", f, f))
}

// tp ...
type tp string

// SubName ...
func (tp) SubName() string {
	return "tp"
}
