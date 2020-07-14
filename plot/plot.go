package plot

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/text"
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
	// Name is the name of the plot. This name is displayed when entering the plot.
	Name string
	// Description is the description of the plot. Like the name, this description is displayed when entering
	// the plot.
	Description string
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
	var add string
	if p.Description != "" {
		add = white("\n", green(p.Description))
	}
	return white("Now entering plot", green(p.Name)+white("."), add)
}
