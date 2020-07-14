package command

import (
	"github.com/df-mc/dragonfly/dragonfly/cmd"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/df-mc/plots/plot"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"reflect"
	"strconv"
	"strings"
)

// List implements a /plot list command which may be used to check the available plots.
type List struct {
	Sub list
}

// Run ...
func (l List) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	h, _ := plot.LookupHandler(p)
	plots := h.Plots()

	var str strings.Builder
	for i, p := range plots {
		c := p.ColourToFormat()
		str.WriteString(text.White()(strconv.Itoa(i+1)+":", c("■", p.ColourToString())))
		if i != len(plots)-1 {
			str.WriteString("\n")
		}
	}
	output.Printf(text.Green()("Your plots:\n" + str.String()))
}

// list ...
type list string

// Type ...
func (list) Type() string {
	return "list"
}

// Options ...
func (list) Options() []string {
	return []string{"list"}
}

// SetOption ...
func (list) SetOption(string, reflect.Value) {}
