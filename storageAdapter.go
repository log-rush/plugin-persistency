package pluginPersistency

type StorageAdapter interface {
	AppendLogs(key string, logs string)
	ListLogFiles(key string) []string
	GetLogs(key string, logFile string) ([]byte, error)
	Shutdown()
}
