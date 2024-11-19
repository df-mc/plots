package plot

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"slices"
	"sync"
)

// PlayerHandler handles events of a player.Player. It handles things such as preventing players from placing
// in plots that they do not own.
type PlayerHandler struct {
	player.NopHandler
	id       uuid.UUID
	settings Settings
	db       *DB
	plots    []Position
}

// LookupHandler looks up the PlayerHandler of a player.Player passed.
func LookupHandler(p *player.Player) (*PlayerHandler, bool) {
	v, ok := handlers.Load(p.UUID())
	if !ok {
		return nil, false
	}
	return v.(*PlayerHandler), true
}

// handlers holds a list of handlers that are currently open.
var handlers sync.Map

// NewPlayerHandler creates a new PlayerHandler for the player.Player passed. The Settings and DB are used to
// track which plots the player.Player can build in.
func NewPlayerHandler(id uuid.UUID, settings Settings, db *DB) *PlayerHandler {
	positions, _ := db.PlayerPlots(id)
	h := &PlayerHandler{
		id:       id,
		settings: settings,
		db:       db,
		plots:    positions,
	}
	handlers.Store(id, h)
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
	return h.db.StorePlayerPlots(h.id, positions)
}

// HandleMove shows information on the plot that the player enters.
func (h *PlayerHandler) HandleMove(ctx *player.Context, pos mgl64.Vec3, _ cube.Rotation) {
	p := ctx.V()
	newPos, oldPos := cube.PosFromVec3(pos), cube.PosFromVec3(p.Position())
	plotPos := PosFromBlockPos(newPos, h.settings)
	previous := PosFromBlockPos(oldPos, h.settings)

	min, max := plotPos.Bounds(h.settings)
	if (plotPos != previous && Within(newPos, min, max)) ||
		(plotPos == previous && (!Within(oldPos, min, max) && Within(newPos, min, max))) {
		// Player entered a plot that it wasn't in before.
		pl, err := h.db.Plot(plotPos)
		if err != nil {
			pl = &Plot{}
		}
		p.SendTip(pl.Info())
	}
}

// HandleBlockBreak prevents block breaking outside of the player's plots.
func (h *PlayerHandler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, _ *[]item.Stack, _ *int) {
	p := ctx.V()
	if !h.canEdit(pos) {
		p.Tx().PlaySound(pos.Vec3Centre(), sound.Deny{})
		p.Tx().AddParticle(pos.Vec3Centre(), particle.BlockForceField{})
		ctx.Cancel()
	}
}

// HandleBlockPlace prevents block placing outside of the player's plots.
func (h *PlayerHandler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, _ world.Block) {
	p := ctx.V()
	if !h.canEdit(pos) {
		p.Tx().PlaySound(pos.Vec3Centre(), sound.Deny{})
		p.Tx().AddParticle(pos.Vec3Centre(), particle.BlockForceField{})
		ctx.Cancel()
	}
}

// HandleItemUseOnBlock prevents using items on blocks outside of the player's plots.
func (h *PlayerHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, _ mgl64.Vec3) {
	p := ctx.V()
	held, _ := p.HeldItems()
	if _, ok := held.Item().(world.Block); !ok && (!h.canEdit(pos) || !h.canEdit(pos.Side(face))) {
		// For blocks, we don't return here but at HandleBlockPlace.
		ctx.Cancel()
	}
}

// canEdit checks if the player.Player held by the PlayerHandler is permitted to edit the block at the
// cube.Pos passed.
func (h *PlayerHandler) canEdit(pos cube.Pos) bool {
	plotPos := PosFromBlockPos(pos, h.settings)
	min, max := plotPos.Bounds(h.settings)
	if !Within(pos, min, max) {
		return false
	}
	plot, err := h.db.Plot(plotPos)
	if err != nil {
		return false
	}
	return plot.Owner == h.id || slices.Index(plot.Helpers, h.id) != -1
}

// HandleQuit removes the PlayerHandler from the Handlers map.
func (h *PlayerHandler) HandleQuit() {
	handlers.Delete(h.id)
}
