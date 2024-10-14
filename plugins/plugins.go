package plugins

import (
	"rory-pearson/pkg/log"
	"rory-pearson/plugins/commands"
)

type Config struct {
	Log log.Log
}

type Plugins struct {
	log log.Log

	Commands commands.CommandsPlugin
}

var instance *Plugins

func Initialize(c Config) (*Plugins, error) {
	if instance != nil {
		return instance, nil
	}

	instance = &Plugins{
		log:      c.Log,
		Commands: commands.CommandsPlugin{},
	}

	// Initialize plugins here
	_, err := instance.Commands.Initialize(commands.CommandsPluginConfig{})
	if err != nil {
		return nil, err
	}
	c.Log.Info().Msg("commands initialized")

	c.Log.Info().Msg("all plugins initialized")

	return instance, nil
}

func (p *Plugins) Close() {
	p.log.Info().Msg("stopping plugins")
}

func GetInstance() Plugins {
	if instance == nil {
		panic("plugins not initialized")
	}
	return *instance
}
