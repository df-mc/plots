package plot

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
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

	MergedDirections []cube.Direction
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
	c, _ := item.Colour{}.FromString(p.Colour)
	switch c.(item.Colour) {
	default:
		return "white"
	case item.ColourOrange():
		return "gold"
	case item.ColourMagenta():
		return "purple"
	case item.ColourLightBlue():
		return "aqua"
	case item.ColourYellow():
		return "yellow"
	case item.ColourLime():
		return "green"
	case item.ColourPink():
		return "red"
	case item.ColourGrey():
		return "dark-grey"
	case item.ColourLightGrey():
		return "grey"
	case item.ColourCyan():
		return "blue"
	case item.ColourPurple():
		return "dark-purple"
	case item.ColourBlue():
		return "dark-blue"
	case item.ColourBrown():
		return "dark-yellow"
	case item.ColourGreen():
		return "dark-green"
	case item.ColourRed():
		return "dark-red"
	case item.ColourBlack():
		return "black"
	}
}

// ColourToString converts the colour of the plot to a readable representation.
func (p *Plot) ColourToString() string {
	return strings.Title(strings.Replace(p.Colour, "_", " ", -1))
}
