package pluginPersistency

import (
	"strings"

	"github.com/log-rush/distribution-server/domain"
	"github.com/log-rush/distribution-server/pkg/app"
	"github.com/log-rush/distribution-server/pkg/devkit"
)

type Config struct {
	Adapter          StorageAdapter
	LogDelimiter     string
	StreamsBlacklist []string
	StreamsWhitelist []string
	ExposeLogHistory bool
}

type PersistencyPlugin struct {
	config       Config
	logPlugin    *logPlugin
	routerPlugin *routerPlugin
	Plugin       app.Plugin
}

func NewPersistencyPlugin(config Config) PersistencyPlugin {
	plugin := PersistencyPlugin{
		config:       config,
		logPlugin:    newLogPlugin(config),
		routerPlugin: newRouterPlugin(config),
	}

	p := devkit.NewPlugin(
		"persistency",
		devkit.PluginHandlers{
			LogHandler:    plugin.logPlugin.HandleLog,
			RouterHandler: plugin.routerPlugin.SetupRouter,
		},
		devkit.PluginHooks{
			OnDeInit: func(context *app.Context) {
				plugin.logPlugin.config.Adapter.Shutdown()
			},
		},
	)
	plugin.Plugin = p

	return plugin
}

type logPlugin struct {
	config    Config
	allowList map[string]bool
}

func newLogPlugin(config Config) *logPlugin {
	allowList := map[string]bool{}
	for _, whiteListed := range config.StreamsWhitelist {
		allowList[whiteListed] = true
	}
	for _, blackListed := range config.StreamsBlacklist {
		allowList[blackListed] = false
	}
	return &logPlugin{config, allowList}
}

func (p *logPlugin) HandleLog(log domain.Log) {
	if allowed, ok := p.allowList[log.Stream]; !ok || allowed {
		message := log.Message
		if !strings.HasSuffix(message, p.config.LogDelimiter) {
			message = message + p.config.LogDelimiter
		}
		p.config.Adapter.AppendLogs(log.Stream, message)
	}
}
