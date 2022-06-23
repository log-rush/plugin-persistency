package pluginPersistency

import (
	logRush "github.com/log-rush/server-devkit"
	"github.com/log-rush/server-devkit/pkg/devkit"
)

type Config struct {
	Adapter          StorageAdapter
	StoragePath      string
	StreamsBlacklist []string
	StreamsWhitelist []string
}

type PersistencyPlugin struct {
	config    Config
	logPlugin logPlugin
	Plugin    logRush.Plugin
}

func NewPersistencyPlugin(config Config) (PersistencyPlugin, error) {
	plugin := PersistencyPlugin{
		config:    config,
		logPlugin: newLogPlugin(config),
	}
	plugin.Plugin = devkit.NewPlugin(plugin.logPlugin.HandleLog)

	err := config.Adapter.Initialize(config.StoragePath)

	return plugin, err
}

type logPlugin struct {
	config    Config
	allowList map[string]bool
}

func newLogPlugin(config Config) logPlugin {
	allowList := map[string]bool{}
	for _, whiteListed := range config.StreamsWhitelist {
		allowList[whiteListed] = true
	}
	for _, blackListed := range config.StreamsBlacklist {
		allowList[blackListed] = false
	}
	return logPlugin{config, allowList}
}

func (p logPlugin) HandleLog(log logRush.Log) {
	if allowed, ok := p.allowList[log.Stream]; !ok || allowed {
		p.config.Adapter.AppendLogs(log.Stream, log.Message)
	}
}
