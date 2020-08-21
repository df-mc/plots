package plot

import (
	"github.com/df-mc/dragonfly/dragonfly/event"
	"github.com/df-mc/dragonfly/dragonfly/player"
	"github.com/df-mc/dragonfly/dragonfly/world"
	"github.com/df-mc/dragonfly/dragonfly/world/particle"
	"github.com/df-mc/dragonfly/dragonfly/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"sync"
)

// PlayerHandler handles events of a player.Player. It handles things such as preventing players from placing
// in plots that they do not own.
type PlayerHandler struct {
	player.NopHandler
	settings Settings
	db       *DB
	p        *player.Player
	plots    []Position
}

// LookupHandler looks up the PlayerHandler of a player.Player passed.
func LookupHandler(p *player.Player) (*PlayerHandler, bool) {
	v, ok := handlers.Load(p)
	if !ok {
		return nil, false
	}
	return v.(*PlayerHandler), true
}

// handlers holds a list of handlers that are currently open.
var handlers sync.Map

// NewPlayerHandler creates a new PlayerHandler for the player.Player passed. The Settings and DB are used to
// track which plots the player.Player can build in.
func NewPlayerHandler(p *player.Player, settings Settings, db *DB) *PlayerHandler {
	positions, _ := db.PlayerPlots(p)
	h := &PlayerHandler{
		settings: settings,
		db:       db,
		p:        p,
		plots:    positions,
	}
	handlers.Store(p, h)
	return h
}

// Settings returns the Settings of the PlayerHandler.
func (h *PlayerHandler) Settings() Settings {
	return h.settings
}

// DB returns the plot DB of the PlayerHandler.
func (h *PlayerHandler) DB() *DB {
	return h.db
}

// PlotPositions returns positions of all plots that the PlayerHandler holds.
func (h *PlayerHandler) PlotPositions() []Position {
	return h.plots
}

// Plots returns a list of all Plots that the PlayerHandler owns.
func (h *PlayerHandler) Plots() []*Plot {
	plots := make([]*Plot, 0, len(h.plots))
	for _, pos := range h.plots {
		plot, err := h.db.Plot(pos)
		if err != nil {
			continue
		}
		plots = append(plots, plot)
	}
	return plots
}

// SetPlotPositions sets the positions of all plots that the PlayerHandler holds.
func (h *PlayerHandler) SetPlotPositions(positions []Position) error {
	h.plots = positions
	return h.db.StorePlayerPlots(h.p, positions)
}

// HandleMove shows information on the plot that the player enters.
func (h *PlayerHandler) HandleMove(_ *event.Context, pos mgl64.Vec3, _, _ float64) {
	newPos, oldPos := world.BlockPosFromVec3(pos), world.BlockPosFromVec3(h.p.Position())
	plotPos := PosFromBlockPos(newPos, h.settings)
	previous := PosFromBlockPos(oldPos, h.settings)

	min, max := plotPos.Bounds(h.settings)
	if (plotPos != previous && Within(newPos, min, max)) ||
		(plotPos == previous && (!Within(oldPos, min, max) && Within(newPos, min, max))) {
		// Player entered a plot that it wasn't in before.
		p, err := h.db.Plot(plotPos)
		if err != nil {
			p = &Plot{}
		}
		h.p.SendTip(p.Info())
	}
}

// HandleBlockBreak prevents block breaking outside of the player's plots.
func (h *PlayerHandler) HandleBlockBreak(ctx *event.Context, pos world.BlockPos) {
	if !h.canEdit(pos) {
		h.p.World().PlaySound(pos.Vec3Centre(), sound.Deny{})
		h.p.World().AddParticle(pos.Vec3Centre(), particle.BlockForceField{})
		ctx.Cancel()
	}
}

// HandleBlockPlace prevents block placing outside of the player's plots.
func (h *PlayerHandler) HandleBlockPlace(ctx *event.Context, pos world.BlockPos, _ world.Block) {
	if !h.canEdit(pos) {
		h.p.World().PlaySound(pos.Vec3Centre(), sound.Deny{})
		h.p.World().AddParticle(pos.Vec3Centre(), particle.BlockForceField{})
		ctx.Cancel()
	}
}

// HandleItemUseOnBlock prevents using items on blocks outside of the player's plots.
func (h *PlayerHandler) HandleItemUseOnBlock(ctx *event.Context, pos world.BlockPos, face world.Face, _ mgl64.Vec3) {
	held, _ := h.p.HeldItems()
	if _, ok := held.Item().(world.Block); !ok && (!h.canEdit(pos) || !h.canEdit(pos.Side(face))) {
		// For blocks, we don't return here but at HandleBlockPlace.
		ctx.Cancel()
	}
}

// canEdit checks if the player.Player held by the PlayerHandler is permitted to edit the block at the
// world.BlockPos passed.
func (h *PlayerHandler) canEdit(pos world.BlockPos) bool {
	plotPos := PosFromBlockPos(pos, h.settings)
	min, max := plotPos.Bounds(h.settings)
	if !Within(pos, min, max) {
		return false
	}
	plot, err := h.db.Plot(plotPos)
	if err != nil {
		return false
	}
	if plot.Owner == h.p.UUID() {
		return true
	}
	for _, helper := range plot.Helpers {
		if h.p.UUID() == helper {
			return true
		}
	}
	return false
}

// HandleQuit removes the PlayerHandler from the Handlers map.
func (h *PlayerHandler) HandleQuit() {
	handlers.Delete(h.p)
}
