package pluginPersistency

// TODO: use package fpr type defs
type xLog struct {
	Message   string `json:"message"`
	TimeStamp int    `json:"timestamp"`
	Stream    string `json:"stream"`
}

type xLogPlugin interface {
	HandleLog(streamId string, log xLog)
}
type xPlugin interface {
	LogPlugin() xLogPlugin
}
type Config struct {
	Adapter          StorageAdapter
	StoragePath      string
	StreamsBlacklist []string
	StreamsWhitelist []string
}

type PersistencyPlugin struct {
	config    Config
	logPlugin xLogPlugin
}

func NewPersistencyPlugin(config Config) (xPlugin, error) {
	plugin := PersistencyPlugin{
		config:    config,
		logPlugin: newLogPlugin(config),
	}

	err := config.Adapter.Initialize(config.StoragePath)

	return plugin, err
}

func (p PersistencyPlugin) LogPlugin() xLogPlugin {
	return p.logPlugin
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

func (p logPlugin) HandleLog(streamId string, log xLog) {
	if allowed, ok := p.allowList[streamId]; !ok || allowed {
		p.config.Adapter.AppendLogs(streamId, log.Message)
	}
}
