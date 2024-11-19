package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// List implements a /plot list command which may be used to check the available plots.
type List struct {
	List cmd.SubCommand `cmd:"list"`
}

// Run ...
func (l List) Run(source cmd.Source, output *cmd.Output, _ *world.Tx) {
	p := source.(*player.Player)
	h, _ := plot.LookupHandler(p)
	plots := h.Plots()

	var str strings.Builder
	for i, p := range plots {
		c := p.ColourToFormat()
		str.WriteString(text.Colourf("<white>%v:</white> <%v>■ %v</%v>", i+1, c, p.ColourToString(), c))
		if i != len(plots)-1 {
			str.WriteString("\n")
		}
	}
	output.Printf(text.Colourf("<green>Your plots:</green>\n" + str.String()))
}
