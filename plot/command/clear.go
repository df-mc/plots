package command

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Clear implements the /plot clear command. It may be used to clear one's plot without removing the claim
// from it.
type Clear struct {
	Clear cmd.SubCommand `cmd:"clear"`
}

// Run ...
func (r Clear) Run(source cmd.Source, output *cmd.Output, tx *world.Tx) {
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
		output.Errorf("You cannot clear this plot because you do not own it.")
		return
	}
	pos.Reset(tx, h.Settings())
	f := current.ColourToFormat()
	output.Printf(text.Colourf("<%v>■</%v> <green>Successfully cleared the plot.</green>", f, f))
}
