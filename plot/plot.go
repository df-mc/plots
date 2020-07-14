package plot

import (
	"github.com/df-mc/dragonfly/dragonfly/block/colour"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// Plot represents a plot in the world. Each plot has an owner
type Plot struct {
	// Owner is the UUID of the owner of the plot. The owner has administrative permissions over the plot such
	// as being able to add helpers to the plot.
	Owner uuid.UUID
	// OwnerName is the name last recorded for the owner.
	OwnerName string
	// Helpers is a list of helpers added to the plot. These helpers may edit the plot, but are unable to, for
	// example, add other helpers.
	Helpers []uuid.UUID
	// Colour is the colour of the plot. The border of the plot will have this colour and the colour will be
	// used to refer to different chunks owned by the player.
	Colour string
}

// Owned checks if the Plot is currently owned.
func (p *Plot) Owned() bool {
	return p.Owner != uuid.UUID{}
}

var white, green = text.White(), text.Green()

// Info returns a string of info about the Plot.
func (p *Plot) Info() string {
	if !p.Owned() {
		return white(" This plot is currently", green("free")+white("."), "\n", white("   Use", green("/p claim"), "to claim it."))
	}
	c := p.ColourToFormat()
	return white(c("■"), "Now entering", green(p.OwnerName)+white("'s plot."))
}

// ColourToFormat converts the colour of the plot to a text.FormatFunc and returns it.
func (p *Plot) ColourToFormat() text.FormatFunc {
	c, _ := colour.Colour{}.FromString(p.Colour)
	switch c.(colour.Colour) {
	default:
		return text.White()
	case colour.Orange():
		return text.Gold()
	case colour.Magenta():
		return text.Purple()
	case colour.LightBlue():
		return text.Aqua()
	case colour.Yellow():
		return text.Yellow()
	case colour.Lime():
		return text.Green()
	case colour.Pink():
		return text.Red()
	case colour.Grey():
		return text.DarkGrey()
	case colour.LightGrey():
		return text.Grey()
	case colour.Cyan():
		return text.Blue()
	case colour.Purple():
		return text.DarkPurple()
	case colour.Blue():
		return text.DarkBlue()
	case colour.Brown():
		return text.DarkYellow()
	case colour.Green():
		return text.DarkGreen()
	case colour.Red():
		return text.DarkRed()
	case colour.Black():
		return text.Black()
	}
}

// ColourToString converts the colour of the plot to a readable representation.
func (p *Plot) ColourToString() string {
	return strings.Title(strings.Replace(p.Colour, "_", " ", -1))
}
