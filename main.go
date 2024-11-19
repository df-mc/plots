package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/plots/plot"
	"github.com/df-mc/plots/plot/command"
	"github.com/pelletier/go-toml"
	"log"
	"log/slog"
	"os"
)

func main() {
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	conf, err := readConfig(slog.Default())
	if err != nil {
		log.Fatalf("error reading conf file: %v", err)
	}

	settings := plot.Settings{
		FloorBlock:    block.Grass{},
		BoundaryBlock: block.StainedTerracotta{Colour: item.ColourCyan()},
		RoadBlock:     block.Concrete{Colour: item.ColourGrey()},
		PlotWidth:     32,
		MaximumPlots:  16,
	}
	conf.Generator = func(dim world.Dimension) world.Generator {
		return plot.NewGenerator(settings)
	}

	s := conf.New()
	s.CloseOnProgramEnd()

	w := s.World()
	w.SetDefaultGameMode(world.GameModeCreative)
	w.SetSpawn(cube.Pos{2, plot.RoadHeight, 2})
	w.SetTime(5000)
	w.StopTime()

	db, err := plot.OpenDB("plots", settings)
	if err != nil {
		log.Fatalf("error opening plot database: %v", err)
	}
	w.Handle(plot.NewWorldHandler(settings))
	cmd.Register(cmd.New("plot", "Manages plots and their settings.", []string{"p", "plot"},
		command.Claim{},
		command.List{},
		command.Teleport{},
		command.Delete{},
		command.Clear{},
		command.Auto{},
	))

	s.Listen()

	for p := range s.Accept() {
		p.Handle(plot.NewPlayerHandler(p.UUID(), settings, db))
	}
	_ = db.Close()
}

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return zero, nil
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
