package main

import (
	"fmt"
	"github.com/df-mc/dragonfly/dragonfly"
	"github.com/df-mc/dragonfly/dragonfly/block"
	"github.com/df-mc/dragonfly/dragonfly/block/colour"
	"github.com/df-mc/dragonfly/dragonfly/cmd"
	"github.com/df-mc/dragonfly/dragonfly/player/chat"
	"github.com/df-mc/dragonfly/dragonfly/world/gamemode"
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

	server := dragonfly.New(&config, log)
	server.CloseOnProgramEnd()
	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
	w := server.World()
	w.SetDefaultGameMode(gamemode.Creative{})

	spawn := w.Spawn()
	spawn[1] = 23
	w.SetSpawn(spawn)

	settings := plot.Settings{
		FloorBlock:     block.Grass{},
		BoundaryBlock:  block.Concrete{Colour: colour.LightGrey()},
		RoadBlockOuter: block.Concrete{Colour: colour.Grey()},
		RoadBlockInner: block.Wool{Colour: colour.Grey()},
		PlotWidth:      128,
	}
	db, err := plot.OpenDB("plots", settings)
	if err != nil {
		log.Fatalf("error opening plot database: %v", err)
	}
	w.Generator(plot.NewGenerator(settings))
	w.Handle(plot.NewWorldHandler(w, settings))
	cmd.Register(cmd.New("plot", "Manages plots and their settings", []string{"p", "plot"}, command.Claim{}))

	for {
		p, err := server.Accept()
		if err != nil {
			return
		}
		p.Handle(plot.NewPlayerHandler(p, settings, db))
	}
}

// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (dragonfly.Config, error) {
	c := dragonfly.DefaultConfig()
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
