package plot

import (
	"github.com/df-mc/dragonfly/server/block"
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

// Info returns a string of info about the Plot.
func (p *Plot) Info() string {
	if !p.Owned() {
		return text.Colourf("<white> This plot is currently <green>free</green>.</white>\n<white>   Use <green>/p claim</green> to claim it.")
	}
	c := p.ColourToFormat()
	return text.Colourf("<%v>■</%v> <white>Now entering <green>%v</green>'s plot.", c, c, p.OwnerName)
}

// ColourToFormat converts the colour of the plot to a text.FormatFunc and returns it.
func (p *Plot) ColourToFormat() string {
	c, _ := block.Colour{}.FromString(p.Colour)
	switch c.(block.Colour) {
	default:
		return "white"
	case block.ColourOrange():
		return "gold"
	case block.ColourMagenta():
		return "purple"
	case block.ColourLightBlue():
		return "aqua"
	case block.ColourYellow():
		return "yellow"
	case block.ColourLime():
		return "green"
	case block.ColourPink():
		return "red"
	case block.ColourGrey():
		return "dark-grey"
	case block.ColourLightGrey():
		return "grey"
	case block.ColourCyan():
		return "blue"
	case block.ColourPurple():
		return "dark-purple"
	case block.ColourBlue():
		return "dark-blue"
	case block.ColourBrown():
		return "dark-yellow"
	case block.ColourGreen():
		return "dark-green"
	case block.ColourRed():
		return "dark-red"
	case block.ColourBlack():
		return "black"
	}
}

// ColourToString converts the colour of the plot to a readable representation.
func (p *Plot) ColourToString() string {
	return strings.Title(strings.Replace(p.Colour, "_", " ", -1))
}
