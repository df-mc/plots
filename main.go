package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/plots/plot"
	"github.com/df-mc/plots/plot/command"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := readConfig()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}

	s := server.New(&config, log)
	s.CloseOnProgramEnd()
	if err := s.Start(); err != nil {
		log.Fatalln(err)
	}
	w := s.World()
	w.SetDefaultGameMode(world.GameModeCreative{})
	w.SetSpawn(cube.Pos{2, plot.RoadHeight, 2})
	w.SetTime(5000)
	w.StopTime()

	settings := plot.Settings{
		FloorBlock:    block.Grass{},
		BoundaryBlock: block.StainedTerracotta{Colour: block.ColourCyan()},
		RoadBlock:     block.Concrete{Colour: block.ColourGrey()},
		PlotWidth:     32,
		MaximumPlots:  16,
	}
	db, err := plot.OpenDB("plots", settings)
	if err != nil {
		log.Fatalf("error opening plot database: %v", err)
	}
	w.Generator(plot.NewGenerator(settings))
	w.Handle(plot.NewWorldHandler(w, settings))
	cmd.Register(cmd.New("plot", "Manages plots and their settings.", []string{"p", "plot"},
		command.Claim{}, command.List{}, command.Teleport{}, command.Delete{}, command.Clear{}, command.Auto{}))

	for {
		p, err := s.Accept()
		if err != nil {
			break
		}
		p.Handle(plot.NewPlayerHandler(p, settings, db))
	}
	_ = db.Close()
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
